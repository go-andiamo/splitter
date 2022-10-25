// Package splitter - Go package for splitting strings (aware of enclosing braces and quotes)
/*
The problem with strings.Split is that it does not take into consideration that the string being split may
contain enclosing braces and/or quotes (where the separator should not be considered where it's inside braces or quotes)

Take for example a string representing a slice of comma separated strings...
	str := `"aaa","bbb","this, for sanity, should not be split"`
running strings.Split on that...
	str := `"aaa","bbb","this, for sanity, should not be split"`
	parts := strings.Split(str, `,`)
	println(len(parts))
would yield 5 - instead of the desired 3

With splitter, the result would be different...
	str := `"aaa","bbb","this, for sanity, should not be split"`
	mySplitter, _ := splitter.NewSplitter(',', splitter.DoubleQuotes)
	parts, _ := mySplitter.Split(str)
	println(len(parts))
which yields the desired 3!

Note: The varargs, after the first separator arg, are the desired 'enclosures' (e.g. quotes, brackets, etc.) to be taken
into consideration

While splitting, any enclosures specified are checked for balancing!
*/
package splitter
