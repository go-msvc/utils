package main

type Config struct{}

func (c Config) Validate() error {
	return nil
}

func (c Config) Loaded() {
	//todo
}
