package splitter

import "strings"

// SubPartType denotes the type of the SubPart (as returned from SubPart.Type)
type SubPartType int

const (
	Fixed SubPartType = iota
	Quotes
	Brackets
)

// SubPart is the interfaces passed to Options.Apply - to enable examination of sub-parts found in a split part
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
	// String returns the actual raw string of the part
	String() string
	// IsWhitespaceOnly returns whether the item is whitespace only (using the given trim cutset)
	IsWhitespaceOnly(cutset ...string) bool

	Enclosure() *Enclosure
}

type subPart struct {
	enc      Enclosure
	openPos  int
	closePos int
	ctx      *splitterContext
	fixed    bool
}

func (s *subPart) StartPos() int {
	return s.openPos
}

func (s *subPart) EndPos() int {
	return s.closePos
}

func (s *subPart) IsQuote() bool {
	return s.enc.IsQuote
}

func (s *subPart) IsBrackets() bool {
	return !s.fixed && !s.enc.IsQuote
}

func (s *subPart) IsFixed() bool {
	return s.fixed
}

func (s *subPart) Type() SubPartType {
	if s.fixed {
		return Fixed
	} else if !s.enc.IsQuote {
		return Brackets
	}
	return Quotes
}

func (s *subPart) Escapable() bool {
	return s.enc.isEscapable()
}

func (s *subPart) StartRune() rune {
	return s.enc.Start
}

func (s *subPart) EndRune() rune {
	return s.enc.End
}

func (s *subPart) EscapeRune() rune {
	return s.enc.Escape
}

func (s *subPart) UnEscaped() string {
	if s.fixed || !s.enc.IsQuote {
		return string(s.ctx.runes[s.openPos : s.closePos+1])
	} else if !s.enc.isEscapable() {
		return string(s.ctx.runes[s.openPos+1 : s.closePos])
	}
	return strings.ReplaceAll(string(s.ctx.runes[s.openPos+1:s.closePos]), string([]rune{s.enc.Escape, s.enc.End}), string(s.enc.End))
}

func (s *subPart) String() string {
	return string(s.ctx.runes[s.openPos : s.closePos+1])
}

func (s *subPart) IsWhitespaceOnly(cutset ...string) bool {
	if !s.fixed {
		return false
	}
	cuts := " \t\n"
	if len(cutset) > 0 {
		cuts = strings.Join(cutset, "")
	}
	return strings.Trim(string(s.ctx.runes[s.openPos:s.closePos+1]), cuts) == ""
}

func (s *subPart) Enclosure() *Enclosure {
	return &s.enc
}
