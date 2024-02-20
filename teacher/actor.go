package teacher

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/student"
)

type Actor struct {
	students   []int
	endOfTests []command.SubmitTest
	replyTo    *actor.PID
	mutex      sync.Mutex
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
	case *command.PrepareTest:
		// 先生が宿題を出す
		context.Logger().Info("先生が", msg.Subject, "テストを出しました")
		a.mutex.Lock()
		for _, st := range a.students {
			sta, err := context.SpawnNamed(
				actor.PropsFromProducer(student.NewActor),
				fmt.Sprintf("student-%d", st))
			if err != nil {
				context.Poison(context.Self())
			}
			context.Send(sta, &command.StartTest{Subject: msg.Subject})
		}
		a.mutex.Unlock()
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
