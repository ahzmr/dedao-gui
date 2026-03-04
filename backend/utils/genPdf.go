package utils

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type PdfOption struct {
	FileName  string
	CoverPath string
	PageSize  string
	Toc       bool
}

func (p *PdfOption) GenPdf(buf *bytes.Buffer) (err error) {
	wkhtmltopdf.SetPath(WkToPdfDir)
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		err = fmt.Errorf("初始化 PDF 生成器失败（请确认 wkhtmltopdf 已安装并在设置中配置路径）: %w", err)
		return
	}
	page := wkhtmltopdf.NewPageReader(buf)
	page.FooterFontSize.Set(10)
	page.FooterRight.Set("[page]")
	page.DisableSmartShrinking.Set(true)

	page.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(page)

	if p.CoverPath != "" {
		pdfg.Cover.EnableLocalFileAccess.Set(true)

		if runtime.GOOS == "windows" {
			pdfg.Cover.Input = p.CoverPath
		} else {
			pdfg.Cover.Input = "file://" + p.CoverPath
		}
	}

	pdfg.Dpi.Set(300)
	if p.Toc {
		pdfg.TOC.Include = true
		pdfg.TOC.TocHeaderText.Set("目 录")
		pdfg.TOC.HeaderFontSize.Set(18)

		pdfg.TOC.TocLevelIndentation.Set(15)
		pdfg.TOC.TocTextSizeShrink.Set(0.9)
		pdfg.TOC.DisableDottedLines.Set(false)
		pdfg.TOC.EnableTocBackLinks.Set(true)
	}

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.MarginTop.Set(15)
	pdfg.MarginBottom.Set(15)
	pdfg.MarginLeft.Set(15)
	pdfg.MarginRight.Set(15)
	err = pdfg.Create()
	if err != nil {
		fmt.Printf("pdfg create err: %#v\n", err)
		return
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile(p.FileName)
	if err != nil {
		return
	}
	if p.CoverPath != "" {
		err = os.Remove(p.CoverPath)
	}
	return
}
