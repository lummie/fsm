package state

type State struct {
	transitions map[int]int
}

func (s *State) AddTransition(transition, newState int) *State {
	s.transitions[transition] = newState
	return s
}

type StateMachine struct {
	states    map[int]State
	current   int
	validated bool
}

func NewStateMachine(initialState int) *StateMachine {
	m := StateMachine{
		states:  make(map[int]State, 0),
		current: initialState,
	}
	return &m
}

func (m *StateMachine) AddState(state int) *State {
	s := State{
		transitions: make(map[int]int, 0),
	}
	m.states[state] = s
	return &s
}

func (m *StateMachine) State() int {
	return m.current
}
