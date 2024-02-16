package teacher

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
)

type Actor struct{}

func NewActor() actor.Actor {
	return &Actor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Restarting:
		ctx.Logger().Info("先生が復活しました")
	case *command.ClassStarts:
		// 先生が宿題を出す
		ctx.Logger().Info("先生が", msg.Subject, "テストを出しました")
		ctx.Respond(&command.TestBegins{Subject: msg.Subject})
		// 宿題を提出した後に先生アクターを意図的にパニックさせる
		panic("teacher panic")
	case *command.EndTest:
		// panic後復活した先生がテストの解答を受け取る
		ctx.Logger().Info("先生が", msg.Subject, "テストの解答を受け取りました")
		ctx.Respond(&command.ReceiveTest{Subject: msg.Subject})
	}
}
