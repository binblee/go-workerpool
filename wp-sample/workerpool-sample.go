package main
import (
  "fmt"
  "github.com/binblee/go-workerpool/workerpool"
)

type SomethingAsJob struct {
  name string
  result chan string
}

func (job SomethingAsJob) Do() {
  job.result <- job.name
}


func main()  {

  var workerPool workerpool.WorkerPool
  workerPool.Run()

  resultChan := make(chan string)
  for i:=0; i<100000; i++ {
    job := SomethingAsJob{ fmt.Sprintf("job #%d",i), resultChan }
    workerPool.Submit(job)
  }

  for i:=0; i<10; i++ {
    for j:=0; j<10000; j++{
      <- resultChan
    }
    fmt.Printf("got result batch: %d\n", i)
  }
}
