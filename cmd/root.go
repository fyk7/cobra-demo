package cmd

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// ルートコマンドを定義(ルートコマンドとは"kubectl get"の"kubectl"みたいなもの)
var RootCmd = &cobra.Command{
	Use:           "hoge",
	Short:         "A hoge CLI written in Go",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	// configファイルの読み込み
	cobra.OnInitialize(initConfig)

	// config fileのパスを--configで指定できるようにする
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.hoge.yml)")
	// エンドポイントURLにデフォルト値を設定しつつ、--urlというオプションで上書きできるようにしている
	RootCmd.PersistentFlags().StringP("url", "", "https://hoge.example.com/api", "hoge endpoint URL")

	// viperのurlキーと関連付けることにより、url: http://~~/api のようにymlファイルを記述することでurlを指定できる
	// コマンドのデフォルト値 < 設定ファイル < コマンドラインオプションの優先順位
	// viper.GetString("url")で透過的にurlの値を取得可能
	viper.BindPFlag("url", RootCmd.PersistentFlags().Lookup("url"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".hoge")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

// 明示的にimportしなくても同じpackage内に定義が存在すれば問題ない??
func newDefaultClient() (*Client, error) {
	endpointURL := viper.GetString("url")
	httpClient := &http.Client{}
	userAgent := fmt.Sprintf("hoge/%s (%s)", Version, runtime.Version())
	return newClient(endpointURL, httpClient, userAgent)
}
