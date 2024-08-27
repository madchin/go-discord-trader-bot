package offer

import "errors"

var ErrActionTypeNotDefined = errors.New("action type is neither Add / Update / Remove / List")

type action struct {
	a string
}

var (
	Add    = action{"add"}
	Update = action{"update"}
	Remove = action{"remove"}
	List   = action{"list"}
)

func ActionTypeFromCandidate(candidate string) (action, error) {
	if Add.a == candidate {
		return Add, nil
	}
	if Update.a == candidate {
		return Update, nil
	}
	if List.a == candidate {
		return List, nil
	}
	if Remove.a == candidate {
		return Remove, nil
	}
	return action{}, ErrOfferTypeNotDefined
}

func (a action) Action() string {
	return a.a
}
