package main

import (
	"os"
	"path/filepath"
	"testing"
)

func util_build_lib_root_path() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "..", "..", "voicevox_core")
}

func TestVoiceVox_init(t *testing.T) {
	var err error
	api, err := NewVocevoxCoreApi(util_build_lib_root_path(), 3)
	if err != nil {
		t.Errorf("NewVocevoxCoreApi error : %v", err.Error())
	}

	err = api.Finalize()
	if err != nil {
		t.Errorf("Finalize error : %v", err.Error())
	}
}

func TestVoiceVox_build(t *testing.T) {
}
