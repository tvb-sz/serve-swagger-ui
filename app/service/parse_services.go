package service

// parseService swagger json parser service
type parseService struct{}

// Swagger parsed swagger file item
type Swagger struct {
	Url    string
	Title  string
	Desc   string
	Author string
	Email  string
	Icon   string
}

func (s *parseService) Parse(path string) (res []Swagger, err error) {
	return
}
