package splitter

import "errors"

// Enclosure is used when creating a splitter to denote what 'enclosures' (e.g. quotes or brackets)
// are to be considered when slitting.
//
// For example, creating a NewSplitter with separator of ',' (comma) and an Enclosure with
// Start and End of `"` (double quotes) - splitting the following...
//	str := `"aaa","dont,split,this"`
// would yield two split items - `"aaa"` and `"dont,split,this"`
type Enclosure struct {
	// the starting rune for the enclosure
	Start rune
	// the ending rune for the enclosure
	End rune
	// whether the enclosure is a quote enclosure
	IsQuote bool
	// whether the close is escapable (only for IsQuote)
	Escapable bool
	// the prefix escape rune (only for IsQuote & and used with Escapable)
	Escape rune
}

func (e *Enclosure) clone() Enclosure {
	return Enclosure{
		Start:     e.Start,
		End:       e.End,
		IsQuote:   e.IsQuote,
		Escapable: e.Escapable,
		Escape:    e.Escape,
	}
}

func (e *Enclosure) isDoubleEscaping() bool {
	return e.IsQuote && e.Escapable && e.End == e.Escape
}

func (e *Enclosure) isEscapable() bool {
	return e.Escapable
}

func (e *Enclosure) isBracketEscapable() bool {
	return !e.IsQuote && e.Escapable
}

// MakeEscapable makes an escapable copy of an enclosure
//
// returns an error if the supplied Enclosure is a brackets
// and the escape rune matches either the start or end rune
// (because brackets cannot be double-escaped - as this would prevent nested brackets)
func MakeEscapable(enc *Enclosure, esc rune) (*Enclosure, error) {
	if !enc.IsQuote && (esc == enc.Start || esc == enc.End) {
		return nil, errors.New("bracket enclosures cannot be double-escaped")
	}
	return &Enclosure{
		Start:     enc.Start,
		End:       enc.End,
		IsQuote:   enc.IsQuote,
		Escapable: true,
		Escape:    esc,
	}, nil
}

// MustMakeEscapable is the same as MakeEscapable, except that it panics on error
func MustMakeEscapable(enc *Enclosure, esc rune) *Enclosure {
	if result, err := MakeEscapable(enc, esc); err != nil {
		panic(err)
	} else {
		return result
	}
}

const (
	escBackslash = '\\'
)

var (
	DoubleQuotes                              = _DoubleQuotes
	DoubleQuotesBackSlashEscaped              = MustMakeEscapable(_DoubleQuotes, escBackslash)
	DoubleQuotesDoubleEscaped                 = MustMakeEscapable(_DoubleQuotes, '"')
	SingleQuotes                              = _SingleQuotes
	SingleQuotesBackSlashEscaped              = MustMakeEscapable(_SingleQuotes, escBackslash)
	SingleQuotesDoubleEscaped                 = MustMakeEscapable(_SingleQuotes, '\'')
	SingleInvertedQuotes                      = _SingleInvertedQuotes
	SingleInvertedQuotesBackSlashEscaped      = MustMakeEscapable(_SingleInvertedQuotes, escBackslash)
	SingleInvertedQuotesDoubleEscaped         = MustMakeEscapable(_SingleInvertedQuotes, '`')
	DoublePointingAngleQuotes                 = _DoublePointingAngleQuotes
	SinglePointingAngleQuotes                 = _SinglePointingAngleQuotes
	SinglePointingAngleQuotesBackSlashEscaped = MustMakeEscapable(_SinglePointingAngleQuotes, escBackslash)
	LeftRightDoubleDoubleQuotes               = _LeftRightDoubleDoubleQuotes
	LeftRightDoubleSingleQuotes               = _LeftRightDoubleSingleQuotes
	LeftRightDoublePrimeQuotes                = _LeftRightDoublePrimeQuotes
	SingleLowHigh9Quotes                      = _SingleLowHigh9Quotes
	DoubleLowHigh9Quotes                      = _DoubleLowHigh9Quotes
	Parenthesis                               = _Parenthesis
	CurlyBrackets                             = _CurlyBrackets
	SquareBrackets                            = _SquareBrackets
	LtGtAngleBrackets                         = _LtGtAngleBrackets
	LeftRightPointingAngleBrackets            = _LeftRightPointingAngleBrackets
	SubscriptParenthesis                      = _SubscriptParenthesis
	SuperscriptParenthesis                    = _SuperscriptParenthesis
	SmallParenthesis                          = _SmallParenthesis
	SmallCurlyBrackets                        = _SmallCurlyBrackets
	DoubleParenthesis                         = _DoubleParenthesis
	MathWhiteSquareBrackets                   = _MathWhiteSquareBrackets
	MathAngleBrackets                         = _MathAngleBrackets
	MathDoubleAngleBrackets                   = _MathDoubleAngleBrackets
	MathWhiteTortoiseShellBrackets            = _MathWhiteTortoiseShellBrackets
	MathFlattenedParenthesis                  = _MathFlattenedParenthesis
	OrnateParenthesis                         = _OrnateParenthesis
	AngleBrackets                             = _AngleBrackets
	DoubleAngleBrackets                       = _DoubleAngleBrackets
	FullWidthParenthesis                      = _FullWidthParenthesis
	FullWidthSquareBrackets                   = _FullWidthSquareBrackets
	FullWidthCurlyBrackets                    = _FullWidthCurlyBrackets
	SubstitutionBrackets                      = _SubstitutionBrackets
	SubstitutionQuotes                        = _SubstitutionQuotes
	DottedSubstitutionBrackets                = _DottedSubstitutionBrackets
	DottedSubstitutionQuotes                  = _DottedSubstitutionQuotes
	TranspositionBrackets                     = _TranspositionBrackets
	TranspositionQuotes                       = _TranspositionQuotes
	RaisedOmissionBrackets                    = _RaisedOmissionBrackets
	RaisedOmissionQuotes                      = _RaisedOmissionQuotes
	LowParaphraseBrackets                     = _LowParaphraseBrackets
	LowParaphraseQuotes                       = _LowParaphraseQuotes
	SquareWithQuillBrackets                   = _SquareWithQuillBrackets
	WhiteParenthesis                          = _WhiteParenthesis
	WhiteCurlyBrackets                        = _WhiteCurlyBrackets
	WhiteSquareBrackets                       = _WhiteSquareBrackets
	WhiteLenticularBrackets                   = _WhiteLenticularBrackets
	WhiteTortoiseShellBrackets                = _WhiteTortoiseShellBrackets
	FullWidthWhiteParenthesis                 = _FullWidthWhiteParenthesis
	BlackTortoiseShellBrackets                = _BlackTortoiseShellBrackets
	BlackLenticularBrackets                   = _BlackLenticularBrackets
	PointingCurvedAngleBrackets               = _PointingCurvedAngleBrackets
	TortoiseShellBrackets                     = _TortoiseShellBrackets
	SmallTortoiseShellBrackets                = _SmallTortoiseShellBrackets
	ZNotationImageBrackets                    = _ZNotationImageBrackets
	ZNotationBindingBrackets                  = _ZNotationBindingBrackets
	MediumOrnamentalParenthesis               = _MediumOrnamentalParenthesis
	LightOrnamentalTortoiseShellBrackets      = _LightOrnamentalTortoiseShellBrackets
	MediumOrnamentalFlattenedParenthesis      = _MediumOrnamentalFlattenedParenthesis
	MediumOrnamentalPointingAngleBrackets     = _MediumOrnamentalPointingAngleBrackets
	MediumOrnamentalCurlyBrackets             = _MediumOrnamentalCurlyBrackets
	HeavyOrnamentalPointingAngleQuotes        = _HeavyOrnamentalPointingAngleQuotes
	HeavyOrnamentalPointingAngleBrackets      = _HeavyOrnamentalPointingAngleBrackets
)

var (
	_DoubleQuotes = &Enclosure{
		Start:   '"',
		End:     '"',
		IsQuote: true,
	}
	_SingleQuotes = &Enclosure{
		Start:   '\'',
		End:     '\'',
		IsQuote: true,
	}
	_SingleInvertedQuotes = &Enclosure{
		Start:   '`',
		End:     '`',
		IsQuote: true,
	}
	_SinglePointingAngleQuotes = &Enclosure{
		Start:   '\u2039',
		End:     '\u203A',
		IsQuote: true,
	}
	_DoublePointingAngleQuotes = &Enclosure{
		Start:   '\u00AB',
		End:     '\u00BB',
		IsQuote: true,
	}
	_LeftRightDoubleDoubleQuotes = &Enclosure{
		Start:   '\u201C',
		End:     '\u201D',
		IsQuote: true,
	}
	_LeftRightDoubleSingleQuotes = &Enclosure{
		Start:   '\u2018',
		End:     '\u2019',
		IsQuote: true,
	}
	_LeftRightDoublePrimeQuotes = &Enclosure{
		Start:   '\u301D',
		End:     '\u301E',
		IsQuote: true,
	}
	_SingleLowHigh9Quotes = &Enclosure{
		Start:   '\u201A', // ‚
		End:     '\u201B', // ‛
		IsQuote: true,
	}
	_DoubleLowHigh9Quotes = &Enclosure{
		Start:   '\u201E', // „
		End:     '\u201F', // ‟
		IsQuote: true,
	}
	_Parenthesis = &Enclosure{
		Start: '(',
		End:   ')',
	}
	_CurlyBrackets = &Enclosure{
		Start: '{',
		End:   '}',
	}
	_SquareBrackets = &Enclosure{
		Start: '[',
		End:   ']',
	}
	_LtGtAngleBrackets = &Enclosure{
		Start: '<',
		End:   '>',
	}
	_LeftRightPointingAngleBrackets = &Enclosure{
		Start: '\u2329',
		End:   '\u232A',
	}
	_SubscriptParenthesis = &Enclosure{
		Start: '\u208D',
		End:   '\u208E',
	}
	_SuperscriptParenthesis = &Enclosure{
		Start: '\u207d',
		End:   '\u207e',
	}
	_SmallParenthesis = &Enclosure{
		Start: '\uFE59',
		End:   '\uFE5A',
	}
	_SmallCurlyBrackets = &Enclosure{
		Start: '\uFE5B',
		End:   '\uFE5C',
	}
	_DoubleParenthesis = &Enclosure{
		Start: '\u2E28',
		End:   '\u2E29',
	}
	_MathWhiteSquareBrackets = &Enclosure{
		Start: '\u27E6',
		End:   '\u27E7',
	}
	_MathAngleBrackets = &Enclosure{
		Start: '\u27E8',
		End:   '\u27E9',
	}
	_MathDoubleAngleBrackets = &Enclosure{
		Start: '\u27EA',
		End:   '\u27EB',
	}
	_MathWhiteTortoiseShellBrackets = &Enclosure{
		Start: '\u27EC',
		End:   '\u27ED',
	}
	_MathFlattenedParenthesis = &Enclosure{
		Start: '\u27EE',
		End:   '\u27EF',
	}
	_OrnateParenthesis = &Enclosure{
		Start: '\uFD3E',
		End:   '\uFD3F',
	}
	_AngleBrackets = &Enclosure{
		Start: '\u3008',
		End:   '\u3009',
	}
	_DoubleAngleBrackets = &Enclosure{
		Start: '\u300A',
		End:   '\u300B',
	}
	_FullWidthParenthesis = &Enclosure{
		Start: '\uFF08',
		End:   '\uFF09',
	}
	_FullWidthSquareBrackets = &Enclosure{
		Start: '\uFF3B',
		End:   '\uFF3D',
	}
	_FullWidthCurlyBrackets = &Enclosure{
		Start: '\uFF5B',
		End:   '\uFF5D',
	}
	_SubstitutionBrackets = &Enclosure{
		Start: '\u2E02',
		End:   '\u2E03',
	}
	_SubstitutionQuotes = &Enclosure{
		Start:   '\u2E02',
		End:     '\u2E03',
		IsQuote: true,
	}
	_DottedSubstitutionBrackets = &Enclosure{
		Start: '\u2E04',
		End:   '\u2E05',
	}
	_DottedSubstitutionQuotes = &Enclosure{
		Start:   '\u2E04',
		End:     '\u2E05',
		IsQuote: true,
	}
	_TranspositionBrackets = &Enclosure{
		Start: '\u2E09',
		End:   '\u2E0A',
	}
	_TranspositionQuotes = &Enclosure{
		Start:   '\u2E09',
		End:     '\u2E0A',
		IsQuote: true,
	}
	_RaisedOmissionBrackets = &Enclosure{
		Start: '\u2E0C',
		End:   '\u2E0D',
	}
	_RaisedOmissionQuotes = &Enclosure{
		Start:   '\u2E0C',
		End:     '\u2E0D',
		IsQuote: true,
	}
	_LowParaphraseBrackets = &Enclosure{
		Start: '\u2E1C',
		End:   '\u2E1D',
	}
	_LowParaphraseQuotes = &Enclosure{
		Start:   '\u2E1C',
		End:     '\u2E1D',
		IsQuote: true,
	}
	_SquareWithQuillBrackets = &Enclosure{
		Start: '\u2045',
		End:   '\u2046',
	}
	_WhiteParenthesis = &Enclosure{
		Start: '\u2985',
		End:   '\u2986',
	}
	_WhiteCurlyBrackets = &Enclosure{
		Start: '\u2983',
		End:   '\u2984',
	}
	_WhiteSquareBrackets = &Enclosure{
		Start: '\u301A',
		End:   '\u301B',
	}
	_WhiteLenticularBrackets = &Enclosure{
		Start: '\u3016',
		End:   '\u3017',
	}
	_WhiteTortoiseShellBrackets = &Enclosure{
		Start: '\u3018',
		End:   '\u3019',
	}
	_FullWidthWhiteParenthesis = &Enclosure{
		Start: '\uFF5F',
		End:   '\uFF60',
	}
	_BlackTortoiseShellBrackets = &Enclosure{
		Start: '\u2997',
		End:   '\u2998',
	}
	_BlackLenticularBrackets = &Enclosure{
		Start: '\u3010',
		End:   '\u3011',
	}
	_PointingCurvedAngleBrackets = &Enclosure{
		Start: '\u29FC',
		End:   '\u29FD',
	}
	_TortoiseShellBrackets = &Enclosure{
		Start: '\u3014',
		End:   '\u3015',
	}
	_SmallTortoiseShellBrackets = &Enclosure{
		Start: '\uFE5D',
		End:   '\uFE5E',
	}
	_ZNotationImageBrackets = &Enclosure{
		Start: '\u2987',
		End:   '\u2988',
	}
	_ZNotationBindingBrackets = &Enclosure{
		Start: '\u2989',
		End:   '\u298A',
	}
	_MediumOrnamentalParenthesis = &Enclosure{
		Start: '\u2768',
		End:   '\u2769',
	}
	_LightOrnamentalTortoiseShellBrackets = &Enclosure{
		Start: '\u2772',
		End:   '\u2773',
	}
	_MediumOrnamentalFlattenedParenthesis = &Enclosure{
		Start: '\u276A',
		End:   '\u276B',
	}
	_MediumOrnamentalPointingAngleBrackets = &Enclosure{
		Start: '\u276C',
		End:   '\u276D',
	}
	_MediumOrnamentalCurlyBrackets = &Enclosure{
		Start: '\u2774',
		End:   '\u2775',
	}
	_HeavyOrnamentalPointingAngleQuotes = &Enclosure{
		Start:   '\u276E',
		End:     '\u276F',
		IsQuote: true,
	}
	_HeavyOrnamentalPointingAngleBrackets = &Enclosure{
		Start: '\u2770',
		End:   '\u2771',
	}
)
