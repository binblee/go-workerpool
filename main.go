package main
import (
  "fmt"
)

type Job interface {
  do()
}

type SomethingAsJob struct {
  name string
  result chan string
}

func (job SomethingAsJob) do() {
  job.result <- job.name
}

type Worker struct {
}

func NewWorker() *Worker{
  return &Worker{}
}

func (worker *Worker) do(job Job) {
    job.do()
}

type WorkerPool struct {
  workerChannel chan Worker
  jobQueue chan Job
}

func (wp *WorkerPool) run(){
  wp.workerChannel = make(chan Worker,2)
  wp.jobQueue = make(chan Job)

  for i:=0; i<2; i++ {
    wp.workerChannel <- *NewWorker()
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

func (wp *WorkerPool) submit(job Job) {
  wp.jobQueue <- job
}

func main()  {

  var workerPool WorkerPool
  workerPool.run()

  resultChan := make(chan string)
  for i:=0; i<100000; i++ {
    job := SomethingAsJob{ fmt.Sprintf("job #%d",i), resultChan }
    workerPool.submit(job)
  }

  for i:=0; i<10; i++ {
    for j:=0; j<10000; j++{
      <- resultChan
    }
    fmt.Printf("got result batch: %d\n", i)
  }
}
