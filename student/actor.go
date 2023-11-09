package student

import "github.com/asynkron/protoactor-go/actor"

type Actor struct {
}

func NewActor() actor.Actor {
	return &Actor{}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (a *Actor) Receive(ctx actor.Context) {

}

var studentList = []string{
	"John",
}

func (a *Actor) GetStudentList(ctx actor.Context) {
	ctx.Respond(studentList)
}
