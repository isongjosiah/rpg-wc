package wcflag

import (
	"flag"
)

// Options defines all the flag options for gowc
type Options struct {
	CountBytes *bool
}

// ParseFlag parses the flag set by the user and returns the FlagOptions
// containing the specified flags if any
func ParseFlag() *Options {

	options := Options{
		CountBytes: new(bool),
	}

	// define the flags
	flag.BoolVar(options.CountBytes, "c", false, "The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage of the -m option")

	// parse the flags
	flag.Parse()

	return &options
}
