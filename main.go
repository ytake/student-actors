package main

import (
	"fmt"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
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
	cr := system.Root.Spawn(
		actor.PropsFromProducer(classroom.NewActor(th, students), actor.WithSupervisor(supervisor)))
	system.Root.Send(cr, &message.BeginClassRequest{Subject: "算数"})
	system.Root.Send(cr, &actor.Watch{Watcher: th})
	_, _ = console.ReadLine()
}
