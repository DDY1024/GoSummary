package main

import (
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
)

func imageTest() {
	doc := document.New()
	defer doc.Close()

	img1, err := common.ImageFromFile("gophercolor.png")
	if err != nil {
		panic(err)
	}

	img1Ref, err := doc.AddImage(img1)
	if err != nil {
		panic(err)
	}

	para := doc.AddParagraph()
	para.AddRun().AddDrawingAnchored(img1Ref) // floating drawing
	para.AddRun().AddDrawingInline(img1Ref)   // inline drawing
	// anchored.SetName("Gopher")
	// anchored.SetSize(2*measurement.Inch, 2*measurement.Inch)
	// anchored.SetOrigin(wml.WdST_RelFromHPage, wml.WdST_RelFromVTopMargin)
	// anchored.SetHAlignment(wml.WdST_AlignHCenter)
	// anchored.SetYOffset(3 * measurement.Inch)
	// anchored.SetTextWrapSquare(wml.WdST_WrapTextBothSides)

	doc.SaveToFile("image.docx")
}

/*
Floating Drawing:
在这种插入方式中，图形被视为一个"浮动"的对象，它可以自由地放置在文档中的某个位置，并且可以在文本的周围进行调整。这种方式适用于需要在文本内容之间插入大型图形，或者需要更多自由布局的情况。

Inline Drawing:
在这种插入方式中，图形被嵌入到文本中的特定位置，就像一个字符一样。这意味着图形会随着文本的流动而移动，而不会独立于文本位置。这种方式适用于需要图形与文本内容紧密结合的情况，例如插入图标、小的装饰性图像等。
*/
