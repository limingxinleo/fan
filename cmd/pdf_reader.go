package cmd

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/logrusorgru/aurora"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

var pdfInvoiceCmd = &cobra.Command{
	Use:   "pdf:invoice {dir}",
	Short: "读取 PDF 发票信息",
	Long:  `读取某文件夹内 PDF 发票金额，并计算总额`,
	Run: func(cmd *cobra.Command, args []string) {
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
						invoices = append(invoices, getInvoice(readPdf(filename)))
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

func readPdf(path string) string {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var buf bytes.Buffer
	b, _ := r.GetPlainText()
	buf.ReadFrom(b)
	str := buf.String()
	if str == "" {
		if getInvoice(path) == "" {
			fmt.Println(path + " 无法正常读取金额，请直接命名为 (小写 xx.xx元)")
		}
		return path
	}

	return str
}

func getInvoice(text string) string {
	r, _ := regexp.Compile(`小写.*`)
	res := r.FindString(text)

	r2, _ := regexp.Compile(`\d+\.\d+`)
	result := r2.FindString(res)
	return result
}

func init() {
	rootCmd.AddCommand(pdfInvoiceCmd)
}
