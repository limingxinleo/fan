/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/limingxinleo/fan/config"
	"github.com/logrusorgru/aurora"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/spf13/cobra"
)

// qiniuUploaderCmd represents the qiniuUploader command
var qiniuUploaderCmd = &cobra.Command{
	Use:   "qiniu:upload",
	Short: "上传文件到七牛",
	Long:  `上传文件到七牛`,
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
			cf.QiniuConfig.Bucket = bucket
		}

		if cf.QiniuConfig.BaseUri == "" {
			fmt.Println(errors.New("请配置 qiniu_config.base_uri"))
			os.Exit(1)
		}

		stat, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		target, err := cmd.Flags().GetString("target")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		mac := credentials.NewCredentials(cf.QiniuConfig.AccessKey, cf.QiniuConfig.SecretKey)

		client := uploader.NewUploadManager(&uploader.UploadManagerOptions{
			Options: http_client.Options{
				Credentials: mac,
			},
		})

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
					err := uploadToQiniu(client, &cf.QiniuConfig, s, t)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}
		} else {
			err := uploadToQiniu(client, &cf.QiniuConfig, file, target)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func uploadToQiniu(client *uploader.UploadManager, cf *config.QiniuConfig, file string, target string) error {
	return client.UploadFile(context.Background(), file, &uploader.ObjectOptions{
		BucketName: cf.Bucket,
		ObjectName: &target,
		CustomVars: map[string]string{
			"name": "github logo",
		},
	}, nil)
}

func init() {
	rootCmd.AddCommand(qiniuUploaderCmd)

	qiniuUploaderCmd.Flags().StringP("target", "t", "", "上传路径")
	qiniuUploaderCmd.Flags().StringP("bucket", "b", "", "Qiniu Bucket")
}
