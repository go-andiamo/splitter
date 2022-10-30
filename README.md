# Splitter
[![GoDoc](https://godoc.org/github.com/go-andiamo/splitter?status.svg)](https://pkg.go.dev/github.com/go-andiamo/splitter)
[![Latest Version](https://img.shields.io/github/v/tag/go-andiamo/splitter.svg?sort=semver&style=flat&label=version&color=blue)](https://github.com/go-andiamo/splitter/releases)
[![codecov](https://codecov.io/gh/go-andiamo/splitter/branch/main/graph/badge.svg?token=igjnZdgh0e)](https://codecov.io/gh/go-andiamo/splitter)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-andiamo/splitter)](https://goreportcard.com/report/github.com/go-andiamo/splitter)

## Overview

Go package for splitting strings (aware of enclosing braces and quotes)

The problem with standard Golang `strings.Split` is that it does not take into consideration that the string being split may
contain enclosing braces and/or quotes (where the separator should not be considered where it's inside braces or quotes)

Take for example a string representing a slice of comma separated strings...
```go
    str := `"aaa","bbb","this, for sanity, should not be split"`
```
running `strings.Split` on that...
```go
package main

import "strings"

func main() {
    str := `"aaa","bbb","this, for sanity, should not be parts"`
    parts := strings.Split(str, `,`)
    println(len(parts))
}
```
would yield 5 ([try on go-playground](https://go.dev/play/p/bEnwjc-gfQS)) - instead of the desired 3

However, with splitter, the result would be different...
```go
package main

import "github.com/go-andiamo/splitter"

func main() {
    commaSplitter, _ := splitter.NewSplitter(',', splitter.DoubleQuotes)

    str := `"aaa","bbb","this, for sanity, should not be split"`
    parts, _ := commaSplitter.Split(str)
    println(len(parts))
}
```
which yields the desired 3! [try on go-playground](https://go.dev/play/p/lIae-RjzSe6)

Note: The varargs, after the first separator arg, are the desired 'enclosures' (e.g. quotes, brackets, etc.) to be taken
into consideration

While splitting, any enclosures specified are checked for balancing!

## Installation
To install Splitter, use go get:

    go get github.com/go-andiamo/splitter

To update Splitter to the latest version, run:

    go get -u github.com/go-andiamo/splitter

## Enclosures
Enclosures instruct the splitter specific start/end sequences within which the separator is not to be considered.  An enclosure can be one of two types: quotes or brackets.

Quote type enclosures only differ from bracket type enclosures in the way that their optional escaping works -
* Quote enclosures can be:
  * escaped by escape prefix - e.g. a quote enclosure starting with `"` and ending with `"` but `\"` is not seen as ending
  * escaped by doubles - e.g. a quote enclosure starting with `'` and ending with `'` but any doubles `''` are not seen as ending 
* Bracket enclosures can only be:
    * escaped by escape prefix - e.g. a bracket enclosure starting with `(` and ending with `)` and escape set to <code>&#92;</code>
      * `\(` is not seen as a start
      * `\)` is not seen as an end

Note that brackets are ignored inside quotes - but quotes can exist within brackets.  And when splitting, separators found within any specified quote or bracket enclosure are not considered. 


The Splitter provides many pre-defined enclosures:
<table>
    <tr>
        <th>
            Var Name
        </th>
        <th>
            Type
        </th>
        <th>
            Start - End
        </th>
        <th>
            Escaped end
        </th>
    </tr>
    <tr>
        <td><code>DoubleQuotes</code></td>
        <td>Quote</td>
        <td><code>"</code> <code>"</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>DoubleQuotesBackSlashEscaped</code></td>
        <td>Quote</td>
        <td><code>"</code> <code>"</code></td>
        <td><code>\"</code></td>
    </tr>
    <tr>
        <td><code>DoubleQuotesDoubleEscaped</code></td>
        <td>Quote</td>
        <td><code>"</code> <code>"</code></td>
        <td><code>""</code></td>
    </tr>
    <tr>
        <td><code>SingleQuotes</code></td>
        <td>Quote</td>
        <td><code>'</code> <code>'</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SingleQuotesBackSlashEscaped</code></td>
        <td>Quote</td>
        <td><code>'</code> <code>'</code></td>
        <td><code>\'</code></td>
    </tr>
    <tr>
        <td><code>SingleQuotesDoubleEscaped</code></td>
        <td>Quote</td>
        <td><code>'</code> <code>'</code></td>
        <td><code>''</code></td>
    </tr>
    <tr>
        <td><code>SingleInvertedQuotes</code></td>
        <td>Quote</td>
        <td><code>`</code> <code>`</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SingleInvertedQuotesBackSlashEscaped</code></td>
        <td>Quote</td>
        <td><code>`</code> <code>`</code></td>
        <td><code>\'</code></td>
    </tr>
    <tr>
        <td><code>SingleInvertedQuotesDoubleEscaped</code></td>
        <td>Quote</td>
        <td><code>`</code> <code>`</code></td>
        <td><code>``</code></td>
    </tr>
    <tr>
        <td><code>SinglePointingAngleQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2039;</code> <code>&#x203A;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SinglePointingAngleQuotesBackSlashEscaped</code></td>
        <td>Quote</td>
        <td><code>&#x2039;</code> <code>&#x203A;</code></td>
        <td><code>\&#x203A;</code></td>
    </tr>
    <tr>
        <td><code>DoublePointingAngleQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x00AB;</code> <code>&#x00BB;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LeftRightDoubleDoubleQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x201C;</code> <code>&#x201D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LeftRightDoubleSingleQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2018;</code> <code>&#x2019;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LeftRightDoublePrimeQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x301D;</code> <code>&#x301E;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SingleLowHigh9Quotes</code></td>
        <td>Quote</td>
        <td><code>&#x201A;</code> <code>&#x201B;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>DoubleLowHigh9Quotes</code></td>
        <td>Quote</td>
        <td><code>&#x201E;</code> <code>&#x201F;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>Parenthesis</code></td>
        <td>Brackets</td>
        <td><code>(</code> <code>)</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>CurlyBrackets</code></td>
        <td>Brackets</td>
        <td><code>{</code> <code>}</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SquareBrackets</code></td>
        <td>Brackets</td>
        <td><code>[</code> <code>]</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LtGtAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&lt;</code> <code>&gt;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LeftRightPointingAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2329;</code> <code>&#x232A;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SubscriptParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x208D;</code> <code>&#x208E;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SuperscriptParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x207d;</code> <code>&#x207e;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SmallParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#xFE59;</code> <code>&#xFE5A;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SmallCurlyBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#xFE5B;</code> <code>&#xFE5C;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>DoubleParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x2E28;</code> <code>&#x2E29;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MathWhiteSquareBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x27E6;</code> <code>&#x27E7;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MathAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x27E8;</code> <code>&#x27E9;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MathDoubleAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x27EA;</code> <code>&#x27EB;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MathWhiteTortoiseShellBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x27EC;</code> <code>&#x27ED;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MathFlattenedParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x27EE;</code> <code>&#x27EF;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>OrnateParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#xFD3E;</code> <code>&#xFD3F;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>AngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x3008;</code> <code>&#x3009;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>DoubleAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x300A;</code> <code>&#x300B;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>FullWidthParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#xFF08;</code> <code>&#xFF09;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>FullWidthSquareBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#xFF3B;</code> <code>&#xFF3D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>FullWidthCurlyBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#xFF5B;</code> <code>&#xFF5D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SubstitutionBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2E02;</code> <code>&#x2E03;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SubstitutionQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2E02;</code> <code>&#x2E03;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>DottedSubstitutionBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2E04;</code> <code>&#x2E05;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>DottedSubstitutionQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2E04;</code> <code>&#x2E05;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>TranspositionBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2E09;</code> <code>&#x2E0A;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>TranspositionQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2E09;</code> <code>&#x2E0A;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>RaisedOmissionBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2E0C;</code> <code>&#x2E0D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>RaisedOmissionQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2E0C;</code> <code>&#x2E0D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LowParaphraseBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2E1C;</code> <code>&#x2E1D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LowParaphraseQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x2E1C;</code> <code>&#x2E1D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SquareWithQuillBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2045;</code> <code>&#x2046;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>WhiteParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x2985;</code> <code>&#x2986;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>WhiteCurlyBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2983;</code> <code>&#x2984;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>WhiteSquareBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x301A;</code> <code>&#x301B;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>WhiteLenticularBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x3016;</code> <code>&#x3017;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>WhiteTortoiseShellBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x3018;</code> <code>&#x3019;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>FullWidthWhiteParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#xFF5F;</code> <code>&#xFF60;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>BlackTortoiseShellBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2997;</code> <code>&#x2998;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>BlackLenticularBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x3010;</code> <code>&#x3011;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>PointingCurvedAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x29FC;</code> <code>&#x29FD;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>TortoiseShellBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x3014;</code> <code>&#x3015;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>SmallTortoiseShellBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#xFE5D;</code> <code>&#xFE5E;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>ZNotationImageBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2987;</code> <code>&#x2988;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>ZNotationBindingBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2989;</code> <code>&#x298A;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MediumOrnamentalParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x2768;</code> <code>&#x2769;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>LightOrnamentalTortoiseShellBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2772;</code> <code>&#x2773;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MediumOrnamentalFlattenedParenthesis</code></td>
        <td>Brackets</td>
        <td><code>&#x276A;</code> <code>&#x276B;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MediumOrnamentalPointingAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x276C;</code> <code>&#x276D;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>MediumOrnamentalCurlyBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2774;</code> <code>&#x2775;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>HeavyOrnamentalPointingAngleQuotes</code></td>
        <td>Quote</td>
        <td><code>&#x276E;</code> <code>&#x276F;</code></td>
        <td><em>none</em></td>
    </tr>
    <tr>
        <td><code>HeavyOrnamentalPointingAngleBrackets</code></td>
        <td>Brackets</td>
        <td><code>&#x2770;</code> <code>&#x2771;</code></td>
        <td><em>none</em></td>
    </tr>
</table>

_Note: To convert any of the above enclosures to escaping - use the `MakeEscapable()` or `MustMakeEscapable()` functions._

### Quote enclosures with escaping
Quotes within quotes can be handled by using an enclosure that specifies how the escaping works, for example the following uses \ (backslash) prefixed escaping...
```go
package main

import "github.com/go-andiamo/splitter"

func main() {
    commaSplitter, _ := splitter.NewSplitter(',', splitter.DoubleQuotesBackSlashEscaped)

    str := `"aaa","bbb","this, for sanity, \"should\" not be split"`
    parts, _ := commaSplitter.Split(str)
    println(len(parts))
}
```
[try on go-playground](https://go.dev/play/p/wgJ68hXBp1n)

Or with double escaping...
```go
package main

import "github.com/go-andiamo/splitter"

func main() {
    commaSplitter, _ := splitter.NewSplitter(',', splitter.DoubleQuotesDoubleEscaped)

    str := `"aaa","bbb","this, for sanity, """"should,,,,"" not be split"`
    parts, _ := commaSplitter.Split(str)
    println(len(parts))
}
```
[try on go-playground](https://go.dev/play/p/3BpayDZyaA7)

#### Not separating when separator encountered in quotes or brackets...
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    encs := []*splitter.Enclosure{
        splitter.Parenthesis, splitter.SquareBrackets, splitter.CurlyBrackets,
        splitter.DoubleQuotesDoubleEscaped, splitter.SingleQuotesDoubleEscaped,
    }
    commaSplitter, _ := splitter.NewSplitter(',', encs...)

    str := `do(not,)split,'don''t,split,this',[,{,(a,"this has "" quotes")}]`
    parts, _ := commaSplitter.Split(str)
    println(len(parts))
    for i, pt := range parts {
        fmt.Printf("\t[%d]%s\n", i, pt)
    }
}
```
[try on go-playground](https://go.dev/play/p/bvzC1NXfG3z)

## Options
Options define behaviours that are to be carried out on each found part during splitting.

An option, by virtue of it's return args from `.Apply()`, can do one of three things:
1. return a modified string of what is to be added to the split parts
2. return a `false` to indicate that the split part is not to be added to the split result
3. return an `error` to indicate that the split part is unacceptable (and cease further splitting - the error is returned from the `Split` method)

Options can be added directly to the Splitter using `.AddDefaultOptions()` method.  These options are checked for every call to the splitters `.Split()` method.

Options can also be specified when calling the splitter `.Split()` method - these options are only carried out for this call (and after any options already specified on the splitter)

### Option Examples
#### 1. Stripping empty parts
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    s := splitter.MustCreateSplitter('/').
        AddDefaultOptions(splitter.IgnoreEmpties)

    parts, _ := s.Split(`/a//c/`)
    println(len(parts))
    fmt.Printf("%+v", parts)
}
```
[try on go-playground](https://go.dev/play/p/l1YnMoeA9Jm)

#### 2. Stripping empty first/last parts
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    s := splitter.MustCreateSplitter('/').
        AddDefaultOptions(splitter.IgnoreEmptyFirst, splitter.IgnoreEmptyLast)

    parts, _ := s.Split(`/a//c/`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)

    parts, _ = s.Split(`a//c/`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)

    parts, _ = s.Split(`/a//c`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)
}
```
[try on go-playground](https://go.dev/play/p/n1NEKQhtWsY)

#### 3. Trimming parts
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    s := splitter.MustCreateSplitter('/').
        AddDefaultOptions(splitter.TrimSpaces)

    parts, _ := s.Split(`/a/b/c/`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)

    parts, _ = s.Split(`  / a /b / c/    `)
    println(len(parts))
    fmt.Printf("%+v\n", parts)

    parts, _ = s.Split(`/   a   /   b   /   c   /`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)
}
```
[try on go-playground](https://go.dev/play/p/d8FZXJCBPze)

#### 4. Trimming spaces (and removing empties)
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    s := splitter.MustCreateSplitter('/').
        AddDefaultOptions(splitter.TrimSpaces, splitter.IgnoreEmpties)

    parts, _ := s.Split(`/a/  /c/`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)

    parts, _ = s.Split(`  / a // c/    `)
    println(len(parts))
    fmt.Printf("%+v\n", parts)

    parts, _ = s.Split(`/   a   /      /   c   /`)
    println(len(parts))
    fmt.Printf("%+v\n", parts)
}
```
[try on go-playground](https://go.dev/play/p/S_ald78xtSi)

#### 5. Error for empties found
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    s := splitter.MustCreateSplitter('/').
        AddDefaultOptions(splitter.TrimSpaces, splitter.NoEmpties)

    if parts, err := s.Split(`/a/  /c/`); err != nil {
        println(err.Error())
    } else {
        println(len(parts))
        fmt.Printf("%+v\n", parts)
    }

    if parts, err := s.Split(`  / a // c/    `); err != nil {
        println(err.Error())
    } else {
        println(len(parts))
        fmt.Printf("%+v\n", parts)
    }

    if parts, err := s.Split(`/   a   /      /   c   /`); err != nil {
        println(err.Error())
    } else {
        println(len(parts))
        fmt.Printf("%+v\n", parts)
    }

    if parts, err := s.Split(` a / b/c `); err != nil {
        println(err.Error())
    } else {
        println(len(parts))
        fmt.Printf("%+v\n", parts)
    }
}
```
[try on go-playground](https://go.dev/play/p/LVLkuRMoYJX)