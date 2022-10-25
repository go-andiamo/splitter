# Splitter
[![GoDoc](https://godoc.org/github.com/go-andiamo/splitter?status.svg)](https://pkg.go.dev/github.com/go-andiamo/splitter)
[![Latest Version](https://img.shields.io/github/v/tag/go-andiamo/splitter.svg?sort=semver&style=flat&label=version&color=blue)](https://github.com/go-andiamo/splitter/releases)

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
would yield 5 - instead of the desired 3

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
which yields the desired 3!

Note: The varargs, after the first separator arg, are the desired 'enclosures' (e.g. quotes, brackets, etc.) to be taken
into consideration

While splitting, any enclosures specified are checked for balancing!

#### Quote enclosures can also handling escaping...
Quotes within quotes can be handled by using an enclosure that specifies how the escaping works, for example the following uses `\` (backslash) prefixed escaping...
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

#### Not separating when separator encountered in quotes or brackets...
```go
package main

import (
    "fmt"
    "github.com/go-andiamo/splitter"
)

func main() {
    encs := []*splitter.Enclosure{
        splitter.Brackets, splitter.SquareBrackets, splitter.CurlyBrackets,
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
