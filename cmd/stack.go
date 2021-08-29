package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newStackCmd())
}

func newStackCmd() *cobra.Command {
	// stackコマンドをルートコマンドに追加
	cmd := &cobra.Command{
		Use:   "stack",
		Short: "Manage Stack resorces",
		// 実行する関数
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// cmdは木構造になっている??
	cmd.AddCommand(
		// ルートコマンドの子供の子供に新たにコマンドを追加
		newStackShowCmd(),
	)

	return cmd
}

func newStackShowCmd() *cobra.Command {
	// ルートコマンドの子供の子供にshowコマンドを追加
	// <StackID>はどのように取得する??
	cmd := &cobra.Command{
		// 例) hoge stack show 1
		Use:   "show <StackID>",
		Short: "Show Stack",
		// RunEにしないとerrorが返せない
		// 以下のように実行する関数を切り出すことで可読性が上がりそう
		RunE: runStackShowCmd,
	}

	return cmd
}

// 切り出した関数
// 引数にはcmdとargsの２種類を取ることに注意(切り出そうが切り出すまいが)
func runStackShowCmd(cmd *cobra.Command, args []string) error {
	client, err := newDefaultClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("StackID is required")
	}

	// primary key of stack
	appStackID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to parse StackID: %s", args[0])
	}

	// golangではRequest構造体をよく使用するイメージ
	req := AppStackShowRequest{
		ID: appStackID,
	}

	// 10秒でタイムアウトするコンテキストを作る
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.StackShow(ctx, req)
	if err != nil {
		// 受け取ったエラーにコメントを追加して、呼び出し元にエラーを伝搬させる
		return errors.Wrapf(err, "StackShow was failed: req = %+v, res = %+v", req, res)
	}

	appStack := res.AppStack
	fmt.Printf(
		"id: %d, name: %s, inserted_at: %v, updated_at: %v\n",
		appStack.ID, appStack.Name, appStack.InsertedAt, appStack.UpdatedAt,
	)

	return nil
}

// 以下はclientに定義すべき??
func (client *Client) StackShow(ctx context.Context, apiRequest AppStackShowRequest) (*AppStackShowResponse, error) {
	subPath := fmt.Sprintf("/app_stacks/%d", apiRequest.ID)
	httpRequest, err := client.newRequest(ctx, "GET", subPath, nil)
	if err != nil {
		return nil, err
	}

	httpResponse, err := client.HTTPClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	var apiResponse AppStackShowResponse
	if err := decodeBody(httpResponse, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
