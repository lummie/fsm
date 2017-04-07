// Provides an implementation of a Finite State Machine

package fsm

import (
	"errors"
	"fmt"
)

type State int

type Guard func(current, target State, ctx interface{}) error
type Action func(state State, ctx interface{})

const eInvalidStateChange = "Invalid state change %s -> %s"

// transition holds the guard and action functions that result in the change of a state to a new state `s`
// if a tag is specified this will be used to identify the transition in any encodings.
type transition struct {
	g   Guard
	a   Action
	tag string
}

// Transitions provides the set of rules for changing states
type Transitions struct {
	ts map[State]map[State]transition
}

// NewTransitions creates an instance of Transitions, providing a set of rules for state change.
func NewTransitions() *Transitions {
	t := Transitions{
		ts: make(map[State]map[State]transition),
	}
	return &t
}

// Add adds a new state transition rule to the transitions list
func (t *Transitions) Add(from, to State, guard Guard, action Action, tag string) {
	_, ok := t.ts[from]
	if !ok {
		t.ts[from] = make(map[State]transition, 0)
	}
	t.ts[from][to] = transition{
		g:   guard,
		a:   action,
		tag: tag,
	}
}

// Get retrieves the information about a specific transition rule.
func (t *Transitions) Get(from, to State) (ok bool, guard Guard, action Action) {
	tos, ok := t.ts[from]
	if !ok {
		return
	}

	info, ok := tos[to]
	if !ok {
		return
	}
	return true, info.g, info.a
}

// StateMachine is instance that of the FSM that holds the state and allows changes of state to be requested.
// The StateMachine does not have any transition rules directly associated with it. When you request a change,
// you supply the transition rules along with the requested state. This allows you to model different state rules and
// choose when to apply them to the instance of the StateMachine.
type StateMachine struct {
	current State
}

// NewStateMachine returns a new instance of a FSM
func NewStateMachine(initial State) *StateMachine {
	return &StateMachine{current: initial}
}

// Change, requests a change in state based on the supplied transition rules and the newState requested.
// Additionally a context object may be supplied that will passed on to the guard and action functions for the transition.
// if the guard returns an error this is passed back by Change, in addition an attempt to change to an invalid state
// will cause an non nil error to be returned.

func (sm *StateMachine) Change(rules *Transitions, newState State, context interface{}) error {
	ok, g, a := rules.Get(sm.current, newState)
	if !ok {
		return errors.New(fmt.Sprintf(eInvalidStateChange, sm.current, newState))
	}
	if g != nil {
		if err := g(sm.current, newState, context); err != nil {
			return err
		}
	}
	sm.current = newState
	if a != nil {
		a(sm.current, context)
	}
	return nil
}
