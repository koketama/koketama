package mod

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// MustModCmd return a new mod cmd, this will panic if logger is nil
func MustModCmd(logger *zap.Logger) *cobra.Command {
	if logger == nil {
		panic("logger required")
	}

	cmd := &cobra.Command{
		Use:   "mod",
		Short: "operations on gomod and vendor",
	}

	cmd.AddCommand(
		vendorCmd(logger),
	)

	return cmd
}
