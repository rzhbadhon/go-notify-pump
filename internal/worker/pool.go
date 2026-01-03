package worker

import (
	"fmt"
	"sync"
)

type Job struct {
	ID   string
	Type string
}

type WorkerPool struct {
	MaxWorkers int
	JobQueue   chan Job
	WG         *sync.WaitGroup
}

func NewWorkerPool (maxWorkers, queueSize int) *WorkerPool{
	jobQueue := make(chan Job, queueSize)

	return &WorkerPool{
		MaxWorkers: maxWorkers,
		JobQueue: jobQueue,
		WG: &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Run(){
	
	for i := 0; i < wp.MaxWorkers; i++{
		wp.WG.Add(1)
		
		go func(workerID int){
			defer wp.WG.Done()
			fmt.Println("Worker started", workerID)

			for job := range wp.JobQueue{
				fmt.Printf("Worker %d processing job: %v \n", workerID, job)
			}
			fmt.Println("Worker stopped ", workerID)
		}(i)
	}
}

func (wp *WorkerPool) ShutDown(){
	close(wp.JobQueue)
	wp.WG.Wait()
}