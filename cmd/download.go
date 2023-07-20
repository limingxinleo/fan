package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var downloadM3u8 = &cobra.Command{
	Use:   "d:m3u8",
	Short: "下载 m3u8 到本地",
	Long:  "下载 m3u8 到本地",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(aurora.Yellow("Please input the url about m3u8."))
			os.Exit(1)
		}

		name := uuid.New().String() + ".mp4"
		cmd2 := exec.Command(
			"ffmpeg",
			"-allowed_extensions",
			"ALL",
			"-protocol_whitelist",
			"'file,https,crypto,tcp,tls,http'",
			"-i",
			args[0],
			"-vcodec",
			"copy",
			"-acodec",
			"copy",
			name,
		)

		e := cmd2.Run()
		if e != nil {
			fmt.Println(e)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadM3u8)

	downloadM3u8.Flags().StringP("dir", "D", "", "下载位置")
}
