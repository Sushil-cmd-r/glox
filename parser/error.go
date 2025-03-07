package parser

import "fmt"

type Error struct {
	Msg string
}

func (e Error) Error() string {
	return e.Msg
}

type ErrorList []Error

func (e ErrorList) Error() string {
	switch e.Len() {
	case 0:
		return "no errors"
	case 1:
		return e[0].Error()

	}
	return fmt.Sprintf("%s and (%d more errors)", e[0], len(e)-1)
}

func (e *ErrorList) Add(msg string) {
	*e = append(*e, Error{Msg: msg})
}

func (e ErrorList) Len() int {
	return len(e)
}

func (e ErrorList) Err() error {
	if len(e) == 0 {
		return nil
	}

	return e
}
