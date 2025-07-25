//go:build windows
// +build windows

package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

const COREDLL_DICT_NAME = "voicevox_core.dll"
const ONNXRUNTIME_DICT_NAME = "voicevox_onnxruntime.dll"
const OPENJTALK_DICT_NAME = "open_jtalk_dic_utf_8-1.11"

type VoicevoxAccelerationMode int32
type VoicevoxResultCode int32
type VoicevoxOnnxruntime struct{}

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
	Filename []byte
}

type VoicevoxTtsOptions struct {
	kana                         bool
	enable_interrogative_upspeak bool
}

type VocevoxCoreApi struct {
	dll *windows.DLL

	func_voicevox_wav_free              *windows.Proc
	func_voicevox_onnxruntime_load_once *windows.Proc
	// 多分LD_LIBRARY_PATH上にvoicevox_onnxruntime.dllが無いとダメそう
	// func_voicevox_get_onnxruntime_lib_versioned_filename *windows.Proc

	init_options VoicevoxInitializeOptions
	onnx_options VoicevoxLoadOnnxruntimeOptions
	tts_options  VoicevoxTtsOptions
}

// デフォルトの初期化オプションを生成する
func voicevox_make_default_initialize_options() VoicevoxInitializeOptions {
	return VoicevoxInitializeOptions{
		acceleration_mode: VOICEVOX_ACCELERATION_MODE_AUTO,
		cpu_num_threads:   0,
	}
}

func voicevox_make_default_load_onnxruntime_options(lib_path string) (VoicevoxLoadOnnxruntimeOptions, error) {
	path := filepath.ToSlash(lib_path)
	// fmt.Println(path)
	// filenamePtr, err := windows.BytePtrFromString(path)
	// if err != nil {
	// 	return VoicevoxLoadOnnxruntimeOptions{}, errors.New(fmt.Sprintf("voicevox_make_default_load_onnxruntime_options error %v", err.Error()))
	// }
	return VoicevoxLoadOnnxruntimeOptions{
		Filename: []byte(path),
	}, nil
}

// デフォルトのテキスト音声合成オプションを生成する
func voicevox_make_default_tts_options() VoicevoxTtsOptions {
	return VoicevoxTtsOptions{
		kana:                         true,
		enable_interrogative_upspeak: false,
	}
}

func build_core_dll_path(root string) string {
	return filepath.Join(root, "c_api", "lib", COREDLL_DICT_NAME)
}

func build_onnxruntime_path(root string) string {
	filenameRaw := filepath.Join(root, "onnxruntime", "lib", ONNXRUNTIME_DICT_NAME)
	// filename := filepath.ToSlash(filenameRaw)
	filename := filenameRaw
	return filename
}

func build_dll_dict_path(root string) string {
	return filepath.Join(root, "dict", OPENJTALK_DICT_NAME)
}

func NewVocevoxCoreApi(lib_root string) (*VocevoxCoreApi, error) {
	var err error
	ret := &VocevoxCoreApi{}

	dll_path := build_core_dll_path(lib_root)
	if err := ret.load(dll_path); err != nil {
		return nil, err
	}

	ret.init_options = voicevox_make_default_initialize_options()
	onnx_path := build_onnxruntime_path(lib_root)
	ret.onnx_options, err = voicevox_make_default_load_onnxruntime_options(onnx_path)
	if err != nil {
		return nil, err
	}
	if _, err = ret.voicevox_onnxruntime_load_once(ret.onnx_options); err != nil {
		return nil, err
	}
	ret.tts_options = voicevox_make_default_tts_options()
	return ret, nil
}

func (a *VocevoxCoreApi) load(dll_path string) error {
	var err error

	a.dll, err = windows.LoadDLL(dll_path)
	if err != nil {
		return err
	}
	a.func_voicevox_wav_free, err = a.dll.FindProc("voicevox_wav_free")
	if err != nil {
		return err
	}

	a.func_voicevox_onnxruntime_load_once, err = a.dll.FindProc("voicevox_onnxruntime_load_once")
	if err != nil {
		return err
	}

	return nil
}

func (a *VocevoxCoreApi) voicevox_onnxruntime_load_once(opt VoicevoxLoadOnnxruntimeOptions) (*VoicevoxOnnxruntime, error) {
	var runtime *VoicevoxOnnxruntime

	ret, _, err := a.func_voicevox_onnxruntime_load_once.Call(
		uintptr(unsafe.Pointer(&opt.Filename)),
		uintptr(unsafe.Pointer(&runtime)),
	)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return nil, errors.New(fmt.Sprintf("voicevox_onnxruntime_load_once() error %v", err.Error()))
	}
	return runtime, nil
}
