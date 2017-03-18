package state_test

import (
	. "github.com/lummie/assert"
	"github.com/lummie/state"
	"testing"
)

const (
	// states
	Off = iota
	On

	// transitions
	TurnOn  = iota
	TurnOff
)

func TestNewStateMachine(t *testing.T) {
	m := state.NewStateMachine(Off)
	Assert(t, m.State(), Equal, Off, "Expected inital state of Off")
}

func TestStateMachine_AddState(t *testing.T) {
	m := state.NewStateMachine(Off)
	m.AddState(Off).AddTransition(TurnOn, On)
	m.AddState(On).AddTransition(TurnOff, Off)

}
