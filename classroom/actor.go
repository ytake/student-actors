package classroom

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/event"
	"github.com/ytake/student-actors/teacher"
)

// Actor represents a classroom.
type Actor struct {
	pipe     *actor.PID
	teacher  *actor.PID
	students []int
	mutex    sync.Mutex
}

func NewActor(pipe *actor.PID, students []int) func() actor.Actor {
	return func() actor.Actor {
		return &Actor{
			pipe:     pipe,
			students: students,
		}
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor.
func (class *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *command.StartsClass:
		class.teacher = context.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return teacher.NewActor(class.students, context.Self())
		}))
		context.Request(class.teacher, &command.PrepareTest{Subject: msg.Subject})
	case *command.FinishTest:
		context.Send(class.pipe, &event.ClassFinished{Subject: msg.Subject})
		context.Poison(context.Self())
	}
}
