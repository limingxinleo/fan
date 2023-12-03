package cmd

import (
	"fmt"
	"github.com/limingxinleo/fan/config"
	"github.com/logrusorgru/aurora"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"os"
	"regexp"
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

		var invoices = []string{}

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
						invoices = append(invoices, readInvoice(filename))
					}
				}
			}
		}

		result := decimal.New(0, 0)
		for _, invoice := range invoices {
			d2, _ := decimal.NewFromString(invoice)
			result = result.Add(d2)
		}

		fmt.Println(aurora.Blue("当前目录内发票金额总额为 " + result.String()))
	},
}

func readInvoice(path string) string {
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

	//fmt.Println(text)
	// 价税合计(大写)           贰拾贰圆整 (小写)¥22.00
	r, _ := regexp.Compile(`小写.*`)
	res := r.FindString(text)

	r2, _ := regexp.Compile(`\d+\.\d+`)
	result := r2.FindString(res)

	return result
}

func init() {
	rootCmd.AddCommand(pdfInvoiceCmd)
}
