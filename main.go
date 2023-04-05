package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "wiki-toc-generator [wikiRepoURL]",
		Short: "wiki-toc-generator is a tool to clone a GitLab Wiki repository and generate a table of contents",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			wikiRepoURL := args[0]
			// ioutil.TempDirの代わりにos.MkdirTempを使用
			cloneDir, err := os.MkdirTemp("", "cloned_wiki")
			if err != nil {
				fmt.Printf("Error creating temporary directory: %v\n", err)
				return
			}
			defer os.RemoveAll(cloneDir)

			err = cloneWiki(wikiRepoURL, cloneDir)
			if err != nil {
				fmt.Printf("Error cloning wiki: %v\n", err)
				return
			}

			toc, err := generateTOC(cloneDir)
			if err != nil {
				fmt.Printf("Error generating TOC: %v\n", err)
				return
			}

			fmt.Println("Table of Contents:")
			fmt.Println(toc)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// cloneWiki は、Gitリポジトリをクローンするための関数です。
func cloneWiki(repoURL, destDir string) error {
	cmd := exec.Command("git", "clone", repoURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// generateTOC は、クローンされたディレクトリから目次を生成するための関数です。
// この関数は、与えられたディレクトリ内のすべてのマークダウンファイルを探索し、
// それらの相対パスを使用して整形された目次を作成します。
// さらに、目次のエントリは、ファイルの階層に応じてインデントされます。
// これにより、目次が視覚的に整理され、ファイルの構造が明確になります。
// 最後に、生成された目次は、Markdownリンクの形式で返されます。
func generateTOC(dir string) (string, error) {
	// ファイルの相対パスを格納するスライス
	var paths []string

	// ディレクトリを再帰的に探索し、.mdファイルの相対パスを取得
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if filepath.Ext(relPath) != ".md" {
			return nil
		}

		paths = append(paths, relPath)
		return nil
	})
	if err != nil {
		return "", err
	}

	// パスを辞書順にソート
	sort.Strings(paths)

	// 目次を作成するための strings.Builder を初期化
	var toc strings.Builder

	// パスのスライスをイテレートし、目次のエントリを追加
	for _, relPath := range paths {
		depth := strings.Count(relPath, string(os.PathSeparator))
		indent := ""
		if depth > 0 {
			indent = strings.Repeat("  ", depth-1)
		}
		title := strings.TrimSuffix(filepath.Base(relPath), ".md")
		relPathWithoutExt := strings.TrimSuffix(relPath, ".md")
		toc.WriteString(fmt.Sprintf("%s* [[%s|%s]]\n", indent, title, relPathWithoutExt))
	}

	// 生成された目次を返す
	return toc.String(), nil
}
