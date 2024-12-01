package feature

type Adapter interface {
	GetClientID() string
	GetPhoneNumber() string
}
