package teacher

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/message"
)

type Actor struct{}

func NewActor() actor.Actor {
	return &Actor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *message.BeginClassRequest:
		ctx.Respond(&message.HomeworkRequest{Subject: msg.Subject})
		// ctx.Send(&message.HomeworkRequest{Subject: msg.Subject})
	}
}
