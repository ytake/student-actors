package classroom

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/event"
	"github.com/ytake/student-actors/student"
	"github.com/ytake/student-actors/teacher"
	"google.golang.org/protobuf/proto"
)

// Actor represents a classroom.
type Actor struct {
	pipe          *actor.PID
	teacher       *actor.PID
	students      []int
	endOfHomework []string
	mutex         sync.Mutex
	state         proto.Message
}

func NewActor(pipe *actor.PID, students []int) func() actor.Actor {
	return func() actor.Actor {
		return &Actor{
			pipe:          pipe,
			students:      students,
			endOfHomework: []string{},
			state:         nil,
		}
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor.
func (class *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *command.ClassStarts:
		class.teacher = context.Spawn(actor.PropsFromProducer(teacher.NewActor))
		context.Request(class.teacher, msg)

		class.state = &event.ClassHasStarted{Subject: msg.Subject}
	case *command.TestBegins:
		class.mutex.Lock()
		for _, st := range class.students {
			sta, err := context.SpawnNamed(
				actor.PropsFromProducer(student.NewActor),
				fmt.Sprintf("student-%d", st))
			if err != nil {
				context.Poison(context.Self())
			}
			context.Request(sta, msg)
		}
		class.mutex.Unlock()
		class.state = &event.TestWasGiven{Subject: msg.Subject}
	case *command.SubmitTest:
		endOfHomework := append(class.endOfHomework, msg.Name)
		if len(endOfHomework) == len(class.students) {
			class.state = &event.TestReceived{Subject: msg.Subject}
			context.Request(class.teacher, &command.EndTest{Subject: msg.Subject})
		} else {
			class.endOfHomework = endOfHomework
		}
	case *command.ReceiveTest:
		class.state = &event.TestFinished{Subject: msg.Subject}
		context.Send(class.pipe, class.state)
		context.Poison(context.Self())
	}
}
