package teacher

import (
	"fmt"

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
		// 先生が宿題を出す
		fmt.Println("先生が授業に参加している生徒にテストを出しました")
		ctx.Respond(&message.AchievementTestRequest{Subject: msg.Subject})
	}
}
