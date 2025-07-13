DEST_DIR := $(PWD)/_build
WORK_DIR := $(PWD)/_wk

VOICEVOX_CORE_URL        := https://github.com/VOICEVOX/voicevox_core/releases/download
VOICEVOX_CORE_VERSION    := 0.15.9
VOICEVOX_CORE_DOWNLOADER := download-windows-x64.exe
VOICEVOX_CORE_DL_OPT     := --device directml --cpu-arch x64 --os windows


default: build

build:
	go mod tidy
	go build -o $(DEST_DIR)/server ./cmd/server

test:

setup:
	go get -u golang.org/x/sys/windows
	go get -u github.com/gopxl/beep/v2
	go get -u google.golang.org/genai

s: server

server:
	$(DEST_DIR)/server

c: client

client:

lib: dl_lib copy_lib

dl_lib:
	mkdir -p $(WORK_DIR)
	cd $(WORK_DIR) && wget --no-check-certificate $(VOICEVOX_CORE_URL)/$(VOICEVOX_CORE_VERSION)/$(VOICEVOX_CORE_DOWNLOADER) \
		&& $(WORK_DIR)/$(VOICEVOX_CORE_DOWNLOADER) $(VOICEVOX_CORE_DL_OPT)

copy_lib:
	cp -r $(WORK_DIR)/voicevox_core/model                            $(DEST_DIR)/
	cp -r $(WORK_DIR)/voicevox_core/open_jtalk_dic_utf_8-1.11        $(DEST_DIR)/
	cp -r $(WORK_DIR)/voicevox_core/voicevox_core.dll                $(DEST_DIR)/
	cp -r $(WORK_DIR)/voicevox_core/onnxruntime.dll                  $(DEST_DIR)/

.PHONY: default build test setup s server c client lib dl_lib copy_lib
