package splitter

import (
	"fmt"
)

// Splitter is the actual splitter interface
type Splitter interface {
	// Split performs a split on the supplied string - returns the split parts and any error encountered
	Split(str string) ([]string, error)
	// PostElementFix supplies a function that can fix-up elements prior to being added (e.g. trimming)
	PostElementFix(f PostElementFix) Splitter
}

type PostElementFix func(s string) (string, error)

// NewSplitter creates a new splitter
//
// the `separator` arg is the rune on which to split
//
// the optional `encs` varargs are the enclosures (e.g. brackets, quotes) to be taken into consideration when splitting
//
// An error is returned if any of enclosures specified match any other enclosure `Start`/`End`
func NewSplitter(separator rune, encs ...*Enclosure) (Splitter, error) {
	result := &splitter{
		separator:  separator,
		enclosures: make([]Enclosure, 0, len(encs)),
		openers:    map[rune]Enclosure{},
		closers:    map[rune]bool{},
	}

	for i, enc := range encs {
		if enc != nil {
			if _, ok := result.openers[enc.Start]; ok {
				return nil, fmt.Errorf("existing start encloser ('%s' in Enclosure[%d])", string(enc.Start), i+1)
			}
			if result.closers[enc.End] {
				return nil, fmt.Errorf("existing end encloser ('%s' in Enclosure[%d])", string(enc.End), i+1)
			}
			result.openers[enc.Start] = enc.clone()
			result.closers[enc.End] = true
			result.enclosures = append(result.enclosures, *enc)
		}
	}
	return result, nil
}

// MustCreateSplitter is the same as NewSplitter, except that it panics in case of error
func MustCreateSplitter(separator rune, encs ...*Enclosure) Splitter {
	if s, err := NewSplitter(separator, encs...); err != nil {
		panic(any(err))
	} else {
		return s
	}
}

type splitter struct {
	separator  rune
	enclosures []Enclosure
	openers    map[rune]Enclosure
	closers    map[rune]bool
	postFix    PostElementFix
}

func (s *splitter) Split(str string) ([]string, error) {
	return newSplitterContext(str, s).split()
}

func (s *splitter) PostElementFix(f PostElementFix) Splitter {
	s.postFix = f
	return s
}

type splitterContext struct {
	splitter *splitter
	runes    []rune
	len      int
	lastAt   int
	captured []string
	current  *delimiter
	stack    []*delimiter
}

func newSplitterContext(str string, splitter *splitter) *splitterContext {
	runes := []rune(str)
	return &splitterContext{
		splitter: splitter,
		runes:    runes,
		lastAt:   0,
		len:      len(runes),
		current:  nil,
		stack:    []*delimiter{},
		captured: make([]string, 0),
	}
}

func (c *splitterContext) split() ([]string, error) {
	for i, r := range c.runes {
		if r == c.splitter.separator {
			if !c.inAny() {
				if err := c.purge(i); err != nil {
					return nil, err
				}
			}
		} else if isEnd, inQuote := c.isQuoteEnd(r, i); isEnd {
			c.pop()
		} else if !inQuote {
			if c.isClose(r) {
				c.pop()
			} else if enc, isOpener := c.isOpener(r); isOpener {
				c.push(enc, i)
			} else if c.splitter.closers[r] {
				return nil, fmt.Errorf("unopened '%s' at position %d", string(r), i)
			}
		}
	}
	if c.inAny() {
		return nil, fmt.Errorf("unclosed '%s' at position %d", string(c.current.enc.Start), c.current.openPos)
	}
	if err := c.purge(c.len); err != nil {
		return nil, err
	}
	return c.captured, nil
}

func (c *splitterContext) purge(i int) (err error) {
	if i >= c.lastAt {
		capture := string(c.runes[c.lastAt:i])
		if c.splitter.postFix != nil {
			capture, err = c.splitter.postFix(capture)
		}
		c.captured = append(c.captured, capture)
		c.lastAt = i + 1
	}
	return
}

func (c *splitterContext) inAny() bool {
	return c.current != nil
}

func (c *splitterContext) isQuoteEnd(r rune, pos int) (isEnd bool, inQuote bool) {
	if c.current != nil && c.current.enc.IsQuote {
		inQuote = true
		if c.current.enc.End == r {
			isEnd = true
			if c.current.enc.Escapable {
				escaped := false
				minPos := c.current.openPos
				for i := pos - 1; i > minPos; i-- {
					if c.runes[i] == c.current.enc.Escape {
						escaped = !escaped
					} else {
						break
					}
				}
				isEnd = !escaped
			}
		}
	}
	return
}

func (c *splitterContext) isClose(r rune) bool {
	return c.current != nil && c.current.enc.End == r
}

func (c *splitterContext) isOpener(r rune) (Enclosure, bool) {
	enc, ok := c.splitter.openers[r]
	return enc, ok
}

func (c *splitterContext) push(enc Enclosure, pos int) {
	if c.current != nil {
		c.stack = append(c.stack, c.current)
	}
	c.current = &delimiter{
		openPos: pos,
		enc:     enc,
	}
}

func (c *splitterContext) pop() {
	if l := len(c.stack); l > 0 {
		c.current = c.stack[l-1]
		c.stack = c.stack[0 : l-1]
	} else {
		c.current = nil
	}
}

type delimiter struct {
	enc     Enclosure
	openPos int
}
