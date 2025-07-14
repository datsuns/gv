package main

import "testing"

func TestVoiceVox_init(t *testing.T) {
	_, err := NewVoiceVoxWith("../../_build")
	if err != nil {
		t.Errorf("NewVoiceVox error : %v", err.Error())
	}

	// ロード時間がとても長い
	//if err := v.Initialize(); err != nil {
	//	t.Errorf("Initialize error : %v", err.Error())
	//}
}
