package teacher

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/student-actors/command"
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
			msg:     &command.PrepareTest{Subject: "math", Students: []int{1}},
			want:    &command.StartTest{Subject: "math"},
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

func TestActor_Respond_FinishTest(t *testing.T) {
	system := actor.NewActorSystem()
	pid, _ := system.Root.SpawnNamed(actor.PropsFromProducer(NewActor), "teacher")
	// 授業開始時に先生が宿題を出す
	system.Root.Send(pid, &command.PrepareTest{Subject: "math", Students: []int{1}})
	f := system.Root.RequestFuture(pid, &command.SubmitTest{Subject: "math", Name: "ytake"}, 3*time.Second)
	r, err := f.Result()
	if err != nil {
		t.Error(err)
	}
	// テストの解答を受け取っていることを確認
	expect := &command.FinishTest{Subject: "math"}
	if !reflect.DeepEqual(r, expect) {
		t.Errorf("got: %v, want: %v", r, expect)
	}
}
