package customErrors

import (
	"errors"
	"fmt"
)

type EndOfInteraction struct {
}

var EOI = errors.New("EOI")

func (e *EndOfInteraction) Error() string {
	return fmt.Sprintf("EOI")
}
