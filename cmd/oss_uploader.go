package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"

	"github.com/spf13/cobra"
)

// ossUploaderCmd represents the ossUploader command
var ossUploaderCmd = &cobra.Command{
	Use:   "oss:upload",
	Short: "Upload some objects to OSS",
	Long:  `Upload some objects to OSS`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(aurora.Yellow("请输入待上传文件或目录地址."))
			os.Exit(1)
		}

		path := args[0]
		stat, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if stat.IsDir() {
			fmt.Println("是目录")
		} else {
			fmt.Println("是文件")
		}
	},
}

func init() {
	rootCmd.AddCommand(ossUploaderCmd)

	ossUploaderCmd.Flags().StringP("target", "t", "", "上传路径")
}
