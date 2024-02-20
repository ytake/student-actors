package teacher

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/student-actors/command"
)

func TestActor_Receive(t *testing.T) {
	t.Run("plain string", func(t *testing.T) {
		system := actor.NewActorSystem()
		forwarder := system.Root.Spawn(actor.PropsFromFunc(func(context actor.Context) {
			switch msg := context.Message(); msg {
			case "hello":
				context.Respond(msg)
			}
		}))
		pid, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
			return NewActor([]int{1}, forwarder)
		}), "teacher")
		f := system.Root.RequestFuture(pid, "plain string", 10*time.Second)
		r, err := f.Result()
		if !errors.Is(err, actor.ErrTimeout) {
			t.Error(err)
		}
		if !reflect.DeepEqual(r, nil) {
			t.Errorf("got: %v, want: %v", r, nil)
		}
	})
	t.Run("receive FinishTest", func(t *testing.T) {
		system := actor.NewActorSystem()
		p := stream.NewTypedStream[*command.FinishTest](system)
		pid, _ := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
			return NewActor([]int{1}, p.PID())
		}), "teacher")
		// 授業開始時に先生が宿題を出す
		system.Root.Send(pid, &command.PrepareTest{Subject: "math"})
		system.Root.Send(pid, &command.SubmitTest{Subject: "math", Name: "ytake"})
		r := <-p.C()
		// テストの解答を受け取っていることを確認
		expect := &command.FinishTest{Subject: "math"}
		if !reflect.DeepEqual(r, expect) {
			t.Errorf("got: %v, want: %v", r, expect)
		}
	})
}
