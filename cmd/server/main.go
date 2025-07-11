package main

import (
	"fmt"
	"os"
)

func main() {

	v, err := NewVoiceVox()
	if err != nil {
		panic(err)
	}

	fmt.Print("生成する音声の文字列を入力>>")
	speak_words := ""
	if len(os.Args) == 1 {
		fmt.Scan(&speak_words)
	} else {
		speak_words = os.Args[1]
		fmt.Println(speak_words)
	}

	v.Generate(speak_words)
	v.Finalize()
}
