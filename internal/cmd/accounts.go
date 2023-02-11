package cmd

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdAccounts(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL := viper.GetString("base-url")

			sandbox := viper.GetBool("sandbox")
			if sandbox {
				baseURL = "https://api.sandbox.dnsimple.com"
			}

			accessToken := viper.GetString("access-token")

			if opts.Client == nil {
				return errors.New("invalid client builder")
			}

			client := opts.Client(baseURL, accessToken)

			resp, err := client.Accounts.ListAccounts(context.Background(), nil)
			if err != nil {
				return err
			}

			cmd.Printf("%#q\n", resp.Data)

			return nil
		},
	}

	return cmd
}
