package student

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
)

type Actor struct{}

func NewActor() actor.Actor {
	return &Actor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *command.StartTest:
		// ランダムで解答時間を1-9秒で設定する
		randTime := rand.Intn(9) + 1
		time.Sleep(time.Duration(randTime) * time.Second)
		// 生徒がテストの問題を解く
		context.Logger().Info(fmt.Sprintf("%s が %s テストの解答を提出しました", context.Self().Id, msg.Subject))
		context.Respond(&command.SubmitTest{Subject: msg.Subject, Name: context.Self().Id})
		context.Poison(context.Self())
	}
}
