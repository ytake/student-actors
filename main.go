package main

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/student-actors/classroom"
	"github.com/ytake/student-actors/message"
	"github.com/ytake/student-actors/teacher"
)

func main() {
	system := actor.NewActorSystem()
	system.Root.ActorSystem()

	th := system.Root.Spawn(actor.PropsFromProducer(teacher.NewActor))

	p := stream.NewTypedStream[*message.ReceivedAchievementTest](system)
	go func() {
		cr, _ := system.Root.SpawnNamed(
			actor.PropsFromProducer(
				classroom.NewActor(p.PID(), th, students())),
			"math-classroom")
		system.Root.Send(cr, &message.BeginClassRequest{Subject: "算数"})
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
