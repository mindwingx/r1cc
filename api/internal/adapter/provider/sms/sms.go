package sms

type Provider struct {
	Host      string
	Urls      map[string]interface{}
	AuthToken string
}

func (s *Provider) Send(phone string, msg string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}
