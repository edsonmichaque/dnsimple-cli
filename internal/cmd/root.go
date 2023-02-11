package cmd

import (
	"fmt"
	"github.com/edsonmichaque/dnsimple-cli/internal/build"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	configFile string
	profile    string
)

func NewCmdRoot(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dnsimple",
		Version: build.Version,
	}

	cobra.OnInitialize(lookupConfigFiles)

	cmd.AddGroup(&cobra.Group{
		ID:    "domains",
		Title: "Domains commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "identity",
		Title: "Identity commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "certs",
		Title: "Certificates commands",
	})

	cmd.AddCommand(NewCmdDomain(opts))
	cmd.AddCommand(NewCmdDSR(opts))
	cmd.AddCommand(NewCmdDNSSEC(opts))
	cmd.AddCommand(NewCmdCollaborator(opts))
	cmd.AddCommand(NewCmdEmailForward(opts))
	cmd.AddCommand(NewCmdWhoami(opts))
	cmd.AddCommand(NewCmdPush(opts))
	cmd.AddCommand(NewCmdCert(opts))
	cmd.AddCommand(NewCmdTLD(opts))
	cmd.AddCommand(NewCmdLetsEncrypt(opts))
	cmd.AddCommand(NewCmdTransfer(opts))
	cmd.AddCommand(NewCmdAccounts(opts))

	cmd.PersistentFlags().String("base-url", "", "Base URL")
	cmd.PersistentFlags().String("access-token", "", "Access token")
	cmd.PersistentFlags().Bool("sandbox", false, "Sandbox environment")
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file")
	cmd.PersistentFlags().StringVar(&profile, "profile", "default", "Profile")

	cmd.MarkFlagsMutuallyExclusive("base-url", "sandbox")

	viper.SetEnvPrefix("DNSIMPLE")
	_ = viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func lookupConfigFiles() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, "dnsimple"))
		viper.AddConfigPath("/etc/dnsimple")
		viper.SetConfigType("yaml")
		viper.SetConfigName(profile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Found error: ", err.Error())
		}
	}
}
