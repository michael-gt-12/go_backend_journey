package Infrastructure

import (
	"context"
	"log"

	"task_manager/Domain"
)

// LogNotifier delivers a reminder by writing it to the log.
// It satisfies Usecases.Notifier; main.go checks that when it wires things up.
type LogNotifier struct{}

func NewLogNotifier() *LogNotifier { return &LogNotifier{} }

func (n *LogNotifier) Notify(ctx context.Context, task Domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	log.Printf("REMINDER task_id=%s title=%q due=%s",
		task.ID, task.Title, task.DueDate.Format("2006-01-02 15:04 MST"))
	return nil
}
