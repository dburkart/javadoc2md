package main

import (
	"fmt"
	"strings"
)

type Visitor interface {
	visit(*Document) (bool, string)
}

type MarkdownVisitor struct {

}

func (v *MarkdownVisitor) visit(doc *Document) (err bool, description string) {
	err = false
	description = ""
	needs_newline := false

	for i, v := range doc.Blocks {
		heading := "## "
		if i == 0 {
			heading = "# "
		}

		fmt.Println(heading, v.Name)
		fmt.Println()

		fmt.Println(strings.TrimSpace(v.Description))
		fmt.Println()

		if len(v.Params) > 0 {
			fmt.Println("* **Parameters:**")
			needs_newline = true
		}

		for key, value := range v.Params {
			fmt.Print("\t* `", key, "` - ", value)
			fmt.Println()
		}

		if ret, found := v.Tags["@return"]; found {
			fmt.Println("* **Returns:**", ret)
			needs_newline = true
		}

		if needs_newline {
			fmt.Println()
		}
		needs_newline = false
	}
	return
}
