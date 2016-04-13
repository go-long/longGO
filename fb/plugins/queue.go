package plugins

import (
	"time"

	"github.com/hoysoft/JexGO/logger"
)

type Tasks struct {
	MaxGoroutine  chan int
	tasks  []func() error
	breaktag bool
}

func NewTasks(maxGoroutineCount int)*Tasks{
	t:=new(Tasks)
	t.MaxGoroutine= make(chan int, maxGoroutineCount)
	return t
}


func (t *Tasks)Start( ){
	t.breaktag=false
	go func() {
		for {
			for len(t.tasks)>0 {

				if t.breaktag {break}
				v:=t.tasks[0]
				t.MaxGoroutine <- 1
				go func(r func() error) {
					//fmt.Println("start task:",k)
					logger.Debug("start task")
					err := r()
					logger.Debug("finish task:",err)
					<-t.MaxGoroutine
				}(v)
				//delete(t.tasks,k)
				t.tasks=append(t.tasks[:0], t.tasks[0+1:]...)
			}
			if t.breaktag {break}
			time.Sleep(time.Millisecond * 100)
		}
	}()

}


func (t *Tasks)AddTask(task func() error){
	t.tasks=append(t.tasks, task)
}



func (t *Tasks)Stop(){
	t.breaktag=true
}
