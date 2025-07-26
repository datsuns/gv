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
	"unsafe"
)

type VoicevoxStyleId C.VoicevoxStyleId
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
	init_option        C.struct_VoicevoxInitializeOptions
	onnxruntime_option C.struct_VoicevoxLoadOnnxruntimeOptions
	tts_option         C.struct_VoicevoxTtsOptions

	runtime     *C.struct_VoicevoxOnnxruntime
	dict        *C.OpenJtalkRc
	synthesizer *C.VoicevoxSynthesizer
}

func (v *VocevoxCoreApi) voicevox_make_default_initialize_options() C.struct_VoicevoxInitializeOptions {
	return C.voicevox_make_default_initialize_options()
}

func (v *VocevoxCoreApi) voicevox_make_default_load_onnxruntime_options() C.struct_VoicevoxLoadOnnxruntimeOptions {
	return C.voicevox_make_default_load_onnxruntime_options()
}

func (v *VocevoxCoreApi) voicevox_make_default_tts_options() C.struct_VoicevoxTtsOptions {
	return C.voicevox_make_default_tts_options()
}

func (v *VocevoxCoreApi) voicevox_error_result_to_message(code C.VoicevoxResultCode) string {
	return C.GoString(C.voicevox_error_result_to_message(code))
}

func (v *VocevoxCoreApi) voicevox_onnxruntime_load_once() (*C.struct_VoicevoxOnnxruntime, error) {
	var runtime *C.struct_VoicevoxOnnxruntime
	ret := C.voicevox_onnxruntime_load_once(v.onnxruntime_option, &runtime)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return nil, errors.New(
			fmt.Sprintf("voicevox_onnxruntime_load_once() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return runtime, nil
}

func (v *VocevoxCoreApi) voicevox_open_jtalk_rc_new(open_jtalk_dic_dir string) (*C.OpenJtalkRc, error) {
	var dict *C.OpenJtalkRc
	ret := C.voicevox_open_jtalk_rc_new(C.CString(open_jtalk_dic_dir), &dict)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return nil, errors.New(
			fmt.Sprintf("voicevox_open_jtalk_rc_new() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return dict, nil
}

func (v *VocevoxCoreApi) voicevox_synthesizer_new() (*C.VoicevoxSynthesizer, error) {
	var synthesizer *C.VoicevoxSynthesizer
	ret := C.voicevox_synthesizer_new(v.runtime, v.dict, v.init_option, &synthesizer)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return nil, errors.New(
			fmt.Sprintf("voicevox_synthesizer_new() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return synthesizer, nil
}

func (v *VocevoxCoreApi) voicevox_open_jtalk_rc_delete() {
	C.voicevox_open_jtalk_rc_delete(v.dict)
}

func (v *VocevoxCoreApi) voicevox_voice_model_file_open(path string) (*C.VoicevoxVoiceModelFile, error) {
	var model *C.struct_VoicevoxVoiceModelFile
	ret := C.voicevox_voice_model_file_open(C.CString(path), &model)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return nil, errors.New(
			fmt.Sprintf("voicevox_voice_model_file_open() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return model, nil
}

func (v *VocevoxCoreApi) voicevox_synthesizer_load_voice_model(model *C.VoicevoxVoiceModelFile) error {
	ret := C.voicevox_synthesizer_load_voice_model(v.synthesizer, model)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return errors.New(
			fmt.Sprintf("voicevox_synthesizer_load_voice_model() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return nil
}

func (v *VocevoxCoreApi) voicevox_voice_model_file_delete(model *C.VoicevoxVoiceModelFile) {
	C.voicevox_voice_model_file_delete(model)
}

func (v *VocevoxCoreApi) voicevox_synthesizer_delete() {
	C.voicevox_synthesizer_delete(v.synthesizer)
}

func (v *VocevoxCoreApi) voicevox_synthesizer_tts(text string, id VoicevoxStyleId) (int, []byte, error) {
	var size C.uintptr_t
	var buf *C.uint8_t
	v.tts_option = v.voicevox_make_default_tts_options()
	ret := C.voicevox_synthesizer_tts(v.synthesizer, C.CString(text), C.VoicevoxStyleId(id), v.tts_option, &size, &buf)
	if VoicevoxResultCode(ret) != VOICEVOX_RESULT_OK {
		return 0, nil, errors.New(
			fmt.Sprintf("voicevox_synthesizer_tts() error: %v", v.voicevox_error_result_to_message(ret)),
		)
	}
	return int(C.int(size)), C.GoBytes(unsafe.Pointer(buf), C.int(size)), nil
}

func (v *VocevoxCoreApi) load(lib_root string) error {
	var err error

	v.init_option = v.voicevox_make_default_initialize_options()
	v.init_option.acceleration_mode = C.VoicevoxAccelerationMode(VOICEVOX_ACCELERATION_MODE_CPU)
	v.onnxruntime_option = v.voicevox_make_default_load_onnxruntime_options()
	raw_onnxruntime_path := C.GoString(v.onnxruntime_option.filename)

	real_path := filepath.Join(lib_root, "onnxruntime", "lib", raw_onnxruntime_path)
	v.onnxruntime_option.filename = C.CString(real_path)
	v.runtime, err = v.voicevox_onnxruntime_load_once()
	if err != nil {
		return err
	}
	return nil
}

func (v *VocevoxCoreApi) load_dict(lib_root string) error {
	var err error
	dict_path := filepath.Join(lib_root, "dict", "open_jtalk_dic_utf_8-1.11")
	v.dict, err = v.voicevox_open_jtalk_rc_new(dict_path)
	if err != nil {
		return err
	}
	return nil
}

func (v *VocevoxCoreApi) load_model(lib_root string, style_id VoicevoxStyleId) error {
	var err error

	model_pattern := filepath.Join(lib_root, "models", "vvms", fmt.Sprintf("%v.vvm", style_id))
	model_path_list, err := filepath.Glob(model_pattern)
	if err != nil {
		return err
	}

	for _, path := range model_path_list {
		model, err := v.voicevox_voice_model_file_open(path)
		if err != nil {
			return err
		}

		err = v.voicevox_synthesizer_load_voice_model(model)
		if err != nil {
			return err
		}
		v.voicevox_voice_model_file_delete(model)
	}
	return nil
}

func (v *VocevoxCoreApi) Finalize() error {
	v.voicevox_synthesizer_delete()
	return nil
}

func (v *VocevoxCoreApi) Build(text string, id VoicevoxStyleId) error {
	_, _, err := v.voicevox_synthesizer_tts(text, id)
	if err != nil {
		return err
	}
	return nil
}

func NewVocevoxCoreApi(lib_root string, style_id VoicevoxStyleId) (*VocevoxCoreApi, error) {
	var err error

	v := &VocevoxCoreApi{}

	err = v.load(lib_root)
	if err != nil {
		return nil, err
	}

	err = v.load_dict(lib_root)
	if err != nil {
		return nil, err
	}

	v.synthesizer, err = v.voicevox_synthesizer_new()
	if err != nil {
		return nil, err
	}
	v.voicevox_open_jtalk_rc_delete()

	err = v.load_model(lib_root, style_id)
	if err != nil {
		return nil, err
	}

	return v, nil
}
