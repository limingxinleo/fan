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
	"time"

	"github.com/google/uuid"

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

		if target == "" && !stat.IsDir() {
			target = generateQiniuTarget(file)
		}

		mac := credentials.NewCredentials(cf.QiniuConfig.AccessKey, cf.QiniuConfig.SecretKey)

		client := uploader.NewUploadManager(&uploader.UploadManagerOptions{
			Options: http_client.Options{
				Credentials: mac,
			},
		})

		if stat.IsDir() {
			prefix := strings.Trim(target, "/")
			if prefix == "" {
				fmt.Println("请使用 -t 输入上传的目标目录")
				os.Exit(1)
			}

			err := client.UploadDirectory(context.Background(), file, &uploader.DirectoryOptions{
				BucketName: cf.QiniuConfig.Bucket,
				UpdateObjectName: func(key string) string {
					return prefix + "/" + key
				},
				ObjectConcurrency: 16, // 对象上传并发度,
				BeforeObjectUpload: func(filePath string, info *uploader.ObjectOptions) {
					fmt.Println(strings.TrimSuffix(cf.QiniuConfig.BaseUri, "/") + "/" + *info.ObjectName)
				},
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
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

func generateQiniuTarget(file string) string {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	return fmt.Sprintf("%s/%s.%s", time.Now().Format("2006/01/02"), uuid.NewString(), ext)
}

func uploadToQiniu(client *uploader.UploadManager, cf *config.QiniuConfig, file string, target string) error {
	fmt.Println(cf, file, target)
	return client.UploadFile(context.Background(), file, &uploader.ObjectOptions{
		BucketName: cf.Bucket,
		ObjectName: &target,
		FileName:   file,
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
