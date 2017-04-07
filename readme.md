# fsm - A Finite State Machine in Go

As with MVC, the persistence of the state machine (StateMachine)
 is separate from the controller of the state machine (Transitions)

Here is a traffic light example.

```go
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
```

In addition the transitions that you define can be encoded to GraphViz input by using an Encoder.

```go
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

```

when run through dot will create a nice state diagram like so..

![Traffic Light State Diagram](https://github.com/lummie/fsm/raw/master/example.gv.png)

btw. I know there is no final state for a traffic light, it is just there to demonstrate the terminators.