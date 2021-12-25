package annotations

import "flag"

type options struct {
	Dir     *string
	Verbose *bool
	AbortOnError *bool
	OutputToStdOutOnly *bool
	IgnoreFiles *string
}

var Option options

func Options() {
	Option.Dir = flag.String("dir", ".", "path of go source files which needs to be processed for wish. Default value is current directory.")
	Option.Verbose = flag.Bool("v", false, "verbose option. Default value is false.")
	Option.AbortOnError = flag.Bool("a", true, "Abort on error option. Default value is true.")
	Option.OutputToStdOutOnly = flag.Bool("stdout", false, "Option to specify the output. If true, then the output will be shown only on STDOUT, else it will be written to the corresponding file having name same as original source file name but suffixed with '_wish'. Default is 'false'")
	Option.IgnoreFiles = flag.String("ignore", "", "Option to specify list of comma separated file suffixes to ignore in addition to '_test.go' and '_wish.go'")
	flag.Parse()
}