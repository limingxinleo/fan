package cmd

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
)

var downloadM3u8 = &cobra.Command{
	Use:   "d:m3u8",
	Short: "下载 m3u8 到本地",
	Long:  "下载 m3u8 到本地",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(aurora.Yellow("Please input the url or path about m3u8."))
			os.Exit(1)
		}

		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = "./output"
		}

		resp, err := http.Get(args[0])
		if err != nil {
			fmt.Errorf("文件下载失败 %s", err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		bf := bufio.NewReader(resp.Body)
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Errorf("文件下载失败 %s", err)
			os.Exit(1)
		}

		f, err := os.OpenFile(dir+"/index.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		w := bufio.NewWriter(f)

		for {
			line, _, err := bf.ReadLine()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Printf("文件读取失败 %s", err)
				os.Exit(1)
			}

			url := string(line)

			if strings.HasPrefix(url, "http") {
				fmt.Println(url)
				arr := strings.Split(url, "/")
				file := arr[len(arr)-1]

				download(url, dir+"/"+file)

				w.WriteString("file '" + file + "'\n")
			}
		}

		// ffmpeg -f concat -safe 0 -i index.txt -c copy output.mkv
		w.Flush()
	},
}

func download(url string, path string) int64 {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer r.Body.Close()

	written, err := io.Copy(f, r.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(written)
	return written
}

func init() {
	rootCmd.AddCommand(downloadM3u8)

	downloadM3u8.Flags().StringP("dir", "D", "", "下载位置")
}
