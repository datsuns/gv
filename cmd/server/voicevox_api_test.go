package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVoiceVox_init(t *testing.T) {
	wd, _ := os.Getwd()
	_, err := NewVocevoxCoreApi(filepath.Join(wd, "..", "..", "voicevox_core"))
	if err != nil {
		t.Errorf("NewVocevoxCoreApi error : %v", err.Error())
	}

	// ロード時間がとても長い
	//if err := v.Initialize(); err != nil {
	//	t.Errorf("Initialize error : %v", err.Error())
	//}
}
