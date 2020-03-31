package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/koketama/koketama/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func replaceCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace",
		Short: "replace some imported package(s) to other(s)",
		Long:  "this will ignore vendor directory, and not .go files",
	}

	var (
		path     string
		pkgname  string
		replace2 string
	)
	cmd.Flags().StringVar(&path, "path", "~/go/src/xxx", "project root path")
	cmd.Flags().StringVar(&pkgname, "pkgname", "from", "imported package name")
	cmd.Flags().StringVar(&replace2, "replace2", "someother", "replace package name to what")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		logger.Info("flags", zap.String("path", path), zap.String("pkgname", pkgname), zap.String("replace2", replace2))
		fmt.Println("Are you sure to perform the replacement? Y/n")
		var answer string
		fmt.Scanf("%s", &answer)
		if !strings.HasPrefix(strings.ToUpper(answer), "Y") {
			return
		}

		fmt.Println("Start execution after 5 seconds")
		time.Sleep(time.Second * 5)
		fmt.Println("Starting.....")

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Fatal("walk got err", zap.Error(err))
			}

			if info.IsDir() || strings.Contains(path, "/vendor/") || filepath.Ext(path) != ".go" {
				return nil
			}

			src, err := util.ReadFile(path)
			if err != nil {
				logger.Fatal("read file err", zap.String("file", path), zap.Error(err))
			}

			file, err := parser.ParseFile(token.NewFileSet(), info.Name(), src, 0)
			if err != nil {
				logger.Fatal("parserfile err", zap.String("file", path), zap.Error(err))
			}

			var importEndPosition int
			ast.Inspect(file, func(n ast.Node) bool {
				switch spec := n.(type) {
				case *ast.ImportSpec:
					importEndPosition = int(spec.Path.End())
				}
				return true
			})

			buffer := bytes.NewBuffer(nil)

			scanner := bufio.NewScanner(bytes.NewReader(src[:importEndPosition]))
			for scanner.Scan() {
				if line := scanner.Text(); strings.Contains(line, replace2) {
					buffer.WriteString(line)
				} else {
					buffer.WriteString(strings.Replace(line, pkgname, replace2, 1))
				}
				buffer.WriteString("\n") // scanner.Text() not include LF
			}
			if err = scanner.Err(); err != nil {
				logger.Fatal("scan import content err", zap.String("file", path), zap.Error(err))
			}

			buffer.Write(src[importEndPosition:])

			if err = util.WriteFile(path, buffer.Bytes()); err != nil {
				logger.Fatal("write to file err", zap.String("file", path), zap.Error(err))
			}

			return nil
		})
		if err != nil {
			logger.Fatal("walk err", zap.Error(err))
		}

		fmt.Println("Done")
	}

	return cmd
}
