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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/limingxinleo/fan/config"
	"github.com/logrusorgru/aurora"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
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

		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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

			var upToken uptoken.Provider
			if overwrite {
				upToken = newOverwriteUpToken(mac, cf.QiniuConfig.Bucket, "")
			}

			err := uploadDirectoryToQiniu(client, &cf.QiniuConfig, file, prefix, upToken)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			var upToken uptoken.Provider
			if overwrite {
				upToken = newOverwriteUpToken(mac, cf.QiniuConfig.Bucket, target)
			}

			err := uploadToQiniu(client, &cf.QiniuConfig, file, target, upToken)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(strings.TrimSuffix(cf.QiniuConfig.BaseUri, "/") + "/" + target)
		}
	},
}

func generateQiniuTarget(file string) string {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	return fmt.Sprintf("%s/%s.%s", time.Now().Format("2006/01/02"), uuid.NewString(), ext)
}

// newOverwriteUpToken 构造允许覆盖同名对象的上传凭证(insertOnly=0)
// key 非空时 scope 为 "bucket:key"(官方推荐的覆盖姿势),为空时 scope 为 "bucket"(用于目录上传)
func newOverwriteUpToken(mac *credentials.Credentials, bucket string, key string) uptoken.Provider {
	putPolicy, err := uptoken.NewPutPolicyWithKey(bucket, key, time.Now().Add(time.Hour))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return uptoken.NewSigner(putPolicy.SetInsertOnly(0), mac)
}

// qiniuDirectoryConcurrency 目录上传时的对象并发度
const qiniuDirectoryConcurrency = 16

// uploadDirectoryToQiniu 并发上传目录下的所有文件
// 单个对象上传失败不会中断其他对象的上传,所有对象处理完毕后统一返回失败信息
func uploadDirectoryToQiniu(client *uploader.UploadManager, cf *config.QiniuConfig, dir string, prefix string, upToken uptoken.Provider) error {
	var (
		mutex  sync.Mutex
		failed []string
		wg     sync.WaitGroup
	)

	semaphore := make(chan struct{}, qiniuDirectoryConcurrency)
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		objectName := prefix + "/" + filepath.ToSlash(rel)
		fmt.Println(strings.TrimSuffix(cf.BaseUri, "/") + "/" + objectName)

		wg.Add(1)
		semaphore <- struct{}{}
		go func(path string, objectName string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			err := client.UploadFile(context.Background(), path, &uploader.ObjectOptions{
				BucketName: cf.Bucket,
				ObjectName: &objectName,
				FileName:   filepath.Base(path),
				UpToken:    upToken, // 为 nil 时走 SDK 默认凭证,行为与现状完全一致
			}, nil)
			if err != nil {
				mutex.Lock()
				failed = append(failed, fmt.Sprintf("%s: %v", objectName, err))
				mutex.Unlock()
			}
		}(path, objectName)

		return nil
	})

	wg.Wait()

	var messages []string
	if walkErr != nil {
		messages = append(messages, walkErr.Error())
	}

	if len(failed) > 0 {
		sort.Strings(failed)
		messages = append(messages, fmt.Sprintf("%d 个对象上传失败:\n%s", len(failed), strings.Join(failed, "\n")))
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
}

func uploadToQiniu(client *uploader.UploadManager, cf *config.QiniuConfig, file string, target string, upToken uptoken.Provider) error {
	fmt.Println(cf, file, target)
	return client.UploadFile(context.Background(), file, &uploader.ObjectOptions{
		BucketName: cf.Bucket,
		ObjectName: &target,
		FileName:   file,
		UpToken:    upToken, // 为 nil 时走 SDK 默认凭证,行为与现状完全一致
		CustomVars: map[string]string{
			"name": "github logo",
		},
	}, nil)
}

func init() {
	rootCmd.AddCommand(qiniuUploaderCmd)

	qiniuUploaderCmd.Flags().StringP("target", "t", "", "上传路径")
	qiniuUploaderCmd.Flags().StringP("bucket", "b", "", "Qiniu Bucket")
	qiniuUploaderCmd.Flags().BoolP("overwrite", "f", false, "允许覆盖空间内已存在的同名对象")
}
