# 🎵 Downloader de Clipes Musicais JW.ORG

Uma aplicação em Go para automatizar o download de clipes musicais do site oficial das Testemunhas de Jeová (jw.org), especificamente da seção de clipes musicais em português brasileiro.

## ✨ Funcionalidades

- 🎵 **Download automático** de todos os clipes musicais disponíveis
- 📁 **Organização inteligente** por ano de publicação
- 🚀 **Downloads concorrentes** para maior velocidade
- 🔄 **Retry automático** para falhas de download
- 📊 **Barra de progresso** em tempo real
- 🎯 **Download específico** por título
- 🔍 **Verificação** de novos clipes disponíveis
- ⚙️ **Configuração flexível** via arquivo YAML
- 📝 **Log detalhado** de todas as operações

## 🏗️ Arquitetura

O projeto utiliza **Arquitetura Hexagonal (Ports and Adapters)** para:
- ✅ Separação clara de responsabilidades
- ✅ Fácil testabilidade
- ✅ Flexibilidade para mudanças futuras
- ✅ Código limpo e manutenível

### 📂 Estrutura do Projeto

```
downloader-music/
├── cmd/
│   └── main.go                 # Ponto de entrada
├── internal/
│   ├── domain/                 # Regras de negócio
│   │   ├── clipe.go
│   │   └── repository.go
│   ├── infrastructure/         # Implementações concretas
│   │   ├── web/               # Web scraping
│   │   ├── storage/           # Sistema de arquivos
│   │   ├── download/          # Downloads HTTP
│   │   └── config/            # Configurações
│   ├── application/           # Casos de uso
│   │   └── download_service.go
│   └── interfaces/           # Adapters
│       └── cli/              # Interface CLI
├── pkg/                      # Utilitários
│   └── logger/
├── configs/
│   └── config.yaml           # Configurações
└── build/                    # Binários compilados
```

## 🚀 Instalação

### Pré-requisitos

- Go 1.19 ou superior
- Conexão com a internet

### Instalação via Git

```bash
# Clonar o repositório
git clone https://github.com/sant0x00/downloader-music.git
cd downloader-music

# Instalar dependências e compilar
make deps
make build

# Ou instalar diretamente no sistema
make install
```

### Instalação via Go

```bash
go install github.com/sant0x00/downloader-music/cmd@latest
```

## 📖 Como Usar

### 🎵 Download de Todos os Clipes

```bash
# Usando make
make run-download-all

# Ou diretamente
./build/downloader-music download all
```

### 🎯 Download de Clipe Específico

```bash
# Exemplo
./build/downloader-music download title "Vou até o fim"
```

### 🔍 Verificar Novos Clipes

```bash
# Usando make
make run-check

# Ou diretamente
./build/downloader-music check
```

### ⚙️ Configurar Diretório de Saída

```bash
./build/downloader-music config output-dir "~/Meus_Videos/ClipesJW"
```

### 📋 Ver Ajuda

```bash
./build/downloader-music --help
```

## ⚙️ Configuração

O arquivo `configs/config.yaml` permite personalizar:

```yaml
download:
  concurrent_workers: 8        # Número de downloads simultâneos
  retry_attempts: 3            # Tentativas em caso de falha
  timeout_seconds: 30          # Timeout por download
  output_directory: "~/Downloads/ClipesJW"  # Diretório de saída

scraping:
  base_url: "https://www.jw.org/pt/biblioteca/musica-canticos/clipes-musicais/"
  delay_between_requests: 1s   # Delay entre requisições
  user_agent: "ClipesJW-Downloader/1.0"

logging:
  level: "info"               # debug, info, warn, error
  output_file: "downloader.log"
```

## 📁 Estrutura de Saída

Os clipes são organizados automaticamente:

```
~/Downloads/ClipesJW/
├── 2025/
│   ├── Fazer_tua_vontade_e_o_meu_prazer.mp4
│   └── Vou_ate_o_fim.mp4
├── 2024/
│   ├── As_boas_novas_sobre_Jesus.mp4
│   └── ...
├── 2023/
│   └── ...
└── outros/
    ├── Cada_minuto.mp4
    └── E_tanto_amor.mp4
```

## 🔧 Desenvolvimento

### Comandos Make Disponíveis

```bash
make help            # Mostra todos os comandos
make deps            # Instala dependências
make build           # Compila a aplicação
make run             # Executa a aplicação
make test            # Executa testes
make clean           # Remove arquivos de build
make lint            # Formata e verifica código
make dev             # Preparação para desenvolvimento
make release         # Build de release
```

### Executar em Modo Debug

```bash
# Configurar log level para debug no config.yaml
logging:
  level: "debug"

# Executar com verbose
./build/downloader-music download all --verbose
```

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Executar testes específicos
go test ./internal/domain/...
```

## 📊 Performance

- **Downloads concorrentes**: 8 workers por padrão
- **Retry automático**: 3 tentativas com backoff progressivo
- **Rate limiting**: 1 segundo entre requisições de scraping
- **Timeout**: 30 segundos por download

## 🔒 Considerações Legais

⚠️ **IMPORTANTE**: Esta aplicação deve ser usada apenas para:
- ✅ Fins pessoais e educacionais
- ✅ Respeitando os termos de uso do site jw.org
- ✅ Não sobrecarregar os servidores (rate limiting implementado)

## 🐛 Solução de Problemas

### Erro de Conexão
```bash
# Verificar conectividade
curl -I https://www.jw.org/pt/biblioteca/musica-canticos/clipes-musicais/
```

### Erro de Permissão de Escrita
```bash
# Verificar permissões do diretório
ls -la ~/Downloads/
```

### Downloads Falhando
- Verificar espaço em disco
- Verificar configuração de timeout
- Aumentar retry_attempts na configuração

## 📝 Log

Os logs são salvos em `downloader.log` por padrão:

```bash
# Visualizar logs em tempo real
tail -f downloader.log

# Buscar erros
grep ERROR downloader.log
```

## 🔄 Atualizações

```bash
# Atualizar código
git pull origin main

# Recompilar
make clean
make build
```

## 📞 Suporte

Se encontrar problemas:

1. Verificar os logs: `cat downloader.log`
2. Testar conectividade com o site
3. Verificar configurações em `configs/config.yaml`
4. Criar uma issue no repositório

## 📄 Licença

Este projeto é fornecido "como está" para fins educacionais e pessoais. 

---

**Desenvolvido com ❤️ em Go** | **Arquitetura Hexagonal** | **Downloads Inteligentes**
