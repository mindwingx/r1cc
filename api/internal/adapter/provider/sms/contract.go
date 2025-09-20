package sms

//go:generate mockgen -source=./contract.go -destination=./sms_mock.go -package=sms
type ISmsProvider interface {
	Send(phone string, msg string) (map[string]interface{}, error)
}
