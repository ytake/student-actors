package teacher

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/student"
)

type Actor struct {
	students   []int
	endOfTests []command.SubmitTest
	replyTo    *actor.PID
}

func NewActor(students []int, replyTo *actor.PID) actor.Actor {
	return &Actor{
		students:   students,
		replyTo:    replyTo,
		endOfTests: []command.SubmitTest{},
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Restarting:
		// リスタート時に実行される
		context.Send(context.Self(), &command.PrepareTest{Subject: "math"})
	case *command.PrepareTest:
		context.Logger().Info("先生が", msg.Subject, "テストを出しました")
		for _, st := range a.students {
			sta, err := context.SpawnNamed(
				actor.PropsFromProducer(student.NewActor),
				fmt.Sprintf("student-%d", st))
			if err != nil {
				context.Logger().Error(fmt.Sprintf("生徒 %d 生成できませんでした", st))
				panic(err)
			}
			context.Send(sta, &command.StartTest{Subject: msg.Subject})
		}
		// 生徒がテストを提出する
	case *command.SubmitTest:
		context.Logger().Info(
			fmt.Sprintf("先生が %s の %s テストの解答を受け取りました", msg.Name, msg.Subject))
		a.endOfTests = append(a.endOfTests, *msg)
		// 全員提出したら先生がテストの解答を受け取る
		if len(a.endOfTests) == len(a.students) {
			context.Send(a.replyTo, &command.FinishTest{Subject: msg.Subject})
		}
	}
}
