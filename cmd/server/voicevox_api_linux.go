//go:build linux
// +build linux

package main

/*
#cgo LDFLAGS: -ldl
#cgo CFLAGS: -I${SRCDIR}/../../voicevox_core/c_api/include
#include <stdlib.h>
#include <dlfcn.h>
#include "voicevox_core.h"

typedef VoicevoxResultCode (*onnxruntime_load_once_func)(
  VoicevoxLoadOnnxruntimeOptions options,
  const VoicevoxOnnxruntime **out_onnxruntime
);
*/
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type onnxruntimeLoadOnceFunc func(
	options C.VoicevoxLoadOnnxruntimeOptions,
	out **C.VoicevoxOnnxruntime,
) C.VoicevoxResultCode

func main() {
	// .so のパス
	libPath := "./voicevox_core/c_api/lib/libvoicevox_core.so"

	// ライブラリ読み込み
	lib := C.dlopen(C.CString(libPath), C.RTLD_LAZY)
	if lib == nil {
		fmt.Println("❌ Failed to load .so")
		os.Exit(1)
	}
	defer C.dlclose(lib)

	// シンボル取得
	sym := C.dlsym(lib, C.CString("voicevox_onnxruntime_load_once"))
	if sym == nil {
		fmt.Println("❌ Function not found")
		os.Exit(1)
	}

	// 関数ポインタにキャスト
	loadOnceFunc := *(*onnxruntimeLoadOnceFunc)(unsafe.Pointer(&sym))

	// オプション構築
	onnxPath := filepath.ToSlash("./voicevox_core/onnxruntime/lib/libvoicevox_onnxruntime.so")
	cPath := C.CString(onnxPath)
	defer C.free(unsafe.Pointer(cPath))

	options := C.VoicevoxLoadOnnxruntimeOptions{
		filename: cPath,
	}

	// 結果出力先
	var runtimePtr *C.VoicevoxOnnxruntime

	result := loadOnceFunc(options, (**C.VoicevoxOnnxruntime)(unsafe.Pointer(&runtimePtr)))
	if result != 0 {
		fmt.Printf("❌ voicevox_onnxruntime_load_once failed: result = %d\n", result)
		os.Exit(1)
	}

	fmt.Println("✅ voicevox_onnxruntime_load_once success!")
	fmt.Printf("Returned runtime ptr = %p\n", runtimePtr)
}
