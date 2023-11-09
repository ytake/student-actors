package strategy

import (
	"github.com/asynkron/protoactor-go/actor"
)

// NewDecider returns a new Supervisor strategy which applies the fault Directive from the decider
func NewDecider() func(reason interface{}) actor.Directive {
	return func(reason interface{}) actor.Directive {
		switch reason.(type) {
		case *actor.Restarting:
			return actor.RestartDirective
		case *actor.Stopping:
			return actor.StopDirective
		case *actor.Stopped:
			return actor.StopDirective
		case *actor.Started:
			return actor.RestartDirective
		default:
			return actor.EscalateDirective
		}
	}
}
