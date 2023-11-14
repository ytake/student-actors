package classroom

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/message"
)

// Actor represents a classroom
type Actor struct {
	pipe          *actor.PID
	teacher       *actor.PID
	students      []*actor.PID
	endOfHomework []string
}

func NewActor(pipe *actor.PID, teacher *actor.PID, students []*actor.PID) func() actor.Actor {
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
		for _, student := range state.students {
			context.Request(student, msg)
		}
	case *message.SubmittedAchievementTest:
		endOfHomework := append(state.endOfHomework, msg.Name)
		if len(endOfHomework) == len(state.students) {
			context.Send(state.pipe, &message.EndOfAchievementTest{Subject: msg.Subject})
			context.Poison(context.Self())
		} else {
			state.endOfHomework = endOfHomework
		}
	}
}
