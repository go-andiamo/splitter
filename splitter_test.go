package splitter

import (
	"fmt"
	"github.com/stretchr/testify/require"
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
		{
			``,
			[]string{``},
		},
		{
			` `,
			[]string{` `},
		},
		{
			`/`,
			[]string{``, ``},
		},
		{
			`//`,
			[]string{``, ``, ``},
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

	c := &infoCapture{}
	parts, err := s.Split(` "bbb""" '222'{'a','b','c'} `, c)
	//                        0123456789012345678901234567
	require.NoError(t, err)
	require.Equal(t, 1, len(parts))
	require.Equal(t, 1, c.called)
	require.Equal(t, 6, c.subs)
	require.Equal(t, 6, len(c.startPositions))
	require.Equal(t, 6, len(c.endPositions))
	require.Equal(t, 6, len(c.isQuotes))
	require.Equal(t, 6, len(c.isBrackets))
	require.Equal(t, 6, len(c.isFixeds))
	require.Equal(t, 6, len(c.escapables))
	require.Equal(t, 6, len(c.startRunes))
	require.Equal(t, 6, len(c.endRunes))
	require.Equal(t, 6, len(c.escRunes))
	require.Equal(t, 6, len(c.unescaped))
	require.Equal(t, 6, len(c.whitespacing))
	require.Equal(t, 6, len(c.types))

	require.Equal(t, 0, c.startPositions[0])
	require.Equal(t, 0, c.endPositions[0])
	require.Equal(t, 1, c.startPositions[1])
	require.Equal(t, 7, c.endPositions[1])
	require.Equal(t, 8, c.startPositions[2])
	require.Equal(t, 8, c.endPositions[2])
	require.Equal(t, 9, c.startPositions[3])
	require.Equal(t, 13, c.endPositions[3])
	require.Equal(t, 14, c.startPositions[4])
	require.Equal(t, 26, c.endPositions[4])
	require.Equal(t, 27, c.startPositions[5])
	require.Equal(t, 27, c.endPositions[5])

	require.True(t, c.isFixeds[0])
	require.True(t, c.isQuotes[1])
	require.True(t, c.isFixeds[2])
	require.True(t, c.isQuotes[3])
	require.True(t, c.isBrackets[4])
	require.True(t, c.isFixeds[5])

	require.Equal(t, Fixed, c.types[0])
	require.Equal(t, Quotes, c.types[1])
	require.Equal(t, Fixed, c.types[2])
	require.Equal(t, Quotes, c.types[3])
	require.Equal(t, Brackets, c.types[4])
	require.Equal(t, Fixed, c.types[5])

	require.False(t, c.escapables[0])
	require.True(t, c.escapables[1])
	require.False(t, c.escapables[2])

	require.True(t, c.whitespacing[0])
	require.False(t, c.whitespacing[1])
	require.True(t, c.whitespacing[2])
	require.False(t, c.whitespacing[3])
	require.False(t, c.whitespacing[4])
	require.True(t, c.whitespacing[5])

	require.Equal(t, int32(0), c.startRunes[0])
	require.Equal(t, int32(0), c.endRunes[0])
	require.Equal(t, '"', c.startRunes[1])
	require.Equal(t, '"', c.endRunes[1])
	require.Equal(t, int32(0), c.startRunes[2])
	require.Equal(t, int32(0), c.endRunes[2])
	require.Equal(t, '\'', c.startRunes[3])
	require.Equal(t, '\'', c.endRunes[3])
	require.Equal(t, '{', c.startRunes[4])
	require.Equal(t, '}', c.endRunes[4])
	require.Equal(t, int32(0), c.startRunes[5])
	require.Equal(t, int32(0), c.endRunes[5])

	require.Equal(t, int32(0), c.escRunes[0])
	require.Equal(t, '"', c.escRunes[1])
	require.Equal(t, int32(0), c.escRunes[2])
	require.Equal(t, int32(0), c.escRunes[3])
	require.Equal(t, int32(0), c.escRunes[4])
	require.Equal(t, int32(0), c.escRunes[5])
	require.Equal(t, ` `, c.unescaped[0])
	require.Equal(t, `bbb"`, c.unescaped[1])
	require.Equal(t, ` `, c.unescaped[2])
	require.Equal(t, `222`, c.unescaped[3])
	require.Equal(t, `{'a','b','c'}`, c.unescaped[4])
	require.Equal(t, ` `, c.unescaped[5])
}

type infoCapture struct {
	called         int
	subs           int
	startPositions []int
	endPositions   []int
	isQuotes       []bool
	isBrackets     []bool
	isFixeds       []bool
	escapables     []bool
	startRunes     []rune
	endRunes       []rune
	escRunes       []rune
	unescaped      []string
	whitespacing   []bool
	types          []SubPartType
}

func (o *infoCapture) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	o.called++
	o.subs = len(subParts)
	for _, sub := range subParts {
		o.startPositions = append(o.startPositions, sub.StartPos())
		o.endPositions = append(o.endPositions, sub.EndPos())
		o.isQuotes = append(o.isQuotes, sub.IsQuote())
		o.isBrackets = append(o.isBrackets, sub.IsBrackets())
		o.isFixeds = append(o.isFixeds, sub.IsFixed())
		o.escapables = append(o.escapables, sub.Escapable())
		o.startRunes = append(o.startRunes, sub.StartRune())
		o.endRunes = append(o.endRunes, sub.EndRune())
		o.escRunes = append(o.escRunes, sub.EscapeRune())
		o.unescaped = append(o.unescaped, sub.UnEscaped())
		o.whitespacing = append(o.whitespacing, sub.IsWhitespaceOnly(" "))
		o.types = append(o.types, sub.Type())
	}
	return s, true, nil
}

func TestSplitter_SetDefaultOptions(t *testing.T) {
	s, err := NewSplitter(',')
	require.NoError(t, err)
	rs, ok := s.(*splitter)
	require.True(t, ok)
	require.Equal(t, 0, len(rs.defOptions))

	s.AddDefaultOptions(nil, NotEmptyLast, NotEmptyLast)
	require.Equal(t, 1, len(rs.defOptions))

	s.AddDefaultOptions(nil, NotEmptyFirst, NotEmptyFirst)
	require.Equal(t, 2, len(rs.defOptions))
}

func TestEnsureSplitterContextOptionsSegregated(t *testing.T) {
	s, err := NewSplitter(',')
	require.NoError(t, err)

	opt := &addOptionsOption{
		splitter: s,
	}

	pts, err := s.Split(`a,b,c,`, opt)
	require.NoError(t, err)
	require.Equal(t, 4, len(pts))
}

type addOptionsOption struct {
	splitter Splitter
}

func (o *addOptionsOption) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	// add an option to the original splitter - which should not be seen in the current splitter context...
	o.splitter.AddDefaultOptions(IgnoreEmptyFirst, IgnoreEmptyLast)
	return s, true, nil
}

func TestBracketsEscaping(t *testing.T) {
	unescaped, err := NewSplitter('/', DoubleQuotes,
		Parenthesis, SquareBrackets)
	require.NoError(t, err)
	escaped, err := NewSplitter('/', DoubleQuotes,
		MustMakeEscapable(Parenthesis, '\\'), MustMakeEscapable(SquareBrackets, '\\'))
	testCases := []struct {
		str             string
		expectUnEscPass bool
		expectEscPass   bool
	}{
		{
			`(\(/)`,
			false,
			true,
		},
		{
			`(\)/)`,
			false,
			true,
		},
		{
			`(\[)`,
			false,
			true,
		},
		{
			`(\])`,
			false,
			true,
		},
		{
			`(/\[)`,
			false,
			true,
		},
		{
			`(/\])`,
			false,
			true,
		},
		{
			`\(\(\)`,
			false,
			true,
		},
		{
			`\[\[\]`,
			false,
			true,
		},
		{
			`(/\(\(\))`,
			false,
			true,
		},
		{
			`[/\[\[\]]`,
			false,
			true,
		},
		{
			`[/\(\(\)]`,
			false,
			true,
		},
		{
			`(/\[\[\])`,
			false,
			true,
		},
		{
			`(\(/)(`,
			false,
			false,
		},
		{
			`(\)/))`,
			false,
			false,
		},
		{
			`((\[)`,
			false,
			false,
		},
		{
			`((\])`,
			false,
			false,
		},
		{
			`[(/\[)`,
			false,
			false,
		},
		{
			`[(/\])`,
			false,
			false,
		},
		{
			`\(\(\)]`,
			false,
			false,
		},
		{
			`\[\[\]]`,
			true,
			false,
		},
		{
			`(/\(\(\)))`,
			true,
			false,
		},
		{
			`[/\[\[\]]]`,
			true,
			false,
		},
		{
			`[[/\(\(\)]`,
			false,
			false,
		},
		{
			`((/\[\[\])`,
			false,
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%s", i+1, tc.str), func(t *testing.T) {
			_, err = unescaped.Split(tc.str)
			if tc.expectUnEscPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			pts, err := escaped.Split(tc.str)
			if tc.expectEscPass {
				require.NoError(t, err)
				require.Equal(t, 1, len(pts))
			} else {
				require.Error(t, err)
			}
		})
	}
}
