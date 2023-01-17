package ms

//Server is a synchronous micro-service interface
//it gets requests, call the appropriate handler then sends a response
//examples are:
//	HTTP REST server
//	NATS subscription on request topics (it publishes a response on a topid specified in the request)
//
//see also Consumer for async micro-services

type Server interface {
	Serve(svc MicroService) error
}
