package teacher

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
)

type Actor struct {
	students   []int
	endOfTests []command.SubmitTest
}

func NewActor() actor.Actor {
	return &Actor{
		students:   []int{},
		endOfTests: []command.SubmitTest{},
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *command.PrepareTest:
		a.students = msg.Students
		// 先生が宿題を出す
		context.Logger().Info("先生が", msg.Subject, "テストを出しました")
		context.Respond(&command.StartTest{Subject: msg.Subject})
		// 生徒がテストを提出する
	case *command.SubmitTest:
		a.endOfTests = append(a.endOfTests, *msg)
		// 全員提出したら先生がテストの解答を受け取る
		if len(a.endOfTests) == len(a.students) {
			context.Respond(&command.FinishTest{Subject: msg.Subject})
		}
	}
}
