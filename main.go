package main

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/ytake/student-actors/classroom"
	"github.com/ytake/student-actors/message"
	"github.com/ytake/student-actors/strategy"
	"github.com/ytake/student-actors/student"
	"github.com/ytake/student-actors/teacher"
)

func main() {
	system := actor.NewActorSystem()
	system.Root.ActorSystem()

	supervisor := actor.NewOneForOneStrategy(10, 1000, strategy.NewDecider())
	th := system.Root.Spawn(actor.PropsFromProducer(teacher.NewActor))
	studentProps := actor.PropsFromProducer(student.NewActor)
	students := make([]*actor.PID, 0)
	for v := range [20]int{} {
		st, _ := system.Root.SpawnNamed(studentProps, fmt.Sprintf("student-%d", v))
		students = append(students, st)
	}
	p := stream.NewTypedStream[*message.EndOfAchievementTest](system)
	go func() {
		cr := system.Root.Spawn(
			actor.PropsFromProducer(classroom.NewActor(p.PID(), th, students), actor.WithSupervisor(supervisor)))
		system.Root.Send(cr, &message.BeginClassRequest{Subject: "算数"})
		system.Root.Send(cr, &actor.Watch{Watcher: th})
	}()
	r := <-p.C()
	fmt.Println(fmt.Sprintln("全員が", r.Subject, "テストの解答を提出しました"))
}
