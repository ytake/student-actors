package student

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/message"
)

type Actor struct {
}

func NewActor() actor.Actor {
	return &Actor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *message.AchievementTestRequest:
		// ランダムで解答時間を設定する
		randTime := rand.Intn(10-1) + 1
		time.Sleep(time.Duration(randTime) * time.Second)
		// 生徒がテストの問題を解く
		fmt.Println(ctx.Self().Id, "が", msg.Subject, "テストの解答を提出しました")
		ctx.Respond(&message.SubmittedAchievementTest{Subject: msg.Subject, Name: ctx.Self().Id})
	}
}
