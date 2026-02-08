package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"for-the-love-of-pointers/job-queue/internal/jobs/types"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/cenkalti/backoff/v5"
	"gopkg.in/gomail.v2"
)

var ProcessorFuncs = map[string]func(*Job) error{
	"printer":   func(j *Job) error { return Printer(j) },
	"send_mail": func(j *Job) error { return SendMail(j) },
	"web_hook":  func(j *Job) error { return WebHook(j) },
}

func Printer(j *Job) error {
	fmt.Println(j)
	return nil
}

func SendMail(j *Job) error {

	var payload types.EmailPayload

	if err := json.Unmarshal(j.Payload, &payload); err != nil {
		j.Status = StatusFailed
		j.Error = err.Error()
		return err
	}

	if _, err := mail.ParseAddress(payload.To); err != nil {
		j.Status = StatusFailed
		j.Error = err.Error()
		return err
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
		j.Status = StatusFailed
		j.Error = err.Error()
		return err
	}

	return nil
}

func WebHook(j *Job) error {

	var wh types.WebhookPayload
	if err := json.Unmarshal(j.Payload, &wh); err != nil {
		return err
	}

	attempt := 0

	request := func() (int, error) {
		fmt.Printf("Trying job of id %s for the %d time", j.ID, attempt)
		if wh.MaxRetries >= 0 && attempt > wh.MaxRetries {
			return 1, backoff.Permanent(errors.New("max retries exceeded"))
		}
		attempt++

		req, err := http.NewRequest(wh.Method, wh.URL, bytes.NewReader(wh.Body))
		if err != nil {
			return 2, backoff.Permanent(err)
		}

		req.Header.Set("Content-Type", "application/json")
		for k, v := range wh.Headers {
			req.Header.Set(k, v)
		}

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return 3, err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {

			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return 4, backoff.Permanent(
					fmt.Errorf("non-retryable status %d", resp.StatusCode),
				)
			}
			return 5, fmt.Errorf("retryable status %d", resp.StatusCode)
		}

		return 0, nil
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 500 * time.Millisecond
	bo.MaxInterval = 10 * time.Second
	bo.Reset()

	_, err := backoff.Retry(
		context.TODO(),
		request,
		backoff.WithBackOff(bo),
	)

	if err != nil {
		j.Status = StatusFailed
		j.Error = err.Error()
	}
	return nil
}
