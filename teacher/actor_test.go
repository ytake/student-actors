package teacher

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/event"
)

func TestActor_Receive(t *testing.T) {
	tests := []struct {
		name    string
		msg     interface{}
		want    interface{}
		isError bool
	}{
		{
			name:    "test start",
			msg:     &command.ClassStarts{Subject: "math"},
			want:    &event.TestStarted{Subject: "math"},
			isError: false,
		},
		{
			name:    "plain string",
			msg:     "hello",
			want:    nil,
			isError: true,
		},
	}
	system := actor.NewActorSystem()
	pid, _ := system.Root.SpawnNamed(actor.PropsFromProducer(NewActor), "teacher")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := system.Root.RequestFuture(pid, tt.msg, 3*time.Second)
			r, err := f.Result()
			if tt.isError {
				if !errors.Is(err, actor.ErrTimeout) {
					t.Error(err)
				}
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("got: %v, want: %v", r, tt.want)
			}
		})
	}
}

func TestActor_Receive_panic(t *testing.T) {
	system := actor.NewActorSystem()
	pid, _ := system.Root.SpawnNamed(actor.PropsFromProducer(NewActor), "teacher")
	// 授業開始時に先生が宿題を出す
	// その後に先生アクターを意図的にパニックが発生する
	system.Root.Send(pid, &command.ClassStarts{Subject: "math"})
	// 先生アクターがパニックした後に復活し、テストの解答を受け取る
	f := system.Root.RequestFuture(pid, &event.TestFinished{Subject: "math"}, 3*time.Second)
	r, err := f.Result()
	if err != nil {
		t.Error(err)
	}
	// テストの解答を受け取っていることを確認
	if !reflect.DeepEqual(r, &event.TestReceived{Subject: "math"}) {
		t.Errorf("got: %v, want: %v", r, &event.TestReceived{Subject: "math"})
	}
}
