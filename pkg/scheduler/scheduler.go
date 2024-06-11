package scheduler

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	podListWatcher  listwatch.ListerWatcher
	nodeListWatcher listwatch.ListerWatcher

	nodesToSchedule     []*config.Node
	nodesToScheduleLock sync.Mutex
	rrScheduleCnt       uint64

	podClient  apiClient.Interface
	nodeClient apiClient.Interface
}

var (
	errorStopRequested = errors.New("stop requested")
)

func NewScheduler() *Scheduler {
	podClient := apiClient.NewRESTClient(types.PodObjectType)
	nodeClient := apiClient.NewRESTClient(types.NodeObjectType)

	podListWatcher := listwatch.NewListWatchFromClient(podClient)
	nodeListWatcher := listwatch.NewListWatchFromClient(nodeClient)

	nodesToScheduleLock := make([]*config.Node, 0)

	return &Scheduler{
		podListWatcher:  podListWatcher,
		nodeListWatcher: nodeListWatcher,
		nodesToSchedule: nodesToScheduleLock,
		podClient:       podClient,
		nodeClient:      nodeClient,
	}
}

func (s *Scheduler) Init() {
	log.Printf("[Scheduler] init start\n")
	defer log.Printf("[Scheduler] init finish\n")
}

func (s *Scheduler) Run(ctx context.Context, cancel context.CancelFunc) {
	log.Printf("[Scheduler] run start\n")
	defer log.Printf("[Scheduler] run finish\n")
	syncChan := make(chan struct{})

	go func() {
		defer cancel()
		err := s.listAndWatchNodes(syncChan, ctx.Done())
		if err != nil {
			log.Printf("[Scheduler] listAndWatchNodes failed, err: %v\n", err)
		}
	}()

	// wait for node list finish
	<-syncChan

	go func() {
		defer cancel()
		err := s.listAndWatchPods(ctx.Done())
		if err != nil {
			log.Printf("[Scheduler] listAndWatchPods failed, err: %v\n", err)
		}
	}()
}

func isMaster(node *config.Node) bool {
	return strings.ToLower(node.Metadata.Name) == "master"
}

func (s *Scheduler) listAndWatchNodes(listDone chan<- struct{}, stopCh <-chan struct{}) error {
	// list all nodes and push into nodesQueue
	nodesList, err := s.nodeListWatcher.List(config.ListOptions{})
	if err != nil {
		return err
	}
	nodeItems := nodesList.GetItems()

	s.nodesToScheduleLock.Lock()
	for _, nodeItem := range nodeItems {
		node := nodeItem.(*config.Node)
		if isMaster(node) {
			continue
		}
		s.nodesToSchedule = append(s.nodesToSchedule, node)
	}
	s.nodesToScheduleLock.Unlock()

	close(listDone)

	w, err := s.nodeListWatcher.Watch(config.ListOptions{})
	if err != nil {
		return err
	}

	err = s.handleWatchNodes(w, stopCh)

	if errors.Is(err, errorStopRequested) {
		return nil
	}

	return err
}

func (s *Scheduler) handleWatchNodes(w watch.Interface, stopCh <-chan struct{}) error {
	eventCount := 0
loop:
	for {
		select {
		case <-stopCh:
			return errorStopRequested
		case event, ok := <-w.ResultChan():
			if !ok {
				break loop
			}
			log.Printf("[handleWatchNodes] event %v\n", event)
			log.Printf("[handleWatchNodes] event object %v\n", event.Object)
			eventCount += 1

			switch event.Type {
			case watch.Added:
				newNode := (event.Object).(*config.Node)
				// in future versions may need to check if node is running
				s.addNode(newNode)

			case watch.Modified:
				newNode := (event.Object).(*config.Node)
				// in future versions may need to check if node is running
				nodeUID := newNode.GetUID()
				s.modifyNode(nodeUID, newNode)
			case watch.Deleted:
				oldNode := (event.Object).(*config.Node)
				nodeUID := oldNode.GetUID()
				s.deleteNode(nodeUID)
			case watch.Bookmark:
				panic("[handleWatchNodes] watchHandler Event Type watch.Bookmark received")
			case watch.Error:
				panic("[handleWatchNodes] watchHandler Event Type watch.Error received")
			default:
				panic("[handleWatchNodes] watchHandler Unknown Event Type received")
			}
		}
	}
	return nil
}

func (s *Scheduler) addNode(node *config.Node) {
	if isMaster(node) {
		return
	}
	s.nodesToScheduleLock.Lock()
	s.nodesToSchedule = append(s.nodesToSchedule, node)
	s.nodesToScheduleLock.Unlock()
}

func (s *Scheduler) modifyNode(uid uuid.UUID, node *config.Node) {
	if isMaster(node) {
		return
	}
	s.nodesToScheduleLock.Lock()
	for i, nodeToModify := range s.nodesToSchedule {
		if nodeToModify.GetUID() == uid {
			s.nodesToSchedule[i] = node
			break
		}
	}
	s.nodesToScheduleLock.Unlock()
}

func (s *Scheduler) deleteNode(uid uuid.UUID) {
	s.nodesToScheduleLock.Lock()
	for i, nodeToDelete := range s.nodesToSchedule {
		if nodeToDelete.GetUID() == uid {
			s.nodesToSchedule = append(s.nodesToSchedule[:i], s.nodesToSchedule[i+1:]...)
			break
		}
	}
	s.nodesToScheduleLock.Unlock()
}

func (s *Scheduler) listAndWatchPods(stopCh <-chan struct{}) error {

	// list all pods and push into podsQueue
	podsList, err := s.podListWatcher.List(config.ListOptions{})
	if err != nil {
		return err
	}

	podItems := podsList.GetItems()
	for _, podItem := range podItems {
		pod := podItem.(*config.Pod)
		s.doSchedule(pod)
	}

	// start watch pods change
	var w watch.Interface
	w, err = s.podListWatcher.Watch(config.ListOptions{})
	if err != nil {
		return err
	}

	err = s.handleWatchPods(w, stopCh)
	w.Stop() // stop watch

	if err == errorStopRequested {
		return nil
	}

	return err

}

func (s *Scheduler) handleWatchPods(w watch.Interface, stopCh <-chan struct{}) error {
	eventCount := 0
loop:
	for {
		select {
		case <-stopCh:
			return errorStopRequested
		case event, ok := <-w.ResultChan():
			if !ok {
				break loop
			}
			log.Printf("[handleWatchPods] event %v\n", event)
			log.Printf("[handleWatchPods] event object %v\n", event.Object)
			eventCount += 1

			switch event.Type {
			case watch.Added:
				newPod := (event.Object).(*config.Pod)
				//s.enqueuePod(newPod)
				log.Printf("[handleWatchPods] new Pod event, handle pod %v created\n", newPod.GetUID())
				s.doSchedule(newPod)
			case watch.Modified:
				// ignore
			case watch.Deleted:
				// ignore
			case watch.Bookmark:
				panic("[handleWatchPods] watchHandler Event Type watch.Bookmark received")
			case watch.Error:
				panic("[handleWatchPods] watchHandler Event Type watch.Error received")
			default:
				panic("[handleWatchPods] watchHandler Unknown Event Type received")
			}
		}
	}
	return nil
}

func (s *Scheduler) doSchedule(pod *config.Pod) {
	log.Printf("[doSchedule] pod %v scheduling\n", pod.GetUID())

	nodeName := pod.Spec.NodeName
	if nodeName != "" {
		return
	}
	// wait when no node is available
	var node *config.Node = nil
	for {
		s.nodesToScheduleLock.Lock()
		node = s.roundRobin()
		s.nodesToScheduleLock.Unlock()
		if node != nil {
			break
		}
		log.Println("[doSchedule] no nodes to schedule, waiting")
		time.Sleep(1 * time.Second)
	}

	pod.Spec.NodeName = node.Metadata.Name
	// TODO: 未来改成用name
	code, err := s.podClient.PutObject(pod)

	if err != nil {
		log.Printf("[doSchedule] put pod %v to node %v failed, error %v, status code: %d\n", pod.GetUID(), node.GetUID(), err, code)
		return
	}
	log.Printf("[doSchedule] pod %v scheduled on node %v\n", pod.GetUID(), node.GetUID())
}

// should acquire the lock before calling roundRobin
func (s *Scheduler) roundRobin() *config.Node {
	log.Printf("[roundRobin] nodesToSchedule: %v\n", s.nodesToSchedule)
	if len(s.nodesToSchedule) == 0 {
		log.Println("[roundRobin] no nodes to schedule")
		return nil
	}
	node := s.nodesToSchedule[s.rrScheduleCnt%uint64(len(s.nodesToSchedule))]
	log.Printf("[roundRobin] round robin scheduled %dth node %v, name: %s\n", s.rrScheduleCnt, node.GetUID(), node.Metadata.Name)
	s.rrScheduleCnt++
	return node
}
