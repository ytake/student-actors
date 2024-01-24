package classroom

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/message"
	"github.com/ytake/student-actors/student"
)

// Actor represents a classroom
type Actor struct {
	pipe          *actor.PID
	teacher       *actor.PID
	students      []int
	endOfHomework []string
}

func NewActor(pipe *actor.PID, teacher *actor.PID, students []int) func() actor.Actor {
	return func() actor.Actor {
		return &Actor{
			pipe:          pipe,
			teacher:       teacher,
			students:      students,
			endOfHomework: []string{},
		}
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (state *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.BeginClassRequest:
		context.Request(state.teacher, msg)
	case *message.AchievementTestRequest:
		for _, st := range state.students {
			st, _ := context.SpawnNamed(
				actor.PropsFromProducer(student.NewActor),
				fmt.Sprintf("student-%d", st))
			context.Request(st, msg)
		}
	case *message.SubmittedAchievementTest:
		endOfHomework := append(state.endOfHomework, msg.Name)
		if len(endOfHomework) == len(state.students) {
			context.Request(state.teacher, &message.EndOfAchievementTest{Subject: msg.Subject})
		} else {
			state.endOfHomework = endOfHomework
		}
	case *message.ReceivedAchievementTest:
		context.Send(state.pipe, msg)
		context.Poison(context.Self())
	}
}
