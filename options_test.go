package splitter

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitOptionsMerge(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)
	const testStr = `/a/b/c/`
	pts, err := s.Split(testStr)
	require.NoError(t, err)
	require.Equal(t, 5, len(pts))

	pts, err = s.Split(testStr, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	s.AddDefaultOptions(IgnoreEmpties)
	pts, err = s.Split(testStr)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	pts, err = s.Split(testStr, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	pts, err = s.Split(testStr, IgnoreEmptyFirst)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))
}

func TestOption_TrimSpaces(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(` / / `, TrimSpaces)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))
	require.Equal(t, ``, pts[0])
	require.Equal(t, ``, pts[1])
	require.Equal(t, ``, pts[2])

	pts, err = s.Split(` / / `, TrimSpaces, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 1, len(pts))
	require.Equal(t, ``, pts[0])

	pts, err = s.Split(` / /`, TrimSpaces, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 1, len(pts))
	require.Equal(t, ``, pts[0])

	pts, err = s.Split(` / `, TrimSpaces, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))

	pts, err = s.Split(` `, TrimSpaces, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))

	pts, err = s.Split(``, TrimSpaces, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))
}

func TestOption_Trim(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	trimmer := Trim(" \t\n")
	pts, err := s.Split("\t/\n/ ", trimmer)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))
	require.Equal(t, ``, pts[0])
	require.Equal(t, ``, pts[1])
	require.Equal(t, ``, pts[2])

	pts, err = s.Split("\t/\n/ ", trimmer, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 1, len(pts))
	require.Equal(t, ``, pts[0])

	pts, err = s.Split("\t/\n/", trimmer, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 1, len(pts))
	require.Equal(t, ``, pts[0])

	pts, err = s.Split("\t /\n ", trimmer, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))

	pts, err = s.Split("\t\n ", trimmer, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))

	pts, err = s.Split(``, trimmer, IgnoreEmptyFirst, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))
}

func TestOption_NoEmpties(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(`a/b/c`, NoEmpties)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	_, err = s.Split(`a//c`, NoEmpties)
	require.Error(t, err)
	require.Equal(t, _NoEmpties.message, err.Error())
}

func TestOption_NoEmptiesMsg(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	_, err = s.Split(`a//c`, NoEmptiesMsg("whoops"))
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestOption_IgnoreEmpties(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(`a/b/c`, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	pts, err = s.Split(`/a/b/c`, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	pts, err = s.Split(`//b/c`, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 2, len(pts))

	pts, err = s.Split(`///c`, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 1, len(pts))

	pts, err = s.Split(`///`, IgnoreEmpties)
	require.NoError(t, err)
	require.Equal(t, 0, len(pts))
}

func TestOption_NotEmptyFirst(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(`a/b/c`, NotEmptyFirst)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	_, err = s.Split(`/a/b/c`, NotEmptyFirst)
	require.Error(t, err)
	require.Equal(t, _NotEmptyFirst.message, err.Error())
}

func TestOption_NotEmptyFirstMsg(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	_, err = s.Split(`/a/b/c`, NotEmptyFirstMsg("whoops"))
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestOption_IgnoreEmptyFirst(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(`a/b/c`, IgnoreEmptyFirst)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	pts, err = s.Split(`/a/b/c`, IgnoreEmptyFirst)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))
}

func TestOption_NotEmptyLast(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(`a/b/c`, NotEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	_, err = s.Split(`a/b/c/`, NotEmptyLast)
	require.Error(t, err)
	require.Equal(t, _NotEmptyLast.message, err.Error())
}

func TestOption_NotEmptyLastMsg(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	_, err = s.Split(`a/b/c/`, NotEmptyLastMsg("whoops"))
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestOption_IgnoreEmptyLast(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)

	pts, err := s.Split(`a/b/c`, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))

	pts, err = s.Split(`a/b/c/`, IgnoreEmptyLast)
	require.NoError(t, err)
	require.Equal(t, 3, len(pts))
}

func TestOption_NoContiguousQuotes(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes)
	require.NoError(t, err)

	_, err = s.Split(`a/"b" "b"/c`, NoContiguousQuotes)
	require.NoError(t, err)

	_, err = s.Split(`a/"b""b"/c/`, NoContiguousQuotes)
	require.Error(t, err)
	require.Equal(t, _NoContiguousQuotes.message, err.Error())
}

func TestOption_NoContiguousQuotesMsg(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes)
	require.NoError(t, err)

	_, err = s.Split(`a/"b""b"/c/`, NoContiguousQuotesMsg("whoops"))
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestOption_NoMultiQuotes(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes)
	require.NoError(t, err)

	_, err = s.Split(`a/"b"/c`, NoMultiQuotes)
	require.NoError(t, err)

	_, err = s.Split(`a/"b""b"/c`, NoMultiQuotes)
	require.Error(t, err)
	require.Equal(t, _NoMultiQuotes.message, err.Error())
}

func TestOption_NoMultiQuotesMsg(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes)
	require.NoError(t, err)

	_, err = s.Split(`a/"b""b"/c`, NoMultiQuotesMsg("whoops"))
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestOption_NoMultis(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes, Parenthesis)
	require.NoError(t, err)

	_, err = s.Split(`a/"b"/(,)`, NoMultis)
	require.NoError(t, err)

	_, err = s.Split(`a/"b" ()/(,)`, NoMultis)
	require.Error(t, err)
	require.Equal(t, _NoMultis.message, err.Error())
}

func TestOption_NoMultisMsg(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes, Parenthesis)
	require.NoError(t, err)

	_, err = s.Split(`a/"b" ()/(,)`, NoMultisMsg("whoops"))
	require.Error(t, err)
	require.Equal(t, "whoops", err.Error())
}

func TestOption_StripQuotes(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotes)
	require.NoError(t, err)

	pts, err := s.Split(`a/"b"`, StripQuotes)
	require.NoError(t, err)
	require.Equal(t, 2, len(pts))
	require.Equal(t, `a`, pts[0])
	require.Equal(t, `b`, pts[1])

	pts, err = s.Split(`a/"b""b"`, StripQuotes)
	require.NoError(t, err)
	require.Equal(t, 2, len(pts))
	require.Equal(t, `a`, pts[0])
	require.Equal(t, `bb`, pts[1])

	pts, err = s.Split(`a/"b" "b"`, StripQuotes)
	require.NoError(t, err)
	require.Equal(t, 2, len(pts))
	require.Equal(t, `a`, pts[0])
	require.Equal(t, `b b`, pts[1])

	pts, err = s.Split(`"a"/"b""b""b"`, StripQuotes)
	require.NoError(t, err)
	require.Equal(t, 2, len(pts))
	require.Equal(t, `a`, pts[0])
	require.Equal(t, `bbb`, pts[1])
}

func TestOption_UnescapeQuotes(t *testing.T) {
	s, err := NewSplitter('/', DoubleQuotesDoubleEscaped)
	require.NoError(t, err)

	pts, err := s.Split(`"a"""/"b""b""b"`)
	require.NoError(t, err)
	require.Equal(t, 2, len(pts))
	require.Equal(t, `"a"""`, pts[0])
	require.Equal(t, `"b""b""b"`, pts[1])

	pts, err = s.Split(`"a"""/"b""b""b" "b"`, UnescapeQuotes)
	require.Equal(t, 2, len(pts))
	require.Equal(t, `a"`, pts[0])
	require.Equal(t, `b"b"b b`, pts[1])

	pts, err = s.Split(`a`, UnescapeQuotes)
	require.NoError(t, err)
	require.Equal(t, 1, len(pts))
	require.Equal(t, `a`, pts[0])
}
