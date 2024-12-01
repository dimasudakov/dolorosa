package nirvana

type Adapter interface {
	GetClientID() string
	GetAmount() int64
}
