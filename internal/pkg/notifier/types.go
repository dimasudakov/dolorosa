package notifier

// Notification структура с данными уведомления
type Notification struct {
	Text           string
	AlertName      string
	ClientID       string
	Amount         int64
	IdempotencyKey *string
	RuleName       string
	OperationID    string
}
