package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// ツールのバージョン
	Version string
	// git commitのハッシュ値を -ldflagsでbuild時に埋め込み
	Revision string
)

func init() {
	// ルートコマンドに新しいサブコマンドを追加している
	// ルートコマンド"kubectl"にサブコマンド"get"を追加するみたいなイメージ
	RootCmd.AddCommand(newVersionCmd())
}

// 要は"version"コマンドが叩かれたら、Run以下の関数が実行される
func newVersionCmd() *cobra.Command {
	// Commandインスタンスを作成して返す
	// 以下のようにbuild時に明示的にVersionとRevisionを指定する
	// $ go build -ldflags "-X github.com/minamijoyo/api-cli-go-example/cmd.Version=v0.0.1 -X github.com/minamijoyo/api-cli-go-example/cmd.Revision=2463e27" -o bin/hoge
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		// 実行する処理をRunに関数として記述
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("hoge version: %s, revision: %s\n", Version, Revision)
		},
	}

	return cmd
}
