package teacher

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/message"
)

type Actor struct {
	isError bool
}

func NewActor() actor.Actor {
	return &Actor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Restarting:
		fmt.Println("先生が復活しました")
	case *message.BeginClassRequest:
		// 先生が宿題を出す
		fmt.Println("先生が授業に参加している生徒にテストを出しました")
		ctx.Respond(&message.AchievementTestRequest{Subject: msg.Subject})
		// 宿題を提出した後に先生アクターを意図的にパニックさせる
		a.isError = true
		panic("teacher panic")
	case *message.EndOfAchievementTest:
		// panic後復活した先生がテストの解答を受け取る
		// 全員がテストの解答を提出した
		fmt.Println("先生が", msg.Subject, "テストの解答を受け取りました")
		ctx.Respond(&message.ReceivedAchievementTest{Subject: msg.Subject})
	}
}
