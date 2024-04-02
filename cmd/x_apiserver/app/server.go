package app

import (
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/options"
	"github.com/HugoWw/x_apiserver/pkg/apiserver"
	cliflag "github.com/HugoWw/x_apiserver/pkg/apiserver/cli/flag"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/cli/terminal"
	"github.com/HugoWw/x_apiserver/pkg/signals"
	"github.com/spf13/cobra"
	"strings"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand() *cobra.Command {

	opt := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:  "X-APIServer",
		Long: `The X-APIServer is used to provide The services REST operations.`,

		// stop printing usage when the command errors
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			//fs := cmd.Flags()

			// output server version info
			options.PrintVersionAndExitIfRequested()

			//cliflag.PrintFlags(fs)

			// validate options
			if err := opt.Validate(); len(err) != 0 {
				return Errs(err)
			}

			return Run(opt, signals.GetStopSignal())
		},
	}

	fs := cmd.Flags()
	namedFlagSets := opt.Flags()
	options.VersionAddFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	// get terminal windows width
	terminalWidth, _, _ := terminal.TerminalSize(cmd.OutOrStdout())

	// output usage help info
	cliflag.SetUsageAndHelpFunc(cmd, namedFlagSets, terminalWidth)

	return cmd
}

func Run(completeOptions *options.ServerRunOptions, stopCh <-chan struct{}) error {

	s, err := apiserver.Create(completeOptions)
	if err != nil {
		return err
	}

	s.PrepareRun()

	return s.Run(stopCh)
}

// Errs is used by cli option Validate error
type Errs []error

func (e Errs) Error() string {

	builder := &strings.Builder{}
	for _, err := range e {
		builder.WriteString(err.Error())
		builder.WriteByte('\n')
	}
	return builder.String()
}
