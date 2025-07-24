.PHONY: help build run clean test install deps

BINARY_NAME=downloader-music
BUILD_DIR=build
MAIN_PATH=./cmd/main.go

help:
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps:
	@echo "📦 Instalando dependências..."
	go mod tidy
	go mod download

build: deps
	@echo "🔨 Compilando aplicação..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Aplicação compilada em $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "📦 Instalando aplicação..."
	go install $(MAIN_PATH)
	@echo "✅ Aplicação instalada. Execute com: $(BINARY_NAME)"

run: build
	@echo "🚀 Executando aplicação..."
	./$(BUILD_DIR)/$(BINARY_NAME)

run-download-all: build
	@echo "🎵 Baixando todos os clipes..."
	./$(BUILD_DIR)/$(BINARY_NAME) download all

run-check: build
	@echo "🔍 Verificando novos clipes..."
	./$(BUILD_DIR)/$(BINARY_NAME) check

clean:
	@echo "🧹 Limpando arquivos de build..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f downloader.log

test:
	@echo "🧪 Executando testes..."
	go test -v ./...

fmt:
	@echo "🎨 Formatando código..."
	go fmt ./...

vet:
	@echo "🔍 Verificando código..."
	go vet ./...

lint: fmt vet

dev: deps lint build

release: clean deps lint test build
	@echo "🎉 Build de release concluído!"

config-output: build
	@read -p "Digite o novo diretório de saída: " dir; \
	./$(BUILD_DIR)/$(BINARY_NAME) config output-dir "$$dir"

download-specific: build
	@read -p "Digite o título do clipe: " titulo; \
	./$(BUILD_DIR)/$(BINARY_NAME) download title "$$titulo"
