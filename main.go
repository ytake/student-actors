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
	system.Root.ActorSystem()

	p := stream.NewTypedStream[*event.TestFinished](system)
	go func() {
		cr, _ := system.Root.SpawnNamed(
			actor.PropsFromProducer(
				classroom.NewActor(p.PID(), students())),
			"math-classroom")
		system.Root.Send(cr, &command.ClassStarts{Subject: "算数"})
	}()
	r := <-p.C()
	fmt.Printf("%s テストが終了しました\n", r.Subject)
}

func students() []int {
	var students []int
	for v := range [20]int{} {
		students = append(students, v+1)
	}
	return students
}
