package workerpool

type Job interface {
  Do()
}

type Worker struct {
}

func NewWorker() *Worker{
  return &Worker{}
}

func (worker *Worker) do(job Job) {
    job.Do()
}

type WorkerPool struct {
  workerChannel chan Worker
  jobQueue chan Job
}

func (wp *WorkerPool) Run(){
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

func (wp *WorkerPool) Submit(job Job) {
  wp.jobQueue <- job
}
