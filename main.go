package main

import (
	"github.com/koketama/koketama/pkg"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	root := &cobra.Command{
		Use:   "koketama",
		Short: "some tools used in develop",
	}

	root.AddCommand(
		pkg.MustPKGCmd(logger),
	)

	if err := root.Execute(); err != nil {
		logger.Fatal("root command execute err", zap.Error(err))
	}
}
