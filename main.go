package main

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/student-actors/classroom"
	"github.com/ytake/student-actors/command"
	"github.com/ytake/student-actors/event"
)

func main() {
	system := actor.NewActorSystem()
	p := stream.NewTypedStream[*event.TestFinished](system)
	cr, err := system.Root.SpawnNamed(
		actor.PropsFromProducer(
			classroom.NewActor(p.PID(), students())),
		"math-classroom")
	if err != nil {
		return
	}
	go func() {
		system.Root.Send(cr, &command.ClassStarts{Subject: "算数"})
	}()
	r := <-p.C()
	fmt.Printf("%s テストが終了しました\n", r.Subject)
}

func students() []int {
	sts := make([]int, 20)
	for i := 0; i < 20; i++ {
		sts[i] = i
	}
	return sts
}
