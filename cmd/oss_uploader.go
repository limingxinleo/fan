package cmd

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/limingxinleo/fan/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
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

		file := args[0]
		cf := config.GetConfig(cmd)
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if bucket != "" {
			cf.OssConfig.Bucket = bucket
		}

		if cf.OssConfig.BaseUri == "" {
			cf.OssConfig.BaseUri = "https://" + cf.OssConfig.Bucket + "." + cf.OssConfig.Endpoint
		}

		stat, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := oss.New(cf.OssConfig.Endpoint, cf.OssConfig.AccessKeyId, cf.OssConfig.AccessKeySecret)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		target, err := cmd.Flags().GetString("target")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if stat.IsDir() {
			if target == "" {
				fmt.Println(aurora.Yellow("上传目录必须通过 -t 指定上传目录."))
				os.Exit(1)
			}
			var files []string
			err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				files = append(files, path)
				return nil
			})

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, s := range files {
				stat, err = os.Stat(s)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if !stat.IsDir() {
					t := strings.TrimSuffix(target, "/") + strings.Replace(s, file, "", 1)
					uploadToOss(client, &cf.OssConfig, s, t)
				}
			}
		} else {
			uploadToOss(client, &cf.OssConfig, file, target)
		}
	},
}

func uploadToOss(client *oss.Client, cf *config.OssConfig, file string, target string) {
	extension := path.Ext(file)
	bucket, err := client.Bucket(cf.Bucket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if target == "" {
		name := uuid.New().String()
		target = path.Join(carbon.Now().Format("Y/m/d"), name) + extension
	}

	err = bucket.PutObjectFromFile(target, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	url := strings.TrimSuffix(cf.BaseUri, "/") + "/" + target
	fmt.Println(url)
}

func init() {
	rootCmd.AddCommand(ossUploaderCmd)

	ossUploaderCmd.Flags().StringP("target", "t", "", "上传路径")
	ossUploaderCmd.Flags().StringP("bucket", "b", "", "OSS Bucket")
}
