package ms

type Address struct {
	Domain    string `json:"domain,omitempty" doc:"Domain name e.g. api.mystuff.co.za"`
	Operation string `json:"operation,omitempty" doc:"Operation name within the domain"`
}

type Client interface {
	Sync(addr Address, req interface{}) (res interface{}, err error)
	ASync(addr Address, req interface{}) (err error)
	Send(addr Address, req interface{}) (err error)
}
