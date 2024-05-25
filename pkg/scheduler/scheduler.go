package scheduler

import (
	"fmt"
)

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Init() {
	fmt.Printf("[Scheduler] init start\n")

	fmt.Printf("[Scheduler] init finish\n")
}
func (s *Scheduler) Run()        {}
func (s *Scheduler) doSchedule() {}
func (s *Scheduler) enqueuePod() {}
func (s *Scheduler) dequeuePod() {}
func (s *Scheduler) roundRobin() {}
