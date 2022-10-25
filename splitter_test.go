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

func TestSplitter_SplitCsv(t *testing.T) {
	enc := &Enclosure{
		Start:     '"',
		End:       '"',
		IsQuote:   true,
		Escapable: true,
		Escape:    '"',
	}
	s, err := NewSplitter(',', enc)
	require.NoError(t, err)
	parts, err := s.Split(`"aaa","bbb","cc""cc"`)
	require.NoError(t, err)
	require.Equal(t, 3, len(parts))
	require.Equal(t, `"aaa"`, parts[0])
	require.Equal(t, `"bbb"`, parts[1])
	require.Equal(t, `"cc""cc"`, parts[2])

	parts, err = s.Split(`"aaa","cc""""cc"`)
	require.NoError(t, err)
	require.Equal(t, 2, len(parts))
	require.Equal(t, `"aaa"`, parts[0])
	require.Equal(t, `"cc""""cc"`, parts[1])

	parts, err = s.Split(`"aaa",""ccc""`)
	require.NoError(t, err)
	require.Equal(t, 2, len(parts))
	require.Equal(t, `"aaa"`, parts[0])
	require.Equal(t, `""ccc""`, parts[1])

	_, err = s.Split(`"aaa","cc"""cc"`)
	//                    012345678901234
	//                    o   c o  __c  o
	require.Error(t, err)
	require.Equal(t, `unclosed '"' at position 14`, err.Error())
}

func TestSplitter_PostElementFix(t *testing.T) {
	s, err := NewSplitter(',')
	require.NoError(t, err)
	s.PostElementFix(func(s string) (string, error) {
		return strings.Trim(s, " "), nil
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
	s.PostElementFix(func(s string) (string, error) {
		if s == "" {
			return "", errors.New("whoops")
		}
		return s, nil
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
		str       string
		expectErr string
	}{
		{
			`{/}}`,
			`unopened '}' at position 3`,
		},
		{
			`{{{/}}`,
			`unclosed '{' at position 0`,
		},
		{
			`{{{/}`,
			`unclosed '{' at position 1`,
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
