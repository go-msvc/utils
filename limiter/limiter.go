package limiter

type Config struct{}

func (c Config) Validate() error {
	//todo
	return nil
}

type Limiter interface {
	Allow(key string) bool
}
