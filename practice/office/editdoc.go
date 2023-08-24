package main

import (
	"fmt"

	"github.com/aymerick/raymond"
	"github.com/unidoc/unioffice/document"
)

func main() {
	// doc := document.New()
	// defer doc.Close()

	// para := doc.AddParagraph()
	// para.AddRun().AddText("{{name}}")

	// para = doc.AddParagraph()
	// para.AddRun().AddText("{{name}}")

	// doc.SaveToFile("d.docx")

	doc, err := document.Open("d.docx")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	ps := make([]document.Paragraph, 0)
	ps = append(ps, doc.Paragraphs()...)

	// // This sample document uses structured document tags, which are not common
	// // except for in document templates.  Normally you can just iterate over the
	// // document's paragraphs.
	for _, sdt := range doc.StructuredDocumentTags() {
		ps = append(ps, sdt.Paragraphs()...)
	}

	data := map[string]interface{}{"name": "wxy"}
	for _, p := range ps {
		for _, r := range p.Runs() {
			fmt.Println("XXX:", r.Text())
			txt, _ := raymond.Render(r.Text(), data)
			r.ClearContent()
			r.AddText(txt)
		}
	}

	doc.SaveToFile("f.docx")

	// for _, p := range ps {
	// 	for _, r := range p.Runs() {
	// 		switch r.Text() {
	// 		case "FIRST NAME":
	// 			// ClearContent clears both text and line breaks within a run,
	// 			// so we need to add the line break back
	// 			r.ClearContent()
	// 			r.AddText("John ")
	// 			r.AddBreak() // 换行符

	// 			// para := doc.InsertParagraphBefore(p)
	// 			// para.AddRun().AddText("Mr.")
	// 			// para.SetStyle("Name") // Name is a default style in this template file

	// 			// para = doc.InsertParagraphAfter(p)
	// 			// para.AddRun().AddText("III")
	// 			// para.SetStyle("Name")

	// 		case "LAST NAME":
	// 			r.ClearContent()
	// 			r.AddText("Smith")
	// 		case "Address | Phone | Email":
	// 			r.ClearContent()
	// 			r.AddText("111 Rustic Rd | 123-456-7890 | jsmith@smith.com")
	// 		case "Date":
	// 			r.ClearContent()
	// 			r.AddText(time.Now().Format("Jan 2, 2006"))
	// 		case "Recipient Name":
	// 			r.ClearContent()
	// 			r.AddText("Mrs. Smith")
	// 			r.AddBreak()
	// 		case "Title":
	// 			// we remove the title content entirely
	// 			p.RemoveRun(r)
	// 		case "Company":
	// 			r.ClearContent()
	// 			r.AddText("Smith Enterprises")
	// 			r.AddBreak()
	// 		case "Address":
	// 			r.ClearContent()
	// 			r.AddText("112 Rustic Rd")
	// 			r.AddBreak()
	// 		case "City, ST ZIP Code":
	// 			r.ClearContent()
	// 			r.AddText("San Francisco, CA 94016")
	// 			r.AddBreak()
	// 		case "Dear Recipient:":
	// 			r.ClearContent()
	// 			r.AddText("Dear Mrs. Smith:")
	// 			r.AddBreak()
	// 		case "Your Name":
	// 			r.ClearContent()
	// 			r.AddText("John Smith")
	// 			r.AddBreak()

	// 			run := p.InsertRunBefore(r)
	// 			run.AddText("---Before----")
	// 			run.AddBreak()
	// 			run = p.InsertRunAfter(r)
	// 			run.AddText("---After----")

	// 		default:
	// 			fmt.Println("not modifying", r.Text())
	// 		}
	// 	}
	// }

	// doc.SaveToFile("edit-document.docx")
}
