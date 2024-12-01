package notifier

import (
	"context"
	"fmt"
)

type Notifier interface {
	SendNotification(ctx context.Context, notification Notification) error
}

type notifier struct {
}

func NewNotifier() Notifier {
	return notifier{}
}

func (n notifier) SendNotification(_ context.Context, notification Notification) error {
	fmt.Printf("[Notifier] Sending notification: %+v\n", notification)
	return nil
}
