package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fan",
	Short: "一个简单的命令集合",
	Long:  "一个简单的命令集合",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringP("config", "c", u.HomeDir+"/.github/.fan.json", "配置文件")
}
