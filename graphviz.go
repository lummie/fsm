package fsm

import (
	"fmt"
	"io"
	"strconv"
)

// GVEncoder encodes a transition rule set to the GraphViz format which can then be used to generate a visualisation
type Encoder struct {
	startState State
	endStates  []State
	w          io.Writer
	Tags       map[State]string
}

// NewEncode returns a new Encoder to write in the Graphviz format.
// w is the io.Writer that the encoding will be written to
// startState is the state you'd like to mark as the initial state in the output. defaults to 0
// endStates is an array of state that you'd like to mark as termination states.
// stateTags is a map holding a string name for each state.
func NewEncoder(w io.Writer, startState State, endStates []State) *Encoder {
	return &Encoder{
		w:          w,
		startState: startState,
		endStates:  endStates,
		Tags:       make(map[State]string),
	}
}

func (enc *Encoder) getTag(s State) string {
	ss, ok := enc.Tags[s]
	if !ok {
		return "State " + strconv.Itoa(int(s))
	}
	return ss
}

func (enc *Encoder) writeDocStart() {
	fmt.Fprintln(enc.w, "digraph fsm {\n\trankdir=TB;\n\tsize=\"8,5\"")
}

func (enc *Encoder) writeDocEnd() {
	fmt.Fprintln(enc.w, "}")
}

func (enc *Encoder) writeStartStates() {
	fmt.Fprintln(enc.w, "\tnode [shape = circle, color=grey, style=filled];")
	fmt.Fprintln(enc.w, "\tSTATE_START [label=\"\"]")
}

func (enc *Encoder) writeTerminationStates() {
	if len(enc.endStates) > 0 {
		fmt.Fprintln(enc.w, "\tnode [shape = circle, color=grey, peripheries=2, style=filled];")
		fmt.Fprintln(enc.w, "\tSTATE_END [label=\"\"]")
	}
}

func (enc *Encoder) writeTransitions(t *Transitions) {
	fmt.Fprintln(enc.w, "\tnode [shape = circle, style=solid, peripheries=1];")
	fmt.Fprintf(enc.w, "\tSTATE_START -> %s [label=\"\"]\n", enc.getTag(enc.startState))
	for ss := range t.ts {
		for es, tr := range t.ts[ss] {
			fmt.Fprintf(enc.w, "\t%s -> %s [label=\"%s\"]\n", enc.getTag(ss), enc.getTag(es), tr.tag)
		}
	}
	for _, s := range enc.endStates {
		fmt.Fprintf(enc.w, "\t%s -> STATE_END [label=\"\"]\n", enc.getTag(s))
	}

}

func (enc *Encoder) Encode(t *Transitions) {
	enc.writeDocStart()
	enc.writeStartStates()
	enc.writeTerminationStates()
	enc.writeTransitions(t)
	enc.writeDocEnd()
}
