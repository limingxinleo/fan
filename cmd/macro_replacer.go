package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var macroReplacer = &cobra.Command{
	Use:   "ad:macro_replace {url} {macro}",
	Short: "覆盖广告链接的宏变量",
	Long:  "覆盖广告链接的宏变量",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(aurora.Yellow("请输入投放链接和宏变量提示"))
			os.Exit(1)
		}

		url := args[0]
		notice := args[1]

		parts := strings.Split(notice, ",")

		for _, part := range parts {
			url = strings.Replace(url, fmt.Sprintf("__%s__", part), "", -1)
		}
		fmt.Println()
		fmt.Println(url)
	},
}

func init() {
	rootCmd.AddCommand(macroReplacer)
}
