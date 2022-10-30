package splitter

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var testEnclosures = map[string]*Enclosure{
	"DoubleQuotes":                              DoubleQuotes,
	"DoubleQuotesBackSlashEscaped":              DoubleQuotesBackSlashEscaped,
	"DoubleQuotesDoubleEscaped":                 DoubleQuotesDoubleEscaped,
	"SingleQuotes":                              SingleQuotes,
	"SingleQuotesBackSlashEscaped":              SingleQuotesBackSlashEscaped,
	"SingleQuotesDoubleEscaped":                 SingleQuotesDoubleEscaped,
	"SingleInvertedQuotes":                      SingleInvertedQuotes,
	"SingleInvertedQuotesBackSlashEscaped":      SingleInvertedQuotesBackSlashEscaped,
	"SingleInvertedQuotesDoubleEscaped":         SingleInvertedQuotesDoubleEscaped,
	"DoublePointingAngleQuotes":                 DoublePointingAngleQuotes,
	"SinglePointingAngleQuotes":                 SinglePointingAngleQuotes,
	"SinglePointingAngleQuotesBackSlashEscaped": SinglePointingAngleQuotesBackSlashEscaped,
	"LeftRightDoubleDoubleQuotes":               LeftRightDoubleDoubleQuotes,
	"LeftRightDoubleSingleQuotes":               LeftRightDoubleSingleQuotes,
	"LeftRightDoublePrimeQuotes":                LeftRightDoublePrimeQuotes,
	"SingleLowHigh9Quotes":                      SingleLowHigh9Quotes,
	"DoubleLowHigh9Quotes":                      DoubleLowHigh9Quotes,
	"Parenthesis":                               Parenthesis,
	"CurlyBrackets":                             CurlyBrackets,
	"SquareBrackets":                            SquareBrackets,
	"LtGtAngleBrackets":                         LtGtAngleBrackets,
	"LeftRightPointingAngleBrackets":            LeftRightPointingAngleBrackets,
	"SubscriptParenthesis":                      SubscriptParenthesis,
	"SuperscriptParenthesis":                    SuperscriptParenthesis,
	"SmallParenthesis":                          SmallParenthesis,
	"SmallCurlyBrackets":                        SmallCurlyBrackets,
	"DoubleParenthesis":                         DoubleParenthesis,
	"MathWhiteSquareBrackets":                   MathWhiteSquareBrackets,
	"MathAngleBrackets":                         MathAngleBrackets,
	"MathDoubleAngleBrackets":                   MathDoubleAngleBrackets,
	"MathWhiteTortoiseShellBrackets":            MathWhiteTortoiseShellBrackets,
	"MathFlattenedParenthesis":                  MathFlattenedParenthesis,
	"OrnateParenthesis":                         OrnateParenthesis,
	"AngleBrackets":                             AngleBrackets,
	"DoubleAngleBrackets":                       DoubleAngleBrackets,
	"FullWidthParenthesis":                      FullWidthParenthesis,
	"FullWidthSquareBrackets":                   FullWidthSquareBrackets,
	"FullWidthCurlyBrackets":                    FullWidthCurlyBrackets,
	"SubstitutionBrackets":                      SubstitutionBrackets,
	"SubstitutionQuotes":                        SubstitutionQuotes,
	"DottedSubstitutionBrackets":                DottedSubstitutionBrackets,
	"DottedSubstitutionQuotes":                  DottedSubstitutionQuotes,
	"TranspositionBrackets":                     TranspositionBrackets,
	"TranspositionQuotes":                       TranspositionQuotes,
	"RaisedOmissionBrackets":                    RaisedOmissionBrackets,
	"RaisedOmissionQuotes":                      RaisedOmissionQuotes,
	"LowParaphraseBrackets":                     LowParaphraseBrackets,
	"LowParaphraseQuotes":                       LowParaphraseQuotes,
	"SquareWithQuillBrackets":                   SquareWithQuillBrackets,
	"WhiteParenthesis":                          WhiteParenthesis,
	"WhiteCurlyBrackets":                        WhiteCurlyBrackets,
	"WhiteSquareBrackets":                       WhiteSquareBrackets,
	"WhiteLenticularBrackets":                   WhiteLenticularBrackets,
	"WhiteTortoiseShellBrackets":                WhiteTortoiseShellBrackets,
	"FullWidthWhiteParenthesis":                 FullWidthWhiteParenthesis,
	"BlackTortoiseShellBrackets":                BlackTortoiseShellBrackets,
	"BlackLenticularBrackets":                   BlackLenticularBrackets,
	"PointingCurvedAngleBrackets":               PointingCurvedAngleBrackets,
	"TortoiseShellBrackets":                     TortoiseShellBrackets,
	"SmallTortoiseShellBrackets":                SmallTortoiseShellBrackets,
	"ZNotationImageBrackets":                    ZNotationImageBrackets,
	"ZNotationBindingBrackets":                  ZNotationBindingBrackets,
	"MediumOrnamentalParenthesis":               MediumOrnamentalParenthesis,
	"LightOrnamentalTortoiseShellBrackets":      LightOrnamentalTortoiseShellBrackets,
	"MediumOrnamentalFlattenedParenthesis":      MediumOrnamentalFlattenedParenthesis,
	"MediumOrnamentalPointingAngleBrackets":     MediumOrnamentalPointingAngleBrackets,
	"MediumOrnamentalCurlyBrackets":             MediumOrnamentalCurlyBrackets,
	"HeavyOrnamentalPointingAngleQuotes":        HeavyOrnamentalPointingAngleQuotes,
	"HeavyOrnamentalPointingAngleBrackets":      HeavyOrnamentalPointingAngleBrackets,
}

func TestEnclosures(t *testing.T) {
	for name, enc := range testEnclosures {
		t.Run(fmt.Sprintf("%s", name), func(t *testing.T) {
			require.NotEqual(t, rune(0), enc.Start)
			require.NotEqual(t, rune(0), enc.End)
			if strings.Contains(name, "Quote") {
				require.True(t, enc.IsQuote)
				require.Equal(t, enc.Escape == rune(0), !enc.Escapable)
				if enc.Escapable {
					require.NotEqual(t, rune(0), enc.Escape)
				}
			} else {
				require.False(t, enc.IsQuote)
				require.False(t, enc.Escapable)
				require.Equal(t, rune(0), enc.Escape)
			}
		})
	}
}

func TestMakeEscapable(t *testing.T) {
	escpd, err := MakeEscapable(DoubleQuotes, '\\')
	require.NoError(t, err)
	require.Equal(t, DoubleQuotes.Start, escpd.Start)
	require.Equal(t, DoubleQuotes.End, escpd.End)
	require.Equal(t, '\\', escpd.Escape)
	require.NotEqual(t, escpd.Escape, DoubleQuotes.Escape)
	require.True(t, escpd.Escapable)
	require.True(t, escpd.isEscapable())
	require.False(t, escpd.isDoubleEscaping())

	escpd, err = MakeEscapable(DoubleQuotes, '"')
	require.NoError(t, err)
	require.Equal(t, DoubleQuotes.Start, escpd.Start)
	require.Equal(t, DoubleQuotes.End, escpd.End)
	require.Equal(t, '"', escpd.Escape)
	require.True(t, escpd.Escapable)
	require.True(t, escpd.isEscapable())
	require.True(t, escpd.isDoubleEscaping())

	_, err = MakeEscapable(Parenthesis, '\\')
	require.NoError(t, err)
	_, err = MakeEscapable(Parenthesis, '(')
	require.Error(t, err)
	require.Equal(t, `bracket enclosures cannot be double-escaped`, err.Error())
	_, err = MakeEscapable(Parenthesis, ')')
	require.Error(t, err)
	require.Equal(t, `bracket enclosures cannot be double-escaped`, err.Error())
}

func TestMustMakeEscapable(t *testing.T) {
	escpd := MustMakeEscapable(Parenthesis, '\\')
	require.Equal(t, Parenthesis.Start, escpd.Start)
	require.Equal(t, Parenthesis.End, escpd.End)
	require.Equal(t, '\\', escpd.Escape)
	require.True(t, escpd.Escapable)
	require.True(t, escpd.isEscapable())
	require.False(t, escpd.isDoubleEscaping())

	require.Panics(t, func() {
		MustMakeEscapable(Parenthesis, '(')
	})
	require.Panics(t, func() {
		MustMakeEscapable(Parenthesis, ')')
	})
}
