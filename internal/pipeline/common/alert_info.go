package common

type AlertInfo struct {
	Msg            string
	Name           string
	ClientID       string
	Amount         int64
	IdempotencyKey string
	OperationID    string
}
