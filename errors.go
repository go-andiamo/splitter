package splitter

import (
	"fmt"
	"strings"
)

// SplittingErrorType is the splitting error type - as used by splittingError
type SplittingErrorType int

const (
	Unopened SplittingErrorType = iota
	Unclosed
	OptionFail
	Wrapped
)

// SplittingError is the error type always returned from Splitter.Split
type SplittingError interface {
	error
	Type() SplittingErrorType
	Position() int
	Rune() rune
	Enclosure() *Enclosure
	Wrapped() error
	Unwrap() error
}

type splittingError struct {
	errorType SplittingErrorType
	position  int
	rune      rune
	enc       *Enclosure
	wrapped   error
	message   string
}

func newSplittingError(t SplittingErrorType, pos int, r rune, enc *Enclosure) SplittingError {
	return &splittingError{
		errorType: t,
		position:  pos,
		rune:      r,
		enc:       enc,
	}
}

const (
	unopenedFmt = "unopened '%s' at position %d"
	unclosedFmt = "unclosed '%s' at position %d"
)

func (e *splittingError) Error() string {
	if e.errorType == Unopened {
		return fmt.Sprintf(unopenedFmt, string(e.rune), e.position)
	} else if e.errorType == Unclosed {
		return fmt.Sprintf(unclosedFmt, string(e.rune), e.position)
	} else if e.wrapped != nil {
		return e.wrapped.Error()
	}
	result := fmt.Sprintf(e.message, e.position)
	if strings.HasSuffix(result, fmt.Sprintf(`%%!(EXTRA int=%d)`, e.position)) {
		result = result[:strings.LastIndex(result, "%!(EXTRA int=")]
	}
	return result
}

func (e *splittingError) Unwrap() error {
	return e.wrapped
}

func (e *splittingError) Type() SplittingErrorType {
	return e.errorType
}
func (e *splittingError) Position() int {
	return e.position
}
func (e *splittingError) Rune() rune {
	return e.rune
}
func (e *splittingError) Enclosure() *Enclosure {
	return e.enc
}
func (e *splittingError) Wrapped() error {
	return e.wrapped
}

func asSplittingError(err error, pos int) SplittingError {
	if err != nil {
		if se, ok := err.(SplittingError); ok {
			return se
		} else {
			return &splittingError{
				errorType: Wrapped,
				position:  pos,
				wrapped:   err,
			}
		}
	}
	return nil
}

func NewOptionFailError(msg string, pos int, subPart SubPart) SplittingError {
	if subPart != nil {
		return &splittingError{
			errorType: OptionFail,
			position:  subPart.StartPos(),
			rune:      subPart.StartRune(),
			enc:       subPart.Enclosure(),
			message:   msg,
		}
	}
	return &splittingError{
		errorType: OptionFail,
		position:  pos,
		message:   msg,
	}
}
