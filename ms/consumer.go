package ms

//Consumer is an asynchronous micro-service interface
//it gets requests and call the appropriate handler
//no response is sent
//
//examples are:
//	NATS subscription (without response)
//	Kafka consumer
//	Rabit MQ consumer
//	Task reader from a DB ...
//
//see also Server for sync micro-services

type ConsumerConfig interface {
	Create() (Consumer, error)
}

type Consumer interface {
	Consume() error
}
