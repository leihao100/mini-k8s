package gpuserver

import (
	"MiniK8S/config"
	apiConfig "MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	gpuclient "MiniK8S/pkg/gpu/client"
	"MiniK8S/utils/class"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

const jobStateCheckInterval = 30 * time.Second
const defaultWorkerSleepInterval = time.Duration(3) * time.Second
const retryDownloadTimes = 10

var (
	errorStopRequested = errors.New("stop requested")
)

type Server interface {
	Run(ctx context.Context, cancel context.CancelFunc)
}

type GpuServer struct {
	cli            gpuclient.Client
	jobClient      apiClient.Client
	jobListWatcher listwatch.ListerWatcher
	jobQueue       class.ConcurrentQueue
}

type resStatus struct {
	Status string `json:"status,omitempty"`
}

func NewServer() Server {
	jobClient := apiClient.NewRESTClient(types.JobObjectType)
	jobListWatcher := listwatch.NewListWatchFromClient(jobClient)
	return &GpuServer{
		cli:            gpuclient.New(),
		jobClient:      *jobClient,
		jobListWatcher: jobListWatcher,
		jobQueue:       *class.NewConcurrentQueue(),
	}
}

func isFinished(status types.JobState) bool {
	return status == types.JobFailed || status == types.JobCompleted
}

func (s *GpuServer) Run(ctx context.Context, cancel context.CancelFunc) {

	log.Printf("[GpuServer] start\n")
	defer log.Printf("[GpuServer] init finish\n")

	s.cli.Run(ctx)

	go func() {
		defer cancel()
		err := s.listAndWatchJobs(ctx.Done())
		if err != nil {
			log.Printf("[GpuServer] listAndWatchJobs failed, err: %v\n", err)
		}
	}()

	go func() {
		defer cancel()
		s.runJobWorker(ctx)
	}()

	go func() {
		defer cancel()
		s.periodicallyCheckJobState()
	}()

}

func (s *GpuServer) periodicallyCheckJobState() {
	for {
		log.Printf("[periodicallyCheckJobState] check start\n")

		time.Sleep(jobStateCheckInterval)
		jobList, err := s.jobListWatcher.List(apiConfig.ListOptions{
			Kind:            string(types.JobObjectType),
			APIVersion:      "",
			LabelSelector:   "",
			FieldSelector:   "",
			Watch:           false,
			ResourceVersion: "",
			TimeoutSeconds:  nil,
		})
		if err != nil {
			log.Printf("[periodicallyCheckJobState] jobListWatcher list failed\n")
			continue
		}

		log.Printf("[periodicallyCheckJobState] jobList %v\n", jobList)

		jobs := jobList.GetItems()
		log.Printf("[periodicallyCheckJobState] %v jobs to check in apiserver storage\n", len(jobs))
		for i, item := range jobs {
			job := item.(*apiConfig.Job)
			jobID := job.Status.JobID
			if jobID == "" {
				log.Printf("[periodicallyCheckJobState] job %v do not have JobID now\n", i)
				continue
			}

			if isFinished(job.Status.State) {
				log.Printf("[periodicallyCheckJobState] job %v, JobID %v has already done, state %v\n", i, jobID, job.Status.State)
				continue
			}

			jobState, _ := s.cli.GetJobState(jobID)
			log.Printf("[periodicallyCheckJobState] job %v, JobID %v, state %v\n", i, jobID, jobState)
			if jobState == "" || jobState == string(types.JobMissing) || jobState == string(job.Status.State) {
				log.Printf("[periodicallyCheckJobState] job %v state not found or unchange\n", i)
				continue
			}

			job.Status.State = types.JobState(jobState)

			serverPath := s.jobClient.BuildURL(apiClient.Create)
			data, _ := job.JsonMarshal()
			res := s.jobClient.Put(serverPath, data)
			var buf bytes.Buffer
			tempBuf := make([]byte, 1024)
			for {
				n, err := res.Read(tempBuf)
				if err != nil && err != io.EOF {
					fmt.Printf("[periodicallyCheckJobState] Error: %v\n", err)
					return
				}
				if n == 0 {
					break
				}
				buf.Write(tempBuf[:n])
			}
			jsonData := buf.Bytes()
			var resStatus resStatus
			err = json.Unmarshal(jsonData, &resStatus)
			if err != nil {
				log.Printf("[periodicallyCheckJobState] JSON Unmarshal error: %v\n", err)
				return
			}
			for resStatus.Status == "FAILED" {
				getServerPath := serverPath + job.GetName()
				resp := s.jobClient.Get(getServerPath, nil)
				for {
					n, err := resp.Read(tempBuf)
					if err != nil && err != io.EOF {
						fmt.Printf("[periodicallyCheckJobState] Error: %v\n", err)
						return
					}
					if n == 0 {
						break
					}
					buf.Write(tempBuf[:n])
				}
				resp.Close()
				jsonData = buf.Bytes()
				_ = job.JsonUnmarshal(jsonData)
				log.Printf("[periodicallyCheckJobState] conflict\n")
				jobState, _ = s.cli.GetJobState(jobID)
				job.Status.State = types.JobState(jobState)
				data, _ := job.JsonMarshal()
				res.Close()
				res = s.jobClient.Put(serverPath, data)
				for {
					n, err := res.Read(tempBuf)
					if err != nil && err != io.EOF {
						fmt.Printf("[periodicallyCheckJobState] Error: %v\n", err)
						return
					}
					if n == 0 {
						break
					}
					buf.Write(tempBuf[:n])
				}
				jsonData := buf.Bytes()
				err = json.Unmarshal(jsonData, &resStatus)
				if err != nil {
					log.Printf("[periodicallyCheckJobState] JSON Unmarshal error: %v\n", err)
					return
				}
			}

			log.Printf("[periodicallyCheckJobState] job %v state update\n", i)
		}

		log.Printf("[periodicallyCheckJobState] jobs state check finish\n")
	}
}

func (s *GpuServer) runJobWorker(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Printf("[worker] ctx.Done() received, worker of GpuServer exit\n")
			return
		default:
			for s.processNextJob() {
			}
			time.Sleep(defaultWorkerSleepInterval)
		}
	}

}

func (s *GpuServer) processNextJob() bool {

	job := s.dequeueJob()
	if job == nil {
		return false
	}

	jobId, err := s.submitJob(job)
	if err != nil {
		log.Printf("[processNextJob] submit job uid %v error: %v\n", job.Metadata.Uid, err)
		return false
	}

	job.Status.JobID = jobId
	jobState, _ := s.cli.GetJobState(jobId)
	job.Status.State = types.JobState(jobState)

	serverPath := s.jobClient.BuildURL(apiClient.Create)
	data, _ := job.JsonMarshal()
	res := s.jobClient.Put(serverPath, data)
	var buf bytes.Buffer
	tempBuf := make([]byte, 1024)
	for {
		n, err := res.Read(tempBuf)
		if err != nil && err != io.EOF {
			fmt.Printf("[processNextJob] Error: %v\n", err)
			return false
		}
		if n == 0 {
			break
		}
		buf.Write(tempBuf[:n])
	}
	jsonData := buf.Bytes()
	var resStatus resStatus
	err = json.Unmarshal(jsonData, &resStatus)
	if err != nil {
		log.Printf("[processNextJob] JSON Unmarshal error: %v\n", err)
		return false
	}
	for resStatus.Status == "FAILED" {
		getServerPath := serverPath + job.GetName()
		resp := s.jobClient.Get(getServerPath, nil)
		for {
			n, err := resp.Read(tempBuf)
			if err != nil && err != io.EOF {
				fmt.Printf("[processNextJob] Error: %v\n", err)
				return false
			}
			if n == 0 {
				break
			}
			buf.Write(tempBuf[:n])
		}
		resp.Close()
		jsonData = buf.Bytes()
		_ = job.JsonUnmarshal(jsonData)
		log.Printf("[processNextJob] conflict\n")
		jobState, _ = s.cli.GetJobState(jobId)
		job.Status.State = types.JobState(jobState)
		data, _ := job.JsonMarshal()
		res.Close()
		res = s.jobClient.Put(serverPath, data)
		for {
			n, err := res.Read(tempBuf)
			if err != nil && err != io.EOF {
				fmt.Printf("[processNextJob] Error: %v\n", err)
				return false
			}
			if n == 0 {
				break
			}
			buf.Write(tempBuf[:n])
		}
		jsonData := buf.Bytes()
		err = json.Unmarshal(jsonData, &resStatus)
		if err != nil {
			log.Printf("[processNextJob] JSON Unmarshal error: %v\n", err)
			return false
		}
	}

	log.Printf("[processNextJob] submit job uid %v finish\n", job.Metadata.Uid)

	return true
}

func (s *GpuServer) listAndWatchJobs(stopCh <-chan struct{}) error {

	jobsList, err := s.jobListWatcher.List(apiConfig.ListOptions{
		Kind:            string(types.JobObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return err
	}

	jobItems := jobsList.GetItems()
	for _, item := range jobItems {
		job := item.(*apiConfig.Job)
		s.enqueueJob(job)
	}

	var w watch.Interface
	w, err = s.jobListWatcher.Watch(apiConfig.ListOptions{
		Kind:            string(types.PodObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           true,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		return err
	}

	err = s.handleWatchJobs(w, stopCh)
	w.Stop()

	if err == errorStopRequested {
		return nil
	}

	return err

}

func (s *GpuServer) enqueueJob(job *apiConfig.Job) {
	s.jobQueue.Enqueue(job)
	log.Printf("[enqueueJob] job %v enqueued\n", job.Metadata.Uid)
}

func (s *GpuServer) dequeueJob() *apiConfig.Job {
	jobItem, exist := s.jobQueue.Dequeue()
	if exist {
		j := jobItem.(*apiConfig.Job)
		log.Printf("[dequeueJob] job %v equeued\n", j.Metadata.Uid)
		return j
	} else {
		log.Printf("[dequeueJob] queue empty\n")
		return nil
	}
}

func (s *GpuServer) handleWatchJobs(w watch.Interface, stopCh <-chan struct{}) error {
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
			log.Printf("[handleWatchJobs] event %v\n", event)
			log.Printf("[handleWatchJobs] event object %v\n", event.Object)
			eventCount += 1

			switch event.Type {
			case watch.Added:
				newJob := (event.Object).(*apiConfig.Job)
				s.enqueueJob(newJob)
				log.Printf("[handleWatchJobs] new Job event, handle job %v created\n", newJob.Metadata.Uid)
			case watch.Modified:
				newJob := (event.Object).(*apiConfig.Job)
				go s.handleJobModified(newJob)
			case watch.Deleted:
				// ignore
			case watch.Bookmark:
				panic("[handleWatchJobs] watchHandler Event Type watch.Bookmark received")
			case watch.Error:
				panic("[handleWatchJobs] watchHandler Event Type watch.Error received")
			default:
				panic("[handleWatchJobs] watchHandler Unknown Event Type received")
			}
		}
	}
	return nil
}

func (s *GpuServer) submitJob(job *apiConfig.Job) (jobId string, err error) {
	slurmFile := GenerateJobScript(job)
	jobId, err = s.cli.SubmitCudaJob(job.Metadata.Uid, job.Spec.CudaFilePath, slurmFile, job.Spec.ResultFileName)
	return jobId, err
}

func (s *GpuServer) handleJobModified(job *apiConfig.Job) {
	if job.Status.State == types.JobCompleted {
		log.Printf("[handleJobModified] handling Job COMPLETED, jobId %v\n", job.Status.JobID)

		downloaded, _ := s.downloadJobResult(job)

		if !downloaded {
			job.Status.State = types.JobRunning

			serverPath := s.jobClient.BuildURL(apiClient.Create)
			data, _ := job.JsonMarshal()
			res := s.jobClient.Put(serverPath, data)
			var buf bytes.Buffer
			tempBuf := make([]byte, 1024)
			for {
				n, err := res.Read(tempBuf)
				if err != nil && err != io.EOF {
					fmt.Printf("[handleJobModified] Error: %v\n", err)
					return
				}
				if n == 0 {
					break
				}
				buf.Write(tempBuf[:n])
			}
			jsonData := buf.Bytes()
			var resStatus resStatus
			err := json.Unmarshal(jsonData, &resStatus)
			if err != nil {
				log.Printf("[handleJobModified] JSON Unmarshal error: %v\n", err)
				return
			}
			for resStatus.Status == "FAILED" {
				getServerPath := serverPath + job.GetName()
				resp := s.jobClient.Get(getServerPath, nil)
				for {
					n, err := resp.Read(tempBuf)
					if err != nil && err != io.EOF {
						fmt.Printf("[handleJobModified] Error: %v\n", err)
						return
					}
					if n == 0 {
						break
					}
					buf.Write(tempBuf[:n])
				}
				resp.Close()
				jsonData = buf.Bytes()
				_ = job.JsonUnmarshal(jsonData)
				job.Status.State = types.JobRunning
				data, _ := job.JsonMarshal()
				res.Close()
				res = s.jobClient.Put(serverPath, data)
				for {
					n, err := res.Read(tempBuf)
					if err != nil && err != io.EOF {
						fmt.Printf("[handleJobModified] Error: %v\n", err)
						return
					}
					if n == 0 {
						break
					}
					buf.Write(tempBuf[:n])
				}
				jsonData := buf.Bytes()
				err = json.Unmarshal(jsonData, &resStatus)
				if err != nil {
					log.Printf("[handleJobModified] JSON Unmarshal error: %v\n", err)
					return
				}
			}
			return
		}

		log.Printf("[handleJobModified] handling Job COMPLETED, jobID %v, result downloaded success\n", job.Status.JobID)

	} else if job.Status.State == types.JobFailed {
		log.Printf("[handleJobModified] handling Job FAILED, jobId %v\n", job.Status.JobID)
	}
}

func (s *GpuServer) downloadJobResult(job *apiConfig.Job) (downloaded bool, err error) {
	downloaded = false
	retry := 0
	for !downloaded && retry < retryDownloadTimes {
		downloaded, err = s.cli.DownloadResult(job.Metadata.Uid, job.Spec.ResultFilePath, job.Spec.ResultFileName)
		if err != nil {
			log.Printf("[handleJobModified] downloadJobResult for jobId %v failed: %v\n", job.Status.JobID, err)
		}
		retry += 1
	}
	return downloaded, err
}

func GenerateJobScript(job *apiConfig.Job) string {

	mailRemindTemplate := `#SBATCH --mail-type=%s
#SBATCH --mail-user=%s%s
`
	mailRemind := ""
	if job.Spec.Args.Mail != nil {
		mailRemind = fmt.Sprintf(
			mailRemindTemplate,
			job.Spec.Args.Mail.Type,
			job.Spec.Args.Mail.UserName,
			config.MailAddressSuffix,
		)
	}

	template := `#!/bin/bash
#SBATCH --job-name=%s
#SBATCH --partition=dgx2
#SBATCH --output=%s.out
#SBATCH --error=%s.err
#SBATCH -N 1
#SBATCH --ntasks-per-node=%d
#SBATCH --cpus-per-task=%d
#SBATCH --gres=gpu:%d
%s

ulimit -s unlimited
ulimit -l unlimited

module load gcc/11.2.0 gromacs/2022.5-gcc-11.2.0-cuda
./%s
`
	numTasksPerNode := 0
	if job.Spec.Args.NumTasksperNode != 0 {
		numTasksPerNode = job.Spec.Args.NumTasksperNode
	} else {
		numTasksPerNode = 1
	}
	cpusPerTask := 0
	if job.Spec.Args.CpusPerTask != 0 {
		cpusPerTask = job.Spec.Args.CpusPerTask
	} else {
		cpusPerTask = 1
	}
	gpuResources := 0
	if job.Spec.Args.GpuResources != 0 {
		gpuResources = job.Spec.Args.GpuResources
	} else {
		gpuResources = 1
	}

	script := fmt.Sprintf(
		template,
		job.Metadata.Uid,
		job.Spec.ResultFileName,
		job.Spec.ResultFileName,
		numTasksPerNode,
		cpusPerTask,
		gpuResources,
		mailRemind,
		job.Spec.ResultFileName,
	)
	return script
}
