package classroom

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/message"
)

// Actor represents a classroom
type Actor struct {
	teacher  *actor.PID
	students []*actor.PID
}

func NewActor(teacher *actor.PID, students []*actor.PID) func() actor.Actor {
	return func() actor.Actor {
		return &Actor{
			teacher:  teacher,
			students: students,
		}
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (state *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.BeginClassRequest:
		f := context.RequestFuture(state.teacher, msg, 2*time.Second)
		r, err := f.Result()
		if err != nil {
			context.Stop(context.Self())
		}
		fmt.Println(r)
		fmt.Println(state.students)
	}
}
