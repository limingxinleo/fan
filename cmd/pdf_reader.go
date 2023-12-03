package cmd

import (
	"fmt"
	"github.com/limingxinleo/fan/config"
	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"os"
	"strings"
)

var pdfInvoiceCmd = &cobra.Command{
	Use:   "pdf:invoice {dir}",
	Short: "读取 PDF 发票信息",
	Long:  `读取某文件夹内 PDF 发票金额，并计算总额`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig(cmd)
		err := license.SetMeteredKey(conf.PdfConfig.MeteredKey)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		for i := 0; i < len(args); i++ {
			path := args[i]
			dir, err := os.ReadDir(path)
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}

			for _, fi := range dir {
				if !fi.IsDir() {
					filename := path + "/" + fi.Name()
					if strings.HasSuffix(filename, ".pdf") {
						fmt.Println(filename)
						readInvoice(filename)
						os.Exit(0)
					}
				}
			}
		}
	},
}

func readInvoice(path string) float32 {
	opt := model.NewReaderOpts()
	reader, f, err := model.NewPdfReaderFromFile(path, opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	page, err := reader.GetPage(1)
	//导出文本
	extract, err := extractor.New(page)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	text, err := extract.ExtractText()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(text)
	fmt.Println(path)

	return 1
}

func init() {
	rootCmd.AddCommand(pdfInvoiceCmd)
}
