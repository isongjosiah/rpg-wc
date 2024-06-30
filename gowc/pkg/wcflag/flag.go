package wcflag

import (
	"flag"
)

// Options defines all the flag options for gowc
type Options struct {
	CountBytes      *bool
	CountLines      *bool
	CountWords      *bool
	CountCharacters *bool
}

// ParseFlag parses the flag set by the user and returns the FlagOptions
// containing the specified flags if any
func ParseFlag() *Options {

	options := Options{
		CountBytes:      new(bool),
		CountLines:      new(bool),
		CountWords:      new(bool),
		CountCharacters: new(bool),
	}

	// define the flags
	flag.BoolVar(options.CountBytes, "c", false, "The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage of the -m option")
	flag.BoolVar(options.CountLines, "l", false, "The number of lines in each input file is written to the standard output")
	flag.BoolVar(options.CountWords, "w", false, "The number of words in each input file is written to the standard output")
	flag.BoolVar(options.CountCharacters, "m", false, "The number of characters in each input file is written to the standard output")

	// parse the flags
	flag.Parse()

	return &options
}
