package splitter

import (
	"strings"
)

type Option interface {
	Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error)
}

var (
	TrimSpaces            Option = _TrimSpaces            // TrimSpaces causes split parts to be trimmed of leading & trailing spaces
	Trim                         = _Trim                  // Trim causes split parts to be trimmed of the leading & trailing custsets specified
	NoEmpties             Option = _NoEmpties             // NoEmpties causes an error if a split part is empty
	NoEmptiesMsg                 = _NoEmptiesMsg          // NoEmptiesMsg is the same as NoEmpties but allows a custom error message
	IgnoreEmpties         Option = _IgnoreEmpties         // IgnoreEmpties causes empty split parts to not be added to the result
	NotEmptyFirst         Option = _NotEmptyFirst         // NotEmptyFirst causes an error if the first split part is empty
	NotEmptyFirstMsg             = _NotEmptyFirstMsg      // NotEmptyFirstMsg is the same as NotEmptyFirst but allows a custom error message
	IgnoreEmptyFirst      Option = _IgnoreEmptyFirst      // IgnoreEmptyFirst causes an empty first split part not to be added to the result
	NotEmptyLast          Option = _NotEmptyLast          // NotEmptyLast causes an error if the last split part is empty
	NotEmptyLastMsg              = _NotEmptyLastMsg       // NotEmptyLastMsg is the same as NotEmptyLast but allows a custom error message
	IgnoreEmptyLast       Option = _IgnoreEmptyLast       // IgnoreEmptyLast causes an empty last split part not to be added to the result
	NotEmptyInners        Option = _NotEmptyInners        // NotEmptyInners causes an error if an inner (i.e. not first or last) split part is empty
	NotEmptyInnersMsg            = _NotEmptyInnersMsg     // NotEmptyInnersMsg is the same as NotEmptyInners but allows a custom error message
	IgnoreEmptyInners     Option = _IgnoreEmptyInners     // IgnoreEmptyInners causes empty inner (i.e. not first or last) split parts not to be added to the result
	NotEmptyOuters        Option = _NotEmptyOuters        // NotEmptyOuters causes an error if an outer (i.e. first or last) split part is empty (same as adding both NotEmptyFirst & NotEmptyLast)
	NotEmptyOutersMsg            = _NotEmptyOutersMsg     // NotEmptyOutersMsg is the same as NotEmptyOuters but allows a custom error message
	IgnoreEmptyOuters     Option = _IgnoreEmptyOuters     // IgnoreEmptyOuters causes empty outer (i.e. first or last) split parts not to be added to the result (same as adding both IgnoreEmptyFirst & IgnoreEmptyLast)
	NoContiguousQuotes    Option = _NoContiguousQuotes    // NoContiguousQuotes causes an error if there are contiguous quotes within a split part
	NoContiguousQuotesMsg        = _NoContiguousQuotesMsg // NoContiguousQuotesMsg is the same as NoContiguousQuotes but allows a custom error message
	NoMultiQuotes         Option = _NoMultiQuotes         // NoMultiQuotes causes an error if there are multiple (not necessarily contiguous) quotes within a split part
	NoMultiQuotesMsg             = _NoMultiQuotesMsg      // NoMultiQuotesMsg is the same as NoMultiQuotes but allows a custom error message
	NoMultis              Option = _NoMultis              // NoMultis causes an error if there are multiple quotes or brackets in a split part
	NoMultisMsg                  = _NoMultisMsg           // NoMultisMsg is the same as NoMultis but allows a custom error message
	StripQuotes           Option = _StripQuotes           // StripQuotes causes quotes within a split part to be stripped
	UnescapeQuotes        Option = _UnescapeQuotes        // UnescapeQuotes causes any quotes within the split part to have any escaped end quotes to be removed
)

var (
	_TrimSpaces = &trim{cutset: " "}
	_Trim       = func(cutset string) Option {
		return &trim{cutset: cutset}
	}
	_NoEmpties    = &noEmpties{message: "split items cannot be empty"}
	_NoEmptiesMsg = func(message string) Option {
		return &noEmpties{message: message}
	}
	_IgnoreEmpties    = &ignoreEmpties{}
	_NotEmptyFirst    = &notEmptyFirst{message: "first split item cannot be empty"}
	_NotEmptyFirstMsg = func(message string) Option {
		return &notEmptyFirst{message: message}
	}
	_IgnoreEmptyFirst = &ignoreEmptyFirst{}
	_NotEmptyLast     = &notEmptyLast{message: "last split item cannot be empty"}
	_NotEmptyLastMsg  = func(message string) Option {
		return &notEmptyLast{message: message}
	}
	_IgnoreEmptyLast   = &ignoreEmptyLast{}
	_NotEmptyInners    = &notEmptyInners{message: "inner items cannot be empty"}
	_NotEmptyInnersMsg = func(message string) Option {
		return &notEmptyInners{message: message}
	}
	_IgnoreEmptyInners = &ignoreEmptyInners{}
	_NotEmptyOuters    = &notEmptyOuters{message: "first/last items cannot be empty"}
	_NotEmptyOutersMsg = func(message string) Option {
		return &notEmptyOuters{message: message}
	}
	_IgnoreEmptyOuters     = &ignoreEmptyOuters{}
	_NoContiguousQuotes    = &noContiguousQuotes{message: "split item cannot have contiguous quotes"}
	_NoContiguousQuotesMsg = func(message string) Option {
		return &noContiguousQuotes{message: message}
	}
	_NoMultiQuotes    = &noMultiQuotes{message: "split item cannot have multiple quotes"}
	_NoMultiQuotesMsg = func(message string) Option {
		return &noMultiQuotes{message: message}
	}
	_NoMultis    = &noMultis{message: "split item cannot have multiple parts"}
	_NoMultisMsg = func(message string) Option {
		return &noMultis{message: message}
	}
	_StripQuotes    = &stripQuotes{}
	_UnescapeQuotes = &unescapeQuotes{}
)

type trim struct {
	cutset string
}

func (o *trim) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	return strings.Trim(s, o.cutset), true, nil
}

type noEmpties struct {
	message string
}

func (o *noEmpties) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" {
		return "", false, NewOptionFailError(o.message, pos, nil)
	}
	return s, true, nil
}

type ignoreEmpties struct {
}

func (o *ignoreEmpties) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" {
		return "", false, nil
	}
	return s, true, nil
}

type notEmptyFirst struct {
	message string
}

func (o *notEmptyFirst) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" && captured == 0 && skipped == 0 {
		return "", false, NewOptionFailError(o.message, pos, nil)
	}
	return s, true, nil
}

type ignoreEmptyFirst struct {
}

func (o *ignoreEmptyFirst) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" && captured == 0 && skipped == 0 {
		return "", false, nil
	}
	return s, true, nil
}

type notEmptyLast struct {
	message string
}

func (o *notEmptyLast) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" && isLast {
		return "", false, NewOptionFailError(o.message, pos, nil)
	}
	return s, true, nil
}

type ignoreEmptyLast struct {
}

func (o *ignoreEmptyLast) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" && isLast {
		return "", false, nil
	}
	return s, true, nil
}

type notEmptyInners struct {
	message string
}

func (o *notEmptyInners) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" && !isLast && !(captured == 0 && skipped == 0) {
		return "", false, NewOptionFailError(o.message, pos, nil)
	}
	return s, true, nil
}

type ignoreEmptyInners struct {
}

func (o *ignoreEmptyInners) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s != "" {
		return s, true, nil
	}
	isInner := !isLast && !(captured == 0 && skipped == 0)
	return s, !isInner, nil
}

type notEmptyOuters struct {
	message string
}

func (o *notEmptyOuters) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s == "" && (isLast || (captured == 0 && skipped == 0)) {
		return "", false, NewOptionFailError(o.message, pos, nil)
	}
	return s, true, nil
}

type ignoreEmptyOuters struct {
}

func (o *ignoreEmptyOuters) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if s != "" {
		return s, true, nil
	}
	isOuter := isLast || (captured == 0 && skipped == 0)
	return s, !isOuter, nil
}

type noContiguousQuotes struct {
	message string
}

func (o *noContiguousQuotes) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	for i := 1; i < len(subParts); i++ {
		if subParts[i].IsQuote() && subParts[i-1].IsQuote() {
			return "", false, NewOptionFailError(o.message, pos, subParts[i])
		}
	}
	return s, true, nil
}

type noMultiQuotes struct {
	message string
}

func (o *noMultiQuotes) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	qs := 0
	var second SubPart
	for _, sp := range subParts {
		if sp.IsQuote() {
			qs++
			if second == nil {
				second = sp
			}
		}
	}
	if qs > 1 {
		return "", false, NewOptionFailError(o.message, pos, second)
	}
	return s, true, nil
}

type noMultis struct {
	message string
}

func (o *noMultis) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if len(subParts) > 1 {
		return "", false, NewOptionFailError(o.message, pos, subParts[1])
	}
	return s, true, nil
}

type stripQuotes struct {
}

func (o *stripQuotes) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if len(subParts) == 1 {
		if subParts[0].IsQuote() {
			str := subParts[0].String()
			return str[1 : len(str)-1], true, nil
		}
		return s, true, nil
	}
	var sb strings.Builder
	sb.Grow(len(s))
	for _, sub := range subParts {
		str := sub.String()
		if sub.IsQuote() {
			str = str[1 : len(str)-1]
		}
		sb.WriteString(str)
	}
	return sb.String(), true, nil
}

type unescapeQuotes struct {
}

func (o *unescapeQuotes) Apply(s string, pos int, totalLen int, captured int, skipped int, isLast bool, subParts ...SubPart) (string, bool, error) {
	if len(subParts) == 1 && subParts[0].IsQuote() {
		return subParts[0].UnEscaped(), true, nil
	} else if len(subParts) == 1 {
		return s, true, nil
	}
	var sb strings.Builder
	sb.Grow(len(s))
	for _, sub := range subParts {
		str := sub.String()
		if sub.IsQuote() {
			str = sub.UnEscaped()
		}
		sb.WriteString(str)
	}
	return sb.String(), true, nil
}
