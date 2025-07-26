//go:build linux
// +build linux

package main

/*
#cgo CFLAGS: -I${SRCDIR}/../../voicevox_core/c_api/include
#cgo LDFLAGS: -L${SRCDIR}/../../voicevox_core/c_api/lib -lvoicevox_core
#include "voicevox_core.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"path/filepath"
)

type VoicevoxAccelerationMode int32
type VoicevoxResultCode C.VoicevoxResultCode
type VoicevoxOnnxruntime struct{}

const (
	VOICEVOX_RESULT_OK                              VoicevoxResultCode = 0
	VOICEVOX_RESULT_NOT_LOADED_OPENJTALK_DICT_ERROR VoicevoxResultCode = 1
	VOICEVOX_RESULT_GET_SUPPORTED_DEVICES_ERROR     VoicevoxResultCode = 3
	VOICEVOX_RESULT_GPU_SUPPORT_ERROR               VoicevoxResultCode = 4
	VOICEVOX_RESULT_INIT_INFERENCE_RUNTIME_ERROR    VoicevoxResultCode = 29
	VOICEVOX_RESULT_STYLE_NOT_FOUND_ERROR           VoicevoxResultCode = 6
	VOICEVOX_RESULT_MODEL_NOT_FOUND_ERROR           VoicevoxResultCode = 7
	VOICEVOX_RESULT_RUN_MODEL_ERROR                 VoicevoxResultCode = 8
	VOICEVOX_RESULT_ANALYZE_TEXT_ERROR              VoicevoxResultCode = 11
	VOICEVOX_RESULT_INVALID_UTF8_INPUT_ERROR        VoicevoxResultCode = 12
	VOICEVOX_RESULT_PARSE_KANA_ERROR                VoicevoxResultCode = 13
	VOICEVOX_RESULT_INVALID_AUDIO_QUERY_ERROR       VoicevoxResultCode = 14
	VOICEVOX_RESULT_INVALID_ACCENT_PHRASE_ERROR     VoicevoxResultCode = 15
	VOICEVOX_RESULT_OPEN_ZIP_FILE_ERROR             VoicevoxResultCode = 16
	VOICEVOX_RESULT_READ_ZIP_ENTRY_ERROR            VoicevoxResultCode = 17
	VOICEVOX_RESULT_INVALID_MODEL_HEADER_ERROR      VoicevoxResultCode = 28
	VOICEVOX_RESULT_MODEL_ALREADY_LOADED_ERROR      VoicevoxResultCode = 18
	VOICEVOX_RESULT_STYLE_ALREADY_LOADED_ERROR      VoicevoxResultCode = 26
	VOICEVOX_RESULT_INVALID_MODEL_DATA_ERROR        VoicevoxResultCode = 27
	VOICEVOX_RESULT_LOAD_USER_DICT_ERROR            VoicevoxResultCode = 20
	VOICEVOX_RESULT_SAVE_USER_DICT_ERROR            VoicevoxResultCode = 21
	VOICEVOX_RESULT_USER_DICT_WORD_NOT_FOUND_ERROR  VoicevoxResultCode = 22
	VOICEVOX_RESULT_USE_USER_DICT_ERROR             VoicevoxResultCode = 23
	VOICEVOX_RESULT_INVALID_USER_DICT_WORD_ERROR    VoicevoxResultCode = 24
	VOICEVOX_RESULT_INVALID_UUID_ERROR              VoicevoxResultCode = 25
)

const (
	VOICEVOX_ACCELERATION_MODE_AUTO VoicevoxAccelerationMode = iota
	VOICEVOX_ACCELERATION_MODE_CPU
	VOICEVOX_ACCELERATION_MODE_GPU
)

type VoicevoxInitializeOptions struct {
	acceleration_mode VoicevoxAccelerationMode
	cpu_num_threads   uint16
}

type VoicevoxLoadOnnxruntimeOptions struct {
	filename string
}

type VocevoxCoreApi struct {
	init_option VoicevoxInitializeOptions

	raw_onnxruntime_option C.struct_VoicevoxLoadOnnxruntimeOptions
}

func (v *VocevoxCoreApi) voicevox_make_default_initialize_options() VoicevoxInitializeOptions {
	opt := C.voicevox_make_default_initialize_options()
	return VoicevoxInitializeOptions{
		acceleration_mode: VoicevoxAccelerationMode(opt.acceleration_mode),
		cpu_num_threads:   uint16(opt.cpu_num_threads),
	}
}

func (v *VocevoxCoreApi) voicevox_make_default_load_onnxruntime_options() C.struct_VoicevoxLoadOnnxruntimeOptions {
	return C.voicevox_make_default_load_onnxruntime_options()
}

func (v *VocevoxCoreApi) voicevox_error_result_to_message(code C.VoicevoxResultCode) string {
	return C.GoString(C.voicevox_error_result_to_message(code))
}

func (v *VocevoxCoreApi) voicevox_onnxruntime_load_once() (*C.struct_VoicevoxOnnxruntime, error) {
	var runtime *C.struct_VoicevoxOnnxruntime
	ret := C.voicevox_onnxruntime_load_once(v.raw_onnxruntime_option, &runtime)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return nil, errors.New(
			fmt.Sprintf("voicevox_onnxruntime_load_once() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return runtime, nil
}

func NewVocevoxCoreApi(lib_root string) (*VocevoxCoreApi, error) {
	ret := &VocevoxCoreApi{}
	ret.init_option = ret.voicevox_make_default_initialize_options()
	ret.raw_onnxruntime_option = ret.voicevox_make_default_load_onnxruntime_options()
	raw_onnxruntime_path := C.GoString(ret.raw_onnxruntime_option.filename)
	real_path := filepath.Join(lib_root, "onnxruntime", "lib", raw_onnxruntime_path)
	ret.raw_onnxruntime_option.filename = C.CString(real_path)
	_, err := ret.voicevox_onnxruntime_load_once()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func main_linux_tmp() {
	// デフォルトのオプションを取得
	opts := C.voicevox_make_default_load_onnxruntime_options()
	fmt.Println(C.GoString(opts.filename))
	fmt.Printf("%T\n", opts)

	// 結果格納用ポインタ
	var onnxruntime *C.VoicevoxOnnxruntime

	// ONNX Runtime をロード
	result := C.voicevox_onnxruntime_load_once(opts, &onnxruntime)
	if result != C.VOICEVOX_RESULT_OK {
		msg := C.GoString(C.voicevox_error_result_to_message(result))
		fmt.Println("エラー:", msg)
		return
	}

	// 成功時のメッセージ
	fmt.Println("ONNX Runtime の初期化に成功しました。")

	// 取得して確認（例として）
	getRuntime := C.voicevox_onnxruntime_get()
	if getRuntime == onnxruntime {
		fmt.Println("同じ ONNX Runtime インスタンスです。")
	} else {
		fmt.Println("異なる ONNX Runtime インスタンスです。")
	}
}
