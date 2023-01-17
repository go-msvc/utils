package streaming

type EventStream interface {
	Chan() <-chan Event //events are processed from this chan by Workflow.Run()
	Stop()              //call to stop queing start events in the chan so that we do not start new sessions
	PushStart(Event) error
	PushContinue(Event) error
}
