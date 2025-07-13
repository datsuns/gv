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

	fmt.Print("初期化完了")
	speak_words := []string{
		"今日の天気は晴れです。",
		"あかさたなはまやらわ",
	}

	for _, word := range speak_words {
		dest_path := GetWaveFileName()
		v.Generate(word, dest_path)
		PlayWav(dest_path)
	}
	v.Finalize()
}
