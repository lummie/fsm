package fsm_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lummie/assert"
	"github.com/lummie/fsm"
	"testing"
)

const (
	// states
	Red      fsm.State = iota
	RedAmber
	Green
	Amber
	Blue     = 101
)

func TestSimpleTransitions(t *testing.T) {
	tr := fsm.NewTransitions()
	assert.NotNil(t, tr, "Transitions hould have been created")

	tr.Add(Red, RedAmber, nil, nil, "")

	ok, g, a := tr.Get(Red, Green)
	assert.False(t, ok, "Red, Green not expected")
	assert.Nil(t, g)
	assert.Nil(t, a)

	ok, g, a = tr.Get(Red, RedAmber)
	assert.True(t, ok, "Expected Red, RedAmber")
	assert.Nil(t, g)
	assert.Nil(t, a)
}

func TestTransitionsThatAreInvalid(t *testing.T) {
	tr := fsm.NewTransitions()
	assert.NotNil(t, tr, "Transitions hould have been created")

	tr.Add(Red, RedAmber, nil, nil, "")
	tr.Add(Red, Green, nil, nil, "")

	ok, _, _ := tr.Get(Red, Blue)
	assert.False(t, ok, "Blue not a valid transtion to")
	ok, _, _ = tr.Get(Blue, Red)
	assert.False(t, ok, "Blue not a valid transtion from")
	ok, _, _ = tr.Get(RedAmber, Green)
	assert.False(t, ok, "Not a valid transtion")

}

func TestTransitionsWithGuardAndAction(t *testing.T) {
	tr := fsm.NewTransitions()
	assert.NotNil(t, tr, "Transitions hould have been created")

	var gc int
	myGuard := func(current, target fsm.State, ctx interface{}) error {
		gc++
		return nil
	}

	var ac int
	myAction := func(state fsm.State, ctx interface{}) {
		ac++
	}

	tr.Add(Red, RedAmber, myGuard, myAction, "")

	ok, g, a := tr.Get(Red, RedAmber)
	assert.True(t, ok, "Red, RedAmber expected")
	assert.NotNil(t, g)
	assert.NotNil(t, a)
	g(Red, RedAmber, nil)
	assert.Equal(t, 1, gc, "expected gc to have been incremented")

	a(RedAmber, nil)
	assert.Equal(t, 1, ac, "expected gc to have been incremented")
}

func TestTransitions_MultipleOnSameState(t *testing.T) {
	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, nil, nil, "")
	tr.Add(Red, Green, nil, nil, "")
	tr.Add(Red, Amber, nil, nil, "")

	ok, _, _ := tr.Get(Red, RedAmber)
	assert.True(t, ok, "Red, RedAmber expected")

	ok, _, _ = tr.Get(Red, Green)
	assert.True(t, ok, "Red, Green expected")

	ok, _, _ = tr.Get(Red, Amber)
	assert.True(t, ok, "Red, Amber expected")
}

func TestTransitions_MultipleDifferentStates(t *testing.T) {
	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, nil, nil, "")
	tr.Add(RedAmber, Green, nil, nil, "")
	tr.Add(Green, Amber, nil, nil, "")
	tr.Add(Amber, Red, nil, nil, "")

	ok, _, _ := tr.Get(Red, RedAmber)
	assert.True(t, ok, "Red, RedAmber expected")

	ok, _, _ = tr.Get(RedAmber, Green)
	assert.True(t, ok, "RedAmber, Green expected")

	ok, _, _ = tr.Get(Green, Amber)
	assert.True(t, ok, "Green, Amber expected")

	ok, _, _ = tr.Get(Amber, Red)
	assert.True(t, ok, "Amber, Red expected")
}

func TestNewStateMachine(t *testing.T) {
	sm := fsm.NewStateMachine(Red)

	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, nil, nil, "")
	tr.Add(RedAmber, Green, nil, nil, "")
	tr.Add(Green, Amber, nil, nil, "")
	tr.Add(Amber, Red, nil, nil, "")

	err := sm.Change(tr, RedAmber, nil)
	assert.NotError(t, err)
	err = sm.Change(tr, Green, nil)
	assert.NotError(t, err)
	err = sm.Change(tr, Amber, nil)
	assert.NotError(t, err)
	err = sm.Change(tr, Red, nil)
	assert.NotError(t, err)

}

func TestNewStateMachineNegative(t *testing.T) {
	sm := fsm.NewStateMachine(Red)

	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, nil, nil, "")
	tr.Add(RedAmber, Green, nil, nil, "")
	tr.Add(Green, Amber, nil, nil, "")
	tr.Add(Amber, Red, nil, nil, "")

	err := sm.Change(tr, Green, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Amber, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Red, nil)
	assert.Error(t, err)
	err = sm.Change(tr, RedAmber, nil)
	assert.NotError(t, err)

	err = sm.Change(tr, RedAmber, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Amber, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Red, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Green, nil)
	assert.NotError(t, err)

	err = sm.Change(tr, RedAmber, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Green, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Red, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Amber, nil)
	assert.NotError(t, err)

	err = sm.Change(tr, RedAmber, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Green, nil)
	assert.Error(t, err)
	err = sm.Change(tr, Amber, nil)
	assert.Error(t, err)

}

func TestTransitionsWithGuardFailureAndAction(t *testing.T) {
	sm := fsm.NewStateMachine(Red)

	var gc int
	myGuard := func(current, target fsm.State, ctx interface{}) error {
		gc++
		return errors.New("Guard failed")
	}

	var ac int
	myAction := func(state fsm.State, ctx interface{}) {
		ac++
	}

	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, myGuard, myAction, "")
	tr.Add(Red, Green, nil, myAction, "")

	err := sm.Change(tr, RedAmber, nil)
	assert.Error(t, err, "Red, RedAmber expected")
	assert.Equal(t, 1, gc, "expected gc to have been incremented")
	assert.Equal(t, 0, ac, "expected ac to have not changed")

	err = sm.Change(tr, Green, nil)
	assert.NotError(t, err, "Red -> Green expected")

	assert.Equal(t, 1, ac, "expected ac to have been incremented")
}

func ExampleNewStateMachine() {
	// Traffic light example
	sm := fsm.NewStateMachine(Red)

	// a guard that always succeeds and prints the state change request
	allowChangeGuard := func(current, target fsm.State, ctx interface{}) error {
		fmt.Printf("Guard: %d -> %d.\n", current, target)
		return nil
	}

	// an action that simply prints the new state
	printStateAction := func(state fsm.State, ctx interface{}) {
		fmt.Printf("Action: %d.\n", state)
	}

	// create the transition rules
	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, allowChangeGuard, printStateAction, "Prepare to Go")
	tr.Add(RedAmber, Green, allowChangeGuard, printStateAction, "Go")
	tr.Add(Green, Amber, allowChangeGuard, printStateAction, "Stop unless unsafe to do so")
	tr.Add(Amber, Red, allowChangeGuard, printStateAction, "Stop")

	// change the state of the machine
	sm.Change(tr, RedAmber, nil)
	sm.Change(tr, Green, nil)
	sm.Change(tr, Amber, nil)
	sm.Change(tr, Red, nil)

	// Output:
	// Guard: 0 -> 1.
	// Action: 1.
	// Guard: 1 -> 2.
	// Action: 2.
	// Guard: 2 -> 3.
	// Action: 3.
	// Guard: 3 -> 0.
	// Action: 0.
}

func TestStateMachine_InvalidState(t *testing.T) {
	sm := fsm.NewStateMachine(Red)
	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, nil, nil, "")
	tr.Add(RedAmber, Green, nil, nil, "")
	tr.Add(Green, Amber, nil, nil, "")
	tr.Add(Amber, Red, nil, nil, "")

	err := sm.Change(tr, fsm.State(99), nil)
	assert.Error(t, err, "Expected an error for an uknown state")
}

func ExampleEncoder_Encode() {
	// create the transition rules
	tr := fsm.NewTransitions()
	tr.Add(Red, RedAmber, nil, nil, "Prepare to Go")
	tr.Add(RedAmber, Green, nil, nil, "Go")
	tr.Add(Green, Amber, nil, nil, "Stop unless unsafe to do so")
	tr.Add(Amber, Red, nil, nil, "Stop")

	buf := bytes.NewBuffer([]byte{})

	e := fsm.NewEncoder(buf, Red, []fsm.State{Green})
	e.Tags[Red] = "Red"
	e.Tags[RedAmber] = "RedAmber"
	e.Tags[Green] = "Green"
	e.Tags[Amber] = "Amber"
	e.Encode(tr)

	fmt.Print(string(buf.Bytes()))

	// Outputs something like:
	//digraph fsm {
	//	rankdir=TB;
	//	size="8,5"
	//	node [shape = circle, color=grey, style=filled];
	//	STATE_START [label=""]
	//	node [shape = circle, color=grey, peripheries=2, style=filled];
	//	STATE_END [label=""]
	//	node [shape = circle, style=empty, peripheries=1];
	//	STATE_START -> Red [label=""]
	//	Green -> Amber [label="Stop unless unsafe to do so"]
	//	Amber -> Red [label="Stop"]
	//	Red -> RedAmber [label="Prepare to Go"]
	//	RedAmber -> Green [label="Go"]
	//	Green -> STATE_END [label=""]
	//}
}
