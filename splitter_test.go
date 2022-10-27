package splitter

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestNewSplitter(t *testing.T) {
	s, err := NewSplitter('/')
	require.NoError(t, err)
	require.NotNil(t, s)
	rs, ok := s.(*splitter)
	require.True(t, ok)
	require.NotNil(t, rs)
	require.Equal(t, '/', rs.separator)
	require.Equal(t, 0, len(rs.enclosures))

	enc := &Enclosure{
		Start: '{',
		End:   '}',
	}
	s, err = NewSplitter('/', enc, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	rs, ok = s.(*splitter)
	require.True(t, ok)
	require.NotNil(t, rs)
	require.Equal(t, '/', rs.separator)
	require.Equal(t, 1, len(rs.enclosures))
	require.Equal(t, *enc, rs.enclosures[0])
}

func TestNewSplitter_Errors(t *testing.T) {
	enc := &Enclosure{
		Start: '{',
		End:   '}',
	}
	_, err := NewSplitter('/', enc, enc)
	require.Error(t, err)
	require.Equal(t, "existing start encloser ('{' in Enclosure[2])", err.Error())

	enc2 := &Enclosure{
		Start: '<',
		End:   '}',
	}
	_, err = NewSplitter('/', enc, enc2)
	require.Error(t, err)
	require.Equal(t, "existing end encloser ('}' in Enclosure[2])", err.Error())
}

func TestMustCreateSplitter_Panics(t *testing.T) {
	enc := &Enclosure{
		Start: '{',
		End:   '}',
	}
	require.NotPanics(t, func() {
		MustCreateSplitter('/', enc)
	})
	require.Panics(t, func() {
		MustCreateSplitter('/', enc, enc)
	})
}

func TestSplitter_Split(t *testing.T) {
	encs := []*Enclosure{
		{
			Start: '{',
			End:   '}',
		},
		{
			Start:     '\'',
			End:       '\'',
			IsQuote:   true,
			Escapable: true,
			Escape:    '\\',
		},
		{
			Start:     '"',
			End:       '"',
			IsQuote:   true,
			Escapable: true,
			Escape:    '\\',
		},
	}
	s, _ := NewSplitter('/', encs...)

	testCases := []struct {
		str    string
		expect []string
	}{
		{
			`/foo/{/}`,
			[]string{``, `foo`, `{/}`},
		},
		{
			`/foo/{{/}}`,
			[]string{``, `foo`, `{{/}}`},
		},
		{
			`foo/bar/"baz/qux"/'qux/"/"/"/"/"/"'/`,
			[]string{`foo`, `bar`, `"baz/qux"`, `'qux/"/"/"/"/"/"'`, ``},
		},
		{
			`foo/"\"/"/bar`,
			[]string{`foo`, `"\"/"`, `bar`},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%s", i+1, tc.str), func(t *testing.T) {
			result, err := s.Split(tc.str)
			require.NoError(t, err)
			require.Equal(t, tc.expect, result)
		})
	}
}

func TestSplitter_Split_DoubleEscapes(t *testing.T) {
	s, err := NewSplitter(',', DoubleQuotesDoubleEscaped)
	require.NoError(t, err)

	parts, err := s.Split(`"aa"","",,,,,,""""""""""bbb"`)
	require.NoError(t, err)
	require.Equal(t, 1, len(parts))
}

func TestSplitter_PostElementFix(t *testing.T) {
	s, err := NewSplitter(',')
	require.NoError(t, err)
	s.SetPostElementFixer(func(s string, pos int, captured int, subParts ...SubPart) (string, bool, error) {
		return strings.Trim(s, " "), true, nil
	})
	parts, err := s.Split(`a, b,  c     `)
	require.NoError(t, err)
	require.Equal(t, 3, len(parts))
	require.Equal(t, `a`, parts[0])
	require.Equal(t, `b`, parts[1])
	require.Equal(t, `c`, parts[2])
}

func TestSplitter_PostElementFix_Errors(t *testing.T) {
	s, err := NewSplitter(',')
	require.NoError(t, err)
	s.SetPostElementFixer(func(s string, pos int, captured int, subParts ...SubPart) (string, bool, error) {
		if s == "" {
			return "", false, errors.New("whoops")
		}
		return s, true, nil
	})
	_, err = s.Split(`a,b,c`)
	require.NoError(t, err)

	_, err = s.Split(`,b,c`)
	require.Error(t, err)
	require.Equal(t, `whoops`, err.Error())

	_, err = s.Split(`a,b,`)
	require.Error(t, err)
	require.Equal(t, `whoops`, err.Error())
}

func TestSplitter_PostElementCheck_Errors(t *testing.T) {
	s, err := NewSplitter(',')
	require.NoError(t, err)
	s.SetPostElementFixer(func(s string, pos int, captured int, subParts ...SubPart) (string, bool, error) {
		if s == "" && captured == 0 {
			return "", false, errors.New("first cannot be empty")
		} else if s == "" {
			return "", false, nil
		}
		return s, true, nil
	})

	parts, err := s.Split(`aaa,bbb,ccc`)
	require.NoError(t, err)
	require.Equal(t, 3, len(parts))

	parts, err = s.Split(`,bbb,ccc`)
	require.Error(t, err)
	require.Equal(t, `first cannot be empty`, err.Error())

	parts, err = s.Split(`aaa,,ccc`)
	require.NoError(t, err)
	require.Equal(t, 2, len(parts))
}

func TestSplitter_Split_Errors(t *testing.T) {
	encs := []*Enclosure{
		{
			Start: '{',
			End:   '}',
		},
		{
			Start:     '\'',
			End:       '\'',
			IsQuote:   true,
			Escapable: true,
			Escape:    '\'',
		},
		{
			Start:     '"',
			End:       '"',
			IsQuote:   true,
			Escapable: true,
			Escape:    '\\',
		},
	}
	s, _ := NewSplitter('/', encs...)

	testCases := []struct {
		str       string
		expectErr string
	}{
		{
			`}`,
			fmt.Sprintf(unopenedFmt, "}", 0),
		},
		{
			`{},{}}`,
			fmt.Sprintf(unopenedFmt, "}", 5),
		},
		{
			`{/}}`,
			fmt.Sprintf(unopenedFmt, "}", 3),
		},
		{
			`{{{/}}`,
			fmt.Sprintf(unclosedFmt, "{", 0),
		},
		{
			`{{{/}`,
			fmt.Sprintf(unclosedFmt, "{", 1),
		},
		{
			`"`,
			fmt.Sprintf(unclosedFmt, `"`, 0),
		},
		{
			`"\"`,
			fmt.Sprintf(unclosedFmt, `"`, 0),
		},
		{
			`"\"""`,
			fmt.Sprintf(unclosedFmt, `"`, 4),
		},
		{
			`'`,
			fmt.Sprintf(unclosedFmt, `'`, 0),
		},
		{
			`'''`,
			fmt.Sprintf(unclosedFmt, `'`, 0),
		},
		{
			`'''''`,
			fmt.Sprintf(unclosedFmt, `'`, 0),
		},
		{
			`'''''''`,
			fmt.Sprintf(unclosedFmt, `'`, 0),
		},
		{
			`{'`,
			fmt.Sprintf(unclosedFmt, `'`, 1),
		},
		{
			`{\'`,
			fmt.Sprintf(unclosedFmt, `'`, 2),
		},
		{
			`{''`,
			fmt.Sprintf(unclosedFmt, `{`, 0),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%s", i+1, tc.str), func(t *testing.T) {
			_, err := s.Split(tc.str)
			require.Error(t, err)
			require.Equal(t, tc.expectErr, err.Error())
		})
	}
}

func TestSplitter_Split_Contiguous(t *testing.T) {
	s, err := NewSplitter(',', DoubleQuotesDoubleEscaped, SingleQuotes, CurlyBrackets)
	require.NoError(t, err)
	subs := 0
	startPositions := make([]int, 0)
	endPositions := make([]int, 0)
	isQuotes := make([]bool, 0)
	isBrackets := make([]bool, 0)
	isFixeds := make([]bool, 0)
	escapables := make([]bool, 0)
	startRunes := make([]rune, 0)
	endRunes := make([]rune, 0)
	escRunes := make([]rune, 0)
	unescaped := make([]string, 0)
	whitespacing := make([]bool, 0)
	types := make([]SubPartType, 0)
	s.SetPostElementFixer(func(s string, pos int, captured int, subParts ...SubPart) (string, bool, error) {
		subs = len(subParts)
		for _, sub := range subParts {
			startPositions = append(startPositions, sub.StartPos())
			endPositions = append(endPositions, sub.EndPos())
			isQuotes = append(isQuotes, sub.IsQuote())
			isBrackets = append(isBrackets, sub.IsBrackets())
			isFixeds = append(isFixeds, sub.IsFixed())
			escapables = append(escapables, sub.Escapable())
			startRunes = append(startRunes, sub.StartRune())
			endRunes = append(endRunes, sub.EndRune())
			escRunes = append(escRunes, sub.EscapeRune())
			unescaped = append(unescaped, sub.UnEscaped())
			whitespacing = append(whitespacing, sub.IsWhitespaceOnly(" "))
			types = append(types, sub.Type())
		}
		return s, true, nil
	})
	parts, err := s.Split(` "bbb""" '222'{'a','b','c'} `)
	//                        0123456789012345678901234567
	require.NoError(t, err)
	require.Equal(t, 1, len(parts))
	require.Equal(t, 6, subs)
	require.Equal(t, 6, len(startPositions))
	require.Equal(t, 6, len(endPositions))
	require.Equal(t, 6, len(isQuotes))
	require.Equal(t, 6, len(isBrackets))
	require.Equal(t, 6, len(isFixeds))
	require.Equal(t, 6, len(escapables))
	require.Equal(t, 6, len(startRunes))
	require.Equal(t, 6, len(endRunes))
	require.Equal(t, 6, len(escRunes))
	require.Equal(t, 6, len(unescaped))
	require.Equal(t, 6, len(whitespacing))
	require.Equal(t, 6, len(types))

	require.Equal(t, 0, startPositions[0])
	require.Equal(t, 0, endPositions[0])
	require.Equal(t, 1, startPositions[1])
	require.Equal(t, 7, endPositions[1])
	require.Equal(t, 8, startPositions[2])
	require.Equal(t, 8, endPositions[2])
	require.Equal(t, 9, startPositions[3])
	require.Equal(t, 13, endPositions[3])
	require.Equal(t, 14, startPositions[4])
	require.Equal(t, 26, endPositions[4])
	require.Equal(t, 27, startPositions[5])
	require.Equal(t, 27, endPositions[5])

	require.True(t, isFixeds[0])
	require.True(t, isQuotes[1])
	require.True(t, isFixeds[2])
	require.True(t, isQuotes[3])
	require.True(t, isBrackets[4])
	require.True(t, isFixeds[5])

	require.Equal(t, Fixed, types[0])
	require.Equal(t, Quotes, types[1])
	require.Equal(t, Fixed, types[2])
	require.Equal(t, Quotes, types[3])
	require.Equal(t, Brackets, types[4])
	require.Equal(t, Fixed, types[5])

	require.False(t, escapables[0])
	require.True(t, escapables[1])
	require.False(t, escapables[2])

	require.True(t, whitespacing[0])
	require.False(t, whitespacing[1])
	require.True(t, whitespacing[2])
	require.False(t, whitespacing[3])
	require.False(t, whitespacing[4])
	require.True(t, whitespacing[5])

	require.Equal(t, int32(0), startRunes[0])
	require.Equal(t, int32(0), endRunes[0])
	require.Equal(t, '"', startRunes[1])
	require.Equal(t, '"', endRunes[1])
	require.Equal(t, int32(0), startRunes[2])
	require.Equal(t, int32(0), endRunes[2])
	require.Equal(t, '\'', startRunes[3])
	require.Equal(t, '\'', endRunes[3])
	require.Equal(t, '{', startRunes[4])
	require.Equal(t, '}', endRunes[4])
	require.Equal(t, int32(0), startRunes[5])
	require.Equal(t, int32(0), endRunes[5])

	require.Equal(t, int32(0), escRunes[0])
	require.Equal(t, '"', escRunes[1])
	require.Equal(t, int32(0), escRunes[2])
	require.Equal(t, int32(0), escRunes[3])
	require.Equal(t, int32(0), escRunes[4])
	require.Equal(t, int32(0), escRunes[5])
	require.Equal(t, ` `, unescaped[0])
	require.Equal(t, `bbb"`, unescaped[1])
	require.Equal(t, ` `, unescaped[2])
	require.Equal(t, `222`, unescaped[3])
	require.Equal(t, `{'a','b','c'}`, unescaped[4])
	require.Equal(t, ` `, unescaped[5])
}
