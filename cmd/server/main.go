package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// 音声ファイル名取得
func get_wave_file_name() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "speech.wav"), nil
}

func play_voicevox(v *VocevoxCore, text, dest_path string) error {
	if err := v.Generate(text, dest_path); err != nil {
		return err
	}
	if err := PlayWav(dest_path); err != nil {
		return err
	}
	return nil
}

func issue_ai_prompt(v *VocevoxCore, text, dest_path string) error {
	generated, err := IssueAiPrompt(text)
	if err != nil {
		play_voicevox(v, "失敗しました", dest_path)
		return err
	} else {
		return play_voicevox(v, generated, dest_path)
	}
}

func main() {

	v, err := NewVoiceVox()
	if err != nil {
		panic(err)
	}
	v.Initialize()

	fmt.Print("初期化完了")
	speak_words := []string{
		"今日の天気は晴れです。",
		"あかさたなはまやらわ",
	}

	dest_path, err := get_wave_file_name()
	if err != nil {
		panic(err)
	}

	for _, word := range speak_words {
		talkBack := fmt.Sprintf("「%v」と、聞いてみます", word)
		fmt.Println(talkBack)
		play_voicevox(v, talkBack, dest_path)
		issue_ai_prompt(v, word, dest_path)
	}
	v.Finalize()
}
