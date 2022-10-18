package ms

type UsedService interface {
	Status() interface{} //include in service health response
}
