package main

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

type VocevoxCore struct {
}

func NewVoiceVox() (*VocevoxCore, error) {
	return NewVoiceVoxWith("./voicevox_core")
}

func NewVoiceVoxWith(dllRoot string) (*VocevoxCore, error) {
	var err error
	ret := &VocevoxCore{}

	return ret, err
}

func (v *VocevoxCore) Finalize() {
}

func (v *VocevoxCore) Initialize() error {
	return nil
}

func (v *VocevoxCore) Generate(speak_words, wav_file_path string) error {
	fmt.Println("音声生成中")
	var output_wav_ptr *uint8
	var output_binary_size uint

	if err := v.save(wav_file_path, output_wav_ptr, output_binary_size); err != nil {
		return err
	}

	return nil
}

func (v *VocevoxCore) save(dest_path string, output_wav_ptr *uint8, output_binary_size uint) error {
	//出力をポインタから取得
	output_wav := unsafe.Slice(output_wav_ptr, output_binary_size)

	//音声ファイルの保存
	f, err := os.Create(dest_path)
	if err != nil {
		return errors.New(fmt.Sprintf("wave Create() error [%v]", err.Error()))
	}
	defer f.Close()
	f.Write(output_wav)
	fmt.Println(dest_path + "に保存されました")
	return nil
}
