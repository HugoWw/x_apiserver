package flag

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Fprintf(os.Stdout, "FLAG: --%s=%q,", flag.Name, flag.Value)
	})
}
