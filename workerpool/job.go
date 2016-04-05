package workerpool

import (
  "runtime"
)

type Job interface {
  Do()
}

type worker struct {
}

func newWorker() *worker{
  return &worker{}
}

func (worker *worker) do(job Job) {
    job.Do()
}

type WorkerPool struct {
  workerChannel chan worker
  jobQueue chan Job
}

func (wp *WorkerPool) Run(){
  maxNumOfWorkers := runtime.NumCPU()
  wp.workerChannel = make(chan worker,maxNumOfWorkers)
  wp.jobQueue = make(chan Job)

  for i:=0; i < maxNumOfWorkers; i++ {
    wp.workerChannel <- *newWorker()
  }
  go func(){
    for{
      select{
      case job := <- wp.jobQueue:
        go func(){
          w := <- wp.workerChannel
          w.do(job)
          wp.workerChannel <- w
        }()
      }
    }
  }()
}

func (wp *WorkerPool) Submit(job Job) {
  wp.jobQueue <- job
}
