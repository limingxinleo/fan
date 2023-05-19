package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "s:md5",
	Short: "计算MD5",
	Long:  "计算MD5",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(aurora.Yellow("Please input the str which will be used to generate md5 string."))
			os.Exit(1)
		}

		hash := md5.New()
		hash.Write([]byte(args[0]))
		ret := hex.EncodeToString(hash.Sum(nil))

		fmt.Println(aurora.Green(ret))
	},
}

func init() {
	rootCmd.AddCommand(md5Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// md5Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// md5Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
