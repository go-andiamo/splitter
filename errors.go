package splitter

import "fmt"

type SplittingErrorType int

const (
	Unopened SplittingErrorType = iota
	Unclosed
)

type SplittingError struct {
	Type      SplittingErrorType
	Position  int
	Rune      rune
	Enclosure *Enclosure
}

func newSplittingError(t SplittingErrorType, pos int, r rune, enc *Enclosure) error {
	return &SplittingError{
		Type:      t,
		Position:  pos,
		Rune:      r,
		Enclosure: enc,
	}
}

const (
	unopenedFmt = "unopened '%s' at position %d"
	unclosedFmt = "unclosed '%s' at position %d"
)

func (e *SplittingError) Error() string {
	if e.Type == Unopened {
		return fmt.Sprintf(unopenedFmt, string(e.Rune), e.Position)
	}
	return fmt.Sprintf(unclosedFmt, string(e.Rune), e.Position)
}
