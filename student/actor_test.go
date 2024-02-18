package student

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
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
			msg:     &event.TestStarted{Subject: "math"},
			want:    &event.TestSubmitted{Subject: "math", Name: "student"},
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
	pid, _ := system.Root.SpawnNamed(actor.PropsFromProducer(NewActor), "student")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := system.Root.RequestFuture(pid, tt.msg, 10*time.Second)
			r, err := f.Result()
			if tt.isError {
				if !errors.Is(err, actor.ErrDeadLetter) {
					t.Error(err)
				}
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("got: %v, want: %v", r, tt.want)
			}
		})
	}
}
