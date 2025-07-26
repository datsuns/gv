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
	"fmt"
)

type VoicevoxAccelerationMode int32

const (
	VOICEVOX_ACCELERATION_MODE_AUTO VoicevoxAccelerationMode = iota
	VOICEVOX_ACCELERATION_MODE_CPU
	VOICEVOX_ACCELERATION_MODE_GPU
)

type VoicevoxInitializeOptions struct {
	acceleration_mode VoicevoxAccelerationMode
	cpu_num_threads   uint16
}

type VocevoxCoreApi struct {
}

func (v *VocevoxCoreApi) voicevox_make_default_initialize_options() VoicevoxInitializeOptions {
	opt := C.voicevox_make_default_initialize_options()
	return VoicevoxInitializeOptions{
		acceleration_mode: VoicevoxAccelerationMode(opt.acceleration_mode),
		cpu_num_threads:   uint16(opt.cpu_num_threads),
	}
}

func NewVocevoxCoreApi(lib_root string) (*VocevoxCoreApi, error) {
	ret := &VocevoxCoreApi{}
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
