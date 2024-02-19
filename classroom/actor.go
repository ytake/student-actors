package classroom

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/event"
	"github.com/ytake/student-actors/student"
	"github.com/ytake/student-actors/teacher"
)

// Actor represents a classroom.
type Actor struct {
	pipe     *actor.PID
	teacher  *actor.PID
	students []int
	mutex    sync.Mutex
}

func NewActor(pipe *actor.PID, students []int) func() actor.Actor {
	return func() actor.Actor {
		return &Actor{
			pipe:     pipe,
			students: students,
		}
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor.
func (class *Actor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *command.StartsClass:
		class.teacher = context.Spawn(actor.PropsFromProducer(teacher.NewActor))
		context.Request(class.teacher, &command.PrepareTest{Subject: msg.Subject, Students: class.students})
	case *command.StartTest:
		class.mutex.Lock()
		for _, st := range class.students {
			sta, err := context.SpawnNamed(
				actor.PropsFromProducer(student.NewActor),
				fmt.Sprintf("student-%d", st))
			if err != nil {
				context.Poison(context.Self())
			}
			context.Request(sta, msg)
		}
		class.mutex.Unlock()
	case *command.SubmitTest:
		// 注意 context.Forward(class.teacher)
		// forwardとの違いは、メッセージを受け取ったアクターがメッセージをそのまま転送することです。
		// 転送された場合、メッセージの送信元は変更されないため送信元の生徒に直接返信することができます。
		context.Request(class.teacher, msg)
	case *command.FinishTest:
		context.Send(class.pipe, &event.ClassFinished{Subject: msg.Subject})
		context.Poison(context.Self())
	}
}
