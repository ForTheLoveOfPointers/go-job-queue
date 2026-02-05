package jobs

import (
	"fmt"
)

var ProcessorFuncs = map[string]func(job *Job){
	"printer": func(job *Job) { Printer(job) },
}

func Printer(job *Job) {
	fmt.Println(job)
	job.Status = StatusCompleted
}
