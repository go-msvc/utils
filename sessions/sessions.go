package sessions

type Session interface {
	Get(name string) (interface{}, error)
	GetAll(names []string) (map[string]interface{}, error)
	Set(name string, value interface{}) error
	SetAll(map[string]interface{}) error
}

type SessionManager interface {
	Get(id string) (Session, error)
}
