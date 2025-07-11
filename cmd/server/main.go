package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// 音声ファイル名取得
func GetWaveFileName() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "speech.wav")
}

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

	dest_path := GetWaveFileName()
	v.Generate(speak_words, dest_path)
	v.Finalize()
}
