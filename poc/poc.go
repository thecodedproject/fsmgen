package poc

import (
)

// to be defined in fsm pkg
type Transition interface{}

type Insert_StateOne struct {
	Transition

	ValueOne string
}
