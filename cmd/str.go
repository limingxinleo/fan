package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var strCmd = &cobra.Command{
	Use:   "s:gen {length}",
	Short: "生成随机字符串",
	Long:  "生成随机字符串",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(aurora.Yellow("Please input the length of str."))
			os.Exit(1)
		}

		num, err := strconv.ParseInt(args[0], 10, 10)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(aurora.Green(randString(int(num))))
	},
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func init() {
	rootCmd.AddCommand(strCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// strCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// strCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
