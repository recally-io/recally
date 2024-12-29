package queue

import (
	"context"
	"recally/internal/pkg/logger"

	"github.com/riverqueue/river"
)

type EmailWorkerArgs struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (EmailWorkerArgs) Kind() string {
	return "send_email"
}

type EmailWorker struct {
	river.WorkerDefaults[EmailWorkerArgs]
}

func (w *EmailWorker) Work(ctx context.Context, job *river.Job[EmailWorkerArgs]) error {
	// Send email
	logger.FromContext(ctx).Info("Sending email", "from", job.Args.From, "to", job.Args.To, "subject", job.Args.Subject, "body", job.Args.Body)
	return nil
}
