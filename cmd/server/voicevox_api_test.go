package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVoiceVox_init(t *testing.T) {
	var err error
	wd, _ := os.Getwd()
	api, err := NewVocevoxCoreApi(filepath.Join(wd, "..", "..", "voicevox_core"), 3)
	if err != nil {
		t.Errorf("NewVocevoxCoreApi error : %v", err.Error())
	}

	err = api.Finalize()
	if err != nil {
		t.Errorf("Finalize error : %v", err.Error())
	}
}
