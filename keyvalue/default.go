package keyvalue

func init() {
	RegisterImplementation("default", MemConfig{})
}

type MemConfig struct{}

func (c MemConfig) Validate() error { return nil }

func (c MemConfig) Create() (Store, error) {
	return inMemStore{
		values: map[string]interface{}{},
	}, nil
}

type inMemStore struct {
	values map[string]interface{}
}

func (s inMemStore) Get(key string) (interface{}, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
	}
	return nil, nil
}

func (s inMemStore) Set(key string, value interface{}) error {
	s.values[key] = value
	return nil
}

// for ms.UsedService interface
func (s inMemStore) Status() interface{} {
	return nil
}
