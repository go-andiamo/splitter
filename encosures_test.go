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
