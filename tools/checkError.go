package tools

import "fmt"

func CheckError(e error, texts ...string) {
	text := "Error happened"
	if len(texts) > 0 {
		text = texts[0]
	}
	if e != nil {
		fmt.Println(text, e)
	}
}
