package domain

type Operation struct {
	OperationID   string
	ClientID      string
	Amount        int64
	SenderPhone   string
	ReceiverPhone string
	ReceiverBic   string
	ReceiverName  string
}

func (sbp Operation) GetClientID() string {
	return sbp.ClientID
}

func (sbp Operation) GetSenderPhone() string {
	return sbp.SenderPhone
}

func (sbp Operation) GetAmount() int64 {
	return sbp.Amount
}

func (sbp Operation) GetReceiverBic() string {
	return sbp.ReceiverBic
}
