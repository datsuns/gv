package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

const COREDLL_DICT_NAME = "voicevox_core.dll"
const ONNXRUNTIME_DICT_NAME = "voicevox_onnxruntime.dll"
const OPENJTALK_DICT_NAME = "open_jtalk_dic_utf_8-1.11"

type VoicevoxAccelerationMode int32
type VoicevoxResultCode int32

const (
	VOICEVOX_ACCELERATION_MODE_AUTO VoicevoxAccelerationMode = iota
	VOICEVOX_ACCELERATION_MODE_CPU
	VOICEVOX_ACCELERATION_MODE_GPU
)

const (
	VOICEVOX_RESULT_OK VoicevoxResultCode = iota
	VOICEVOX_RESULT_NOT_LOADED_OPENJTALK_DICT_ERROR
	VOICEVOX_RESULT_LOAD_MODEL_ERROR
	VOICEVOX_RESULT_GET_SUPPORTED_DEVICES_ERROR
	VOICEVOX_RESULT_GPU_SUPPORT_ERROR
	VOICEVOX_RESULT_LOAD_METAS_ERROR
	VOICEVOX_RESULT_UNINITIALIZED_STATUS_ERROR
	VOICEVOX_RESULT_INVALID_SPEAKER_ID_ERROR
	VOICEVOX_RESULT_INVALID_MODEL_INDEX_ERROR
	VOICEVOX_RESULT_INFERENCE_ERROR
	VOICEVOX_RESULT_EXTRACT_FULL_CONTEXT_LABEL_ERROR
	VOICEVOX_RESULT_INVALID_UTF8_INPUT_ERROR
	VOICEVOX_RESULT_PARSE_KANA_ERROR
	VOICEVOX_RESULT_INVALID_AUDIO_QUERY_ERROR
)

type VoicevoxInitializeOptions struct {
	acceleration_mode VoicevoxAccelerationMode
	cpu_num_threads   uint16
}

type VoicevoxLoadOnnxruntimeOptions struct {
	filename uintptr
}

type VoicevoxTtsOptions struct {
	kana                         bool
	enable_interrogative_upspeak bool
}

type VocevoxCore struct {
	dll             *windows.LazyDLL
	func_initialize *windows.LazyProc
	func_tts        *windows.LazyProc
	func_wav_free   *windows.LazyProc
	func_finalize   *windows.LazyProc

	func_make_default_initialize_options *windows.LazyProc
	// .dllに含まれてない？
	//func_make_default_load_onnxruntime_options  *windows.LazyProc
	// func_get_onnxruntime_lib_versioned_filename *windows.LazyProc

	init_options VoicevoxInitializeOptions
	tts_options  VoicevoxTtsOptions
}

// デフォルトの初期化オプションを生成する
func voicevox_make_default_initialize_options() VoicevoxInitializeOptions {
	return VoicevoxInitializeOptions{
		acceleration_mode: VOICEVOX_ACCELERATION_MODE_AUTO,
		cpu_num_threads:   8,
	}
}

func voicevox_make_default_load_onnxruntime_options(prefix string) VoicevoxLoadOnnxruntimeOptions {
	filename := fmt.Sprintf("%v/onnxruntime/lib", prefix)
	filenamePtr, err := windows.BytePtrFromString(filename)
	if err != nil {
		log.Fatalf("failed to convert filename: %v", err)
	}
	return VoicevoxLoadOnnxruntimeOptions{}
}

// デフォルトのテキスト音声合成オプションを生成する
func voicevox_make_default_tts_options() VoicevoxTtsOptions {
	return VoicevoxTtsOptions{
		kana:                         true,
		enable_interrogative_upspeak: false,
	}
}

// dllファイルのパスを取得
func get_dll_dict() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(exePath)
}

// OpenJTalk辞書のパスを取得
func get_open_JTalk_dict() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, OPENJTALK_DICT_NAME)
}

func NewVoiceVox() (*VocevoxCore, error) {
	return NewVoiceVoxWith(get_dll_dict())
}

func NewVoiceVoxWith(dllRoot string) (*VocevoxCore, error) {

	ret := &VocevoxCore{}
	dll_path := filepath.Join(dllRoot, COREDLL_DICT_NAME)
	ret.load(dll_path)
	if err := ret.load_succeded(); err != nil {
		return nil, err
	}

	ret.init_options = voicevox_make_default_initialize_options()
	//ret.init_options.open_jtalk_dict_dir = []byte(get_open_JTalk_dict())

	ret.tts_options = voicevox_make_default_tts_options()
	return ret, nil
}

func (v *VocevoxCore) Finalize() {
	v.func_finalize.Call()
}

func (v *VocevoxCore) Initialize() error {
	result, _, _ := v.func_initialize.Call(uintptr(unsafe.Pointer(&v.init_options.acceleration_mode)))
	if VoicevoxResultCode(result) != VOICEVOX_RESULT_OK {
		return errors.New(fmt.Sprintf("voicevox_initialize() error [%v]", VoicevoxResultCode(result)))
	}
	return nil
}

func (v *VocevoxCore) Generate(speak_words, wav_file_path string) error {
	fmt.Println("音声生成中")
	speaker_id := 0
	speak_words_byte := []byte(speak_words)
	var output_wav_ptr *uint8
	var output_binary_size uint

	result, _, _ := v.func_tts.Call(uintptr(unsafe.Pointer(&speak_words_byte[0])), uintptr(speaker_id), uintptr(unsafe.Pointer(&v.tts_options.kana)), uintptr(unsafe.Pointer(&output_binary_size)), uintptr(unsafe.Pointer(&output_wav_ptr)))
	if VoicevoxResultCode(result) != VOICEVOX_RESULT_OK {
		return errors.New(fmt.Sprintf("voicevox_tts() error [%v]", VoicevoxResultCode(result)))
	}
	defer func() {
		fmt.Println("音声データの開放")
		v.func_wav_free.Call(uintptr(unsafe.Pointer(output_wav_ptr)))
	}()

	if err := v.save(wav_file_path, output_wav_ptr, output_binary_size); err != nil {
		return err
	}

	return nil
}

func (v *VocevoxCore) load(dll_path string) {
	d, _ := os.Getwd()
	fmt.Println("pwd: ", d, "dll:", dll_path)
	v.dll = windows.NewLazyDLL(dll_path)

	v.func_initialize = v.dll.NewProc("voicevox_initialize")
	v.func_tts = v.dll.NewProc("voicevox_tts")
	v.func_wav_free = v.dll.NewProc("voicevox_wav_free")
	v.func_finalize = v.dll.NewProc("voicevox_finalize")

}

func (v *VocevoxCore) load_succeded() error {
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
