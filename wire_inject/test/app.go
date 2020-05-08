package test

import "utils/wire_inject/test/service"

type Application struct {
	service.Config
	service.News
	service.DB
}

//@inject(set=app)
func NewApplication(c *service.AcmConfig, n service.News, d service.DB) *Application {
	return &Application{
		Config: c,
		News:   n,
		DB:     d,
	}
}
