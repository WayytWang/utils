package service

type Config interface {
	GetConfig()
}

//@inject(Config,set=service)
type AcmConfig struct {
}

func (ac *AcmConfig) GetConfig() {

}

type DB struct {
}

//@inject(set=service)
func NewDB() DB {
	return DB{}
}

type Service struct {
	DB
	Config
}
