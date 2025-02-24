# Nome do binário
BINARY_NAME=json2md

# Diretório de saída
OUTPUT_DIR=bin

# Comando padrão
all: linux macos windows

# Compilação para Linux
linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_amd64

# Compilação para macOS
macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME)_darwin_amd64

# Compilação para Windows
windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_amd64.exe

# Limpeza dos binários
clean:
	rm -rf $(OUTPUT_DIR)/*
