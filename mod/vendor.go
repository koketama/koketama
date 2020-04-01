package mod

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/koketama/koketama/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func vendorCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vendor",
		Short: "match vendor to some git history commit",
		Long: `
when remove vendor towards to gomod, 
this will try to match vendor to some
git hitory commits, got the history
git version used in go.mod.
To compare with git history commits, 
git -C XXX --hard HEAD^ will be used.
`}

	var (
		repo    string
		vendor  string
		pkgs    []string
		ignores []string
	)
	cmd.Flags().StringVar(&repo, "repo", "/repo/src/github.com/XXX", "git repo directory")
	cmd.Flags().StringVar(&vendor, "vendor", "~/go/src/XXX/vendor/github.com/XXX", "project vendor directory")
	cmd.Flags().StringArrayVar(&pkgs, "pkgs", nil, "which package(s) in vendor should be match")
	cmd.Flags().StringArrayVar(&ignores, "ignores", []string{
		".mod", ".sum"}, "ignored files, marked by extension")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		ignoredExt := make(map[string]bool)
		for _, name := range ignores {
			ignoredExt[name] = true
		}

		hash := func(dir string) string {
			var digest []byte
			if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					logger.Fatal("walk got err", zap.Error(err))
				}

				if info.IsDir() || ignoredExt[filepath.Ext(info.Name())] {
					return nil
				}

				src, err := util.ReadFile(path)
				if err != nil {
					logger.Fatal("read file err", zap.String("file", path), zap.Error(err))
				}

				hash := hmac.New(sha256.New, digest)
				hash.Write(src)
				digest = hash.Sum(nil)

				return nil
			}); err != nil {
				logger.Fatal("walk err", zap.String("dir", dir), zap.Error(err))
			}

			return hex.EncodeToString(digest)
		}

		all := hashset.New()
		wantted := make(map[string]string) // digest: pkgname
		for _, pkg := range pkgs {
			all.Add(pkg)
			wantted[hash(filepath.Join(vendor, pkg))] = pkg
		}

		var git []string // digest
		hashGit := func() {
			git = make([]string, 0, len(pkgs))
			for _, pkg := range pkgs {
				git = append(git, hash(filepath.Join(repo, pkg)))
			}
		}

		left := func(bingos []string) []interface{} {
			tmp := hashset.New()
			tmp.Add(all.Values()...)

			for _, bingo := range bingos {
				tmp.Remove(bingo)
			}
			return tmp.Values()
		}

		maxMatch := 0
		for {
			hashGit()

			var match int
			var bingos []string
			for _, digest := range git {
				if name, ok := wantted[digest]; ok {
					match++
					bingos = append(bingos, name)
				}
			}
			if match > maxMatch {
				maxMatch = match
			}
			fmt.Println(">>>>>>>", "maxMatch:", maxMatch, "match:", match, "bingos:", bingos, "left:", left(bingos))
			fmt.Println("") // pretty for see

			output, err := exec.Command("git", "-C", repo, "reset", "--hard", "HEAD^").Output()
			if err != nil {
				logger.Fatal("exec command err", zap.Error(err))
			}
			fmt.Print(string(output))
		}
	}

	return cmd
}
