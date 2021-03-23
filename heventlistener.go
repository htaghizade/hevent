package hevent

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
)

type Dispatcher struct {
	jobs   chan job
	events []string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		jobs:   make(chan job),
	}
}

func (d *Dispatcher) Register(events ...Event) {
	for _, v := range events {
		d.events = append(d.events, reflect.TypeOf(v).String())
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event Event) error {
	name := reflect.TypeOf(event).String()
	for _, v := range d.events {
		if v == name {
			d.jobs <- job{event , ctx}
		}
	}

	return fmt.Errorf("%s is not a registered event", name)
}
func (d *Dispatcher) Consume() {
	for job := range d.jobs {
		err := job.event.Handle(job.ctx)
		if err!=nil {
			addEventLog("Error In Handling Event JOB ::" , job , err)
		}
		addEventLog("Event JOB Handled successfully ::" , job)
	}
}


// All custom events must satisfy this interface.
type Event interface {
	Handle(ctx context.Context) error
}
// job represents events. When a new event is dispatched, it
// gets tuned into a job and put into `Dispatcher.jobs` channel.
type job struct {
	event Event
	ctx context.Context
}

func addEventLog(msg... interface{})  {
	// log to custom file
	LOG_FILE := "tmp/event_log"
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	// Set log out put and enjoy :)
	log.SetOutput(logFile)

	// optional: log date-time, filename, and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	for _ , v := range msg {
		log.Println(v)
	}
	log.Println("----------------------------------------------------------------------------------------------------")

}


