package splitter

import (
	"fmt"
)

// Splitter is the actual splitter interface
type Splitter interface {
	// Split performs a split on the supplied string - returns the split parts and any error encountered
	//
	// If an error is returned, it will always be of type splittingError
	Split(s string, options ...Option) ([]string, error)
	// AddDefaultOptions adds default options for the splitter (other options can also be added when using Split)
	AddDefaultOptions(options ...Option) Splitter
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
		separator:   separator,
		enclosures:  make([]Enclosure, 0, len(encs)),
		openers:     map[rune]Enclosure{},
		closers:     map[rune]Enclosure{},
		defOptions:  make([]Option, 0),
		seenOptions: map[Option]bool{},
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
	separator   rune
	enclosures  []Enclosure
	openers     map[rune]Enclosure
	closers     map[rune]Enclosure
	defOptions  []Option
	seenOptions map[Option]bool
}

func (s *splitter) Split(str string, options ...Option) ([]string, error) {
	return newSplitterContext(str, s, s.mergeOptions(options)).split()
}

func (s *splitter) AddDefaultOptions(options ...Option) Splitter {
	for _, opt := range options {
		if opt != nil && !s.seenOptions[opt] {
			s.defOptions = append(s.defOptions, opt)
			s.seenOptions[opt] = true
		}
	}
	return s
}

func (s *splitter) mergeOptions(addOpts []Option) []Option {
	addLen := len(addOpts)
	defLen := len(s.defOptions)
	if addLen == 0 {
		return s.defOptions
	} else if defLen == 0 && addLen != 0 {
		result := make([]Option, 0, addLen)
		seen := map[Option]bool{}
		for _, opt := range addOpts {
			if opt != nil && !seen[opt] {
				result = append(result, opt)
			}
		}
		return result
	}
	result := make([]Option, 0, len(s.defOptions)+len(addOpts))
	result = append(result, s.defOptions...)
	seen := map[Option]bool{}
	for _, opt := range addOpts {
		if opt != nil && !seen[opt] && !s.seenOptions[opt] {
			result = append(result, opt)
			seen[opt] = true
		}
	}
	return result
}

type splitterContext struct {
	splitter *splitter
	options  []Option
	runes    []rune
	len      int
	lastAt   int
	current  *subPart
	stack    []*subPart
	delims   []SubPart
	captured []string
	skipped  int
}

func newSplitterContext(str string, splitter *splitter, options []Option) *splitterContext {
	runes := []rune(str)
	cp := 1
	for _, r := range runes {
		if r == splitter.separator {
			cp++
		}
	}
	return &splitterContext{
		splitter: splitter,
		options:  options,
		runes:    runes,
		lastAt:   0,
		len:      len(runes),
		current:  nil,
		stack:    make([]*subPart, 0),
		delims:   make([]SubPart, 0),
		captured: make([]string, 0, cp),
	}
}

func (ctx *splitterContext) split() ([]string, error) {
	for i := 0; i < ctx.len; i++ {
		r := ctx.runes[i]
		if r == ctx.splitter.separator {
			if !ctx.inAny() {
				if err := ctx.purge(i, false); err != nil {
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
	if err := ctx.purge(ctx.len, true); err != nil {
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

func (ctx *splitterContext) purge(i int, isLast bool) (err error) {
	if i >= ctx.lastAt {
		ctx.purgeFixed(i)
		capture := string(ctx.runes[ctx.lastAt:i])
		addIt := true
		cLen := len(ctx.captured)
		for _, o := range ctx.options {
			capture, addIt, err = o.Apply(capture, ctx.lastAt, ctx.len, cLen, ctx.skipped, isLast, ctx.delims...)
			if !addIt || err != nil {
				break
			}
		}
		err = asSplittingError(err, ctx.lastAt)
		if addIt {
			ctx.captured = append(ctx.captured, capture)
		} else {
			ctx.skipped++
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
	ctx.current = &subPart{
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
		ctx.delims = append(ctx.delims, &subPart{
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
