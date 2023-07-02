package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/limingxinleo/fan/config"
	"github.com/spf13/cobra"
	"os"
)

// configCmd represents the entity command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show Config",
	Long:  `展示配置信息`,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.GetConfig(cmd)

		res, err := json.Marshal(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(res))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
