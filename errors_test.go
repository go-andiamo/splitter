package splitter

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrappedSplittingError(t *testing.T) {
	sErr := &splittingError{
		errorType: Wrapped,
		wrapped:   errors.New("whoops"),
	}
	require.Error(t, sErr)
	require.Equal(t, "whoops", sErr.Error())
}

func TestWrappedAsSplittingError(t *testing.T) {
	sErr := asSplittingError(errors.New("whoops"), 16)
	require.Error(t, sErr)
	require.Equal(t, "whoops", sErr.Error())

	sErr2 := asSplittingError(sErr, 0)
	require.Error(t, sErr2)
	require.Equal(t, sErr, sErr2)

	var err error
	sErr = asSplittingError(err, 0)
	require.NoError(t, sErr)
	require.Nil(t, sErr)

}

func TestNewSplittingError(t *testing.T) {
	err := newSplittingError(Unopened, 16, '(', Parenthesis)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(unopenedFmt, "(", 16), err.Error())

	err = newSplittingError(Unclosed, 16, ')', Parenthesis)
	require.Equal(t, fmt.Sprintf(unclosedFmt, ")", 16), err.Error())
}

func TestSplittingError_DefaultMessage(t *testing.T) {
	err := &splittingError{
		errorType: -1,
		message:   "whoops",
	}
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestNewOptionFailError(t *testing.T) {
	err := NewOptionFailError("whoops", 16, nil)
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
	sErr, ok := err.(*splittingError)
	require.True(t, ok)
	require.Equal(t, OptionFail, sErr.Type())
	require.Equal(t, 16, sErr.Position())

	subPart := &subPart{
		enc:      *Parenthesis,
		openPos:  5,
		closePos: 10,
	}
	err = NewOptionFailError("whoops", 0, subPart)
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
	sErr, ok = err.(*splittingError)
	require.True(t, ok)
	require.Equal(t, OptionFail, sErr.Type())
	require.Equal(t, 5, sErr.Position())
}

func TestSplitAlwaysReturnsSplittingError(t *testing.T) {
	s, err := NewSplitter(',', Parenthesis)
	require.NoError(t, err)
	s.AddDefaultOptions(NotEmptyFirstMsg("whoops at position %d"))

	_, err = s.Split(`(`)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(unclosedFmt, `(`, 0), err.Error())
	sErr, ok := err.(SplittingError)
	require.True(t, ok)
	require.Equal(t, Unclosed, sErr.Type())
	require.Equal(t, '(', sErr.Rune())
	require.Equal(t, 0, sErr.Position())
	require.Equal(t, Parenthesis, sErr.Enclosure())
	require.Nil(t, sErr.Wrapped())

	_, err = s.Split(`)`)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(unopenedFmt, `)`, 0), err.Error())
	sErr, ok = err.(SplittingError)
	require.True(t, ok)
	require.Equal(t, Unopened, sErr.Type())
	require.Equal(t, ')', sErr.Rune())
	require.Equal(t, 0, sErr.Position())
	require.Equal(t, Parenthesis, sErr.Enclosure())
	require.Nil(t, sErr.Wrapped())
	require.NoError(t, sErr.Unwrap())
	require.NoError(t, errors.Unwrap(sErr))

	_, err = s.Split(`,a`)
	require.Error(t, err)
	require.Equal(t, `whoops at position 0`, err.Error())
	sErr, ok = err.(SplittingError)
	require.True(t, ok)
	require.Equal(t, OptionFail, sErr.Type())
	require.Equal(t, rune(0), sErr.Rune())
	require.Equal(t, 0, sErr.Position())
	require.Nil(t, sErr.Enclosure())
	require.Nil(t, sErr.Wrapped())
	require.NoError(t, sErr.Unwrap())
	require.NoError(t, errors.Unwrap(sErr))

	_, err = s.Split(`a`, &errorOption{})
	require.Error(t, err)
	require.Equal(t, "error option", err.Error())
	sErr, ok = err.(SplittingError)
	require.True(t, ok)
	require.Equal(t, Wrapped, sErr.Type())
	require.Error(t, sErr.Wrapped())
	require.Error(t, sErr.Unwrap())
	require.Error(t, errors.Unwrap(sErr))
	require.Equal(t, rune(0), sErr.Rune())
	require.Nil(t, sErr.Enclosure())
	require.Equal(t, 0, sErr.Position())
}

type errorOption struct {
}

func (o *errorOption) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	return s, false, errors.New("error option")
}
