package splitter

import (
	"fmt"
	"strings"
)

// Splitter is the actual splitter interface
type Splitter interface {
	// Split performs a split on the supplied string - returns the split parts and any error encountered
	Split(s string) ([]string, error)
	// SetPostElementFixer supplies a function that can fix-up elements prior to being added (e.g. trimming)
	// see PostElementFix for more information
	SetPostElementFixer(f PostElementFix) Splitter
}

// PostElementFix is the func that can be set using Splitter.SetPostElementFixer() to enable
// the captured split part to be altered prior to adding to the result
//
// It can also return false, to indicate that the part is not to be added to the result (e.g. stripping empty items)
// or an error to indicate the part is unacceptable (and cease splitting with that error)
//
// s - is the entire found string part found
//
// pos - is the start position (relative to the original string)
//
// captured - is the number of parts already captured
//
// subParts - is the sub-parts of the found string part (such as quote or bracket enclosures)
//
// The subParts can be used to determine if the overall part was made of, for example, contiguous quotes
// (which may need to be rejected - e.g. contiguous quotes in a CSV)
type PostElementFix func(s string, pos int, captured int, subParts ...SubPart) (string, bool, error)

type SubPartType int

const (
	Fixed SubPartType = iota
	Quotes
	Brackets
)

// SubPart is the interfaces passed PostElementCheck - to enable contiguous enclosures to be spotted
type SubPart interface {
	// StartPos returns the start position (relative to the original string) of the part
	StartPos() int
	// EndPos returns the end position (relative to the original string) of the part
	EndPos() int
	// IsQuote returns whether the part is quotes enclosure
	IsQuote() bool
	// IsBrackets returns whether the part is a brackets enclosure
	IsBrackets() bool
	// IsFixed returns whether the part was fixed text (i.e not quotes or brackets)
	IsFixed() bool
	// Type returns the part type (fixed, quotes, brackets)
	Type() SubPartType
	// Escapable returns whether the original enclosure was escapable
	Escapable() bool
	// StartRune returns the start rune of the enclosure
	StartRune() rune
	// EndRune returns the end rune of the enclosure
	EndRune() rune
	// EscapeRune returns the original escape rune for the enclosure
	EscapeRune() rune
	// UnEscaped returns the unescaped string
	//
	// If the part was a quote enclosure, the enclosing quote marks are stripped and, if escapable, any escaped quotes are transposed.
	// If the quote enclosure was not escapable, just the enclosing quote marks are removed
	//
	// If the part was not a quote enclosure, the original string part is returned
	UnEscaped() string
	IsWhitespaceOnly(cutset ...string) bool
}

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
		closers:    map[rune]Enclosure{},
	}
	for i, enc := range encs {
		if enc != nil {
			if _, exists := result.openers[enc.Start]; exists {
				return nil, fmt.Errorf("existing start encloser ('%s' in Enclosure[%d])", string(enc.Start), i+1)
			}
			if _, exists := result.closers[enc.End]; exists {
				return nil, fmt.Errorf("existing end encloser ('%s' in Enclosure[%d])", string(enc.End), i+1)
			}
			cEnc := enc.clone()
			result.openers[enc.Start] = cEnc
			result.closers[enc.End] = cEnc
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
	closers    map[rune]Enclosure
	postFix    PostElementFix
}

func (s *splitter) Split(str string) ([]string, error) {
	return newSplitterContext(str, s).split()
}

func (s *splitter) SetPostElementFixer(f PostElementFix) Splitter {
	s.postFix = f
	return s
}

type splitterContext struct {
	splitter *splitter
	runes    []rune
	len      int
	lastAt   int
	current  *delimitedEnclosure
	stack    []*delimitedEnclosure
	delims   []SubPart
	captured []string
}

func newSplitterContext(str string, splitter *splitter) *splitterContext {
	runes := []rune(str)
	cp := 1
	for _, r := range runes {
		if r == splitter.separator {
			cp++
		}
	}
	return &splitterContext{
		splitter: splitter,
		runes:    runes,
		lastAt:   0,
		len:      len(runes),
		current:  nil,
		stack:    make([]*delimitedEnclosure, 0),
		delims:   make([]SubPart, 0),
		captured: make([]string, 0, cp),
	}
}

func (ctx *splitterContext) split() ([]string, error) {
	for i := 0; i < ctx.len; i++ {
		r := ctx.runes[i]
		if r == ctx.splitter.separator {
			if !ctx.inAny() {
				if err := ctx.purge(i); err != nil {
					return nil, err
				}
			}
		} else if isEnd, inQuote, inc := ctx.isQuoteEnd(r, i); isEnd {
			ctx.pop(i)
		} else {
			i += inc
			if !inQuote {
				if ctx.isClose(r) {
					ctx.pop(i)
				} else if enc, isOpener := ctx.isOpener(r); isOpener {
					ctx.push(enc, i)
				} else if cEnc, ok := ctx.splitter.closers[r]; ok {
					return nil, newSplittingError(Unopened, i, r, &cEnc)
				}
			}
		}
	}
	if ctx.inAny() {
		return nil, newSplittingError(Unclosed, ctx.current.openPos, ctx.current.enc.Start, &ctx.current.enc)
	}
	if err := ctx.purge(ctx.len); err != nil {
		return nil, err
	}
	return ctx.captured, nil
}

func (ctx *splitterContext) isQuoteEnd(r rune, pos int) (isEnd bool, inQuote bool, skip int) {
	if ctx.current != nil && ctx.current.enc.IsQuote {
		inQuote = true
		if ctx.current.enc.End == r {
			isEnd = true
			if ctx.current.enc.isDoubleEscaping() {
				if pos < ctx.len-1 && ctx.runes[pos+1] == r {
					isEnd = false
					skip = 1
				}
			} else if ctx.current.enc.isEscapable() {
				escaped := false
				minPos := ctx.current.openPos
				for i := pos - 1; i > minPos; i-- {
					if ctx.runes[i] == ctx.current.enc.Escape {
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

func (ctx *splitterContext) purge(i int) (err error) {
	if i >= ctx.lastAt {
		ctx.purgeFixed(i)
		capture := string(ctx.runes[ctx.lastAt:i])
		addIt := true
		if ctx.splitter.postFix != nil {
			capture, addIt, err = ctx.splitter.postFix(capture, ctx.lastAt, len(ctx.captured), ctx.delims...)
		}
		if addIt {
			ctx.captured = append(ctx.captured, capture)
		}
		ctx.lastAt = i + 1
		ctx.delims = make([]SubPart, 0)
	}
	return
}

func (ctx *splitterContext) inAny() bool {
	return ctx.current != nil
}

func (ctx *splitterContext) isClose(r rune) bool {
	return ctx.current != nil && ctx.current.enc.End == r
}

func (ctx *splitterContext) isOpener(r rune) (Enclosure, bool) {
	enc, ok := ctx.splitter.openers[r]
	return enc, ok
}

func (ctx *splitterContext) push(enc Enclosure, pos int) {
	if ctx.current != nil {
		ctx.stack = append(ctx.stack, ctx.current)
	}
	ctx.current = &delimitedEnclosure{
		openPos: pos,
		enc:     enc,
		ctx:     ctx,
	}
	if len(ctx.stack) == 0 {
		ctx.purgeFixed(pos)
		ctx.delims = append(ctx.delims, ctx.current)
	}
}

func (ctx *splitterContext) purgeFixed(pos int) {
	last := ctx.lastAt
	if len(ctx.delims) > 0 {
		last = ctx.delims[len(ctx.delims)-1].EndPos() + 1
	}
	if last < pos {
		ctx.delims = append(ctx.delims, &delimitedEnclosure{
			enc:      Enclosure{},
			openPos:  last,
			closePos: pos - 1,
			ctx:      ctx,
			fixed:    true,
		})
	}
}

func (ctx *splitterContext) pop(pos int) {
	ctx.current.closePos = pos
	if l := len(ctx.stack); l > 0 {
		ctx.current = ctx.stack[l-1]
		ctx.stack = ctx.stack[0 : l-1]
	} else {
		ctx.current = nil
	}
}

type delimitedEnclosure struct {
	enc      Enclosure
	openPos  int
	closePos int
	ctx      *splitterContext
	fixed    bool
}

func (d *delimitedEnclosure) StartPos() int {
	return d.openPos
}
func (d *delimitedEnclosure) EndPos() int {
	return d.closePos
}
func (d *delimitedEnclosure) IsQuote() bool {
	return d.enc.IsQuote
}
func (d *delimitedEnclosure) IsBrackets() bool {
	return !d.fixed && !d.enc.IsQuote
}
func (d *delimitedEnclosure) IsFixed() bool {
	return d.fixed
}
func (d *delimitedEnclosure) Type() SubPartType {
	if d.fixed {
		return Fixed
	} else if !d.enc.IsQuote {
		return Brackets
	}
	return Quotes
}
func (d *delimitedEnclosure) Escapable() bool {
	return d.enc.isEscapable()
}
func (d *delimitedEnclosure) StartRune() rune {
	return d.enc.Start
}
func (d *delimitedEnclosure) EndRune() rune {
	return d.enc.End
}
func (d *delimitedEnclosure) EscapeRune() rune {
	return d.enc.Escape
}
func (d *delimitedEnclosure) UnEscaped() string {
	if d.fixed || !d.enc.IsQuote {
		return string(d.ctx.runes[d.openPos : d.closePos+1])
	} else if !d.enc.isEscapable() {
		return string(d.ctx.runes[d.openPos+1 : d.closePos])
	}
	return strings.ReplaceAll(string(d.ctx.runes[d.openPos+1:d.closePos]), string([]rune{d.enc.Escape, d.enc.End}), string(d.enc.End))
}
func (d *delimitedEnclosure) IsWhitespaceOnly(cutset ...string) bool {
	if !d.fixed {
		return false
	}
	cuts := " \t\n"
	if len(cutset) > 0 {
		cuts = strings.Join(cutset, "")
	}
	return strings.Trim(string(d.ctx.runes[d.openPos:d.closePos+1]), cuts) == ""
}
