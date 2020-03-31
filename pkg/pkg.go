package pkg

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// MustPKGCmd return a new pkg cmd, this will panic if logger is nil
func MustPKGCmd(logger *zap.Logger) *cobra.Command {
	if logger == nil {
		panic("logger required")
	}

	cmd := &cobra.Command{
		Use:   "pkg",
		Short: "operations on imported package(s)",
	}

	cmd.AddCommand(
		replaceCmd(logger),
	)

	return cmd
}
