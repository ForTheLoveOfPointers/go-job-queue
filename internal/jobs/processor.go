package jobs

import (
	"encoding/json"
	"fmt"
	"for-the-love-of-pointers/job-queue/internal/jobs/types"
	"os"

	"gopkg.in/gomail.v2"
)

var ProcessorFuncs = map[string]func(*Job){
	"printer":   func(job *Job) { Printer(job) },
	"send_mail": func(job *Job) { SendMail(job) },
}

func Printer(job *Job) {
	fmt.Println(job)
}

func SendMail(job *Job) {

	var payload types.EmailPayload

	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		fmt.Println("Could not unmarshal email payload:", err)
	}

	m := gomail.NewMessage()
	comp_mail := os.Getenv("COMPANY_MAIL")
	m.SetHeader("From", comp_mail)
	m.SetHeader("To", payload.To)
	m.SetHeader("Subject", payload.Subject)
	m.SetBody("text/plain", payload.Body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_SERVER"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	if err := d.DialAndSend(m); err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		fmt.Println("Error:", err)
	}
}
