# 🎵 Exemplos de Uso - Downloader de Clipes Musicais JW.ORG

## 🚀 Início Rápido

```bash
# 1. Compilar a aplicação
make build

# 2. Verificar clipes disponíveis
./build/downloader-music check

# 3. Baixar todos os clipes
./build/downloader-music download all
```

## 📋 Comandos Principais

### 🔍 Verificar Novos Clipes
```bash
# Ver todos os clipes disponíveis
./build/downloader-music check
```

### 🎵 Download de Todos os Clipes
```bash
# Baixar todos os clipes disponíveis
./build/downloader-music download all

# Com logs verbosos
./build/downloader-music download all --verbose
```

### 🎯 Download de Clipe Específico
```bash
# Baixar um clipe específico por título
./build/downloader-music download title "Vou até o fim"

# Outros exemplos
./build/downloader-music download title "Fazer tua vontade é o meu prazer (cântico do congresso de 2025)"
./build/downloader-music download title "As boas novas sobre Jesus (cântico do congresso de 2024)"
```

### ⚙️ Configurações
```bash
# Configurar diretório de saída
./build/downloader-music config output-dir "~/Meus_Videos/ClipesJW"

# Ver ajuda de qualquer comando
./build/downloader-music --help
./build/downloader-music download --help
./build/downloader-music config --help
```

## 🛠️ Comandos Make

```bash
# Ver todos os comandos disponíveis
make help

# Compilar
make build

# Executar verificação
make run-check

# Executar download completo
make run-download-all

# Limpar arquivos de build
make clean

# Preparar para desenvolvimento
make dev
```

## 📁 Estrutura de Saída

Após o download, os clipes serão organizados em:

```
~/Downloads/ClipesJW/
├── 2025/
│   ├── Fazer_tua_vontade_e_o_meu_prazer_cantico_do_congresso_de_2025.mp4
│   └── Vou_ate_o_fim.mp4
├── 2024/
│   ├── As_boas_novas_sobre_Jesus_cantico_do_congresso_de_2024.mp4
│   └── ...
├── 2023/
│   ├── Nao_vai_se_atrasar_cantico_do_congresso_de_2023.mp4
│   └── ...
├── 2022/
│   ├── Paz_enfim_cantico_do_congresso_de_2022.mp4
│   └── ...
└── outros/
    ├── Cada_minuto.mp4
    ├── E_tanto_amor.mp4
    └── ...
```

## 🔧 Personalização

### Modificar Configurações

Edite o arquivo `configs/config.yaml`:

```yaml
download:
  concurrent_workers: 8        # Número de downloads simultâneos
  retry_attempts: 3            # Tentativas em caso de falha
  timeout_seconds: 30          # Timeout por download
  output_directory: "~/Downloads/ClipesJW"

scraping:
  base_url: "https://www.jw.org/pt/biblioteca/musica-canticos/clipes-musicais/"
  delay_between_requests: 1s   # Delay entre requisições
  user_agent: "ClipesJW-Downloader/1.0"

logging:
  level: "info"               # debug, info, warn, error
  output_file: "downloader.log"
```

### Para Downloads Mais Rápidos
```yaml
download:
  concurrent_workers: 12       # Aumentar workers
  timeout_seconds: 60          # Aumentar timeout
```

### Para Modo Debug
```yaml
logging:
  level: "debug"              # Logs mais detalhados
```

## 📊 Monitoramento

### Ver Logs em Tempo Real
```bash
tail -f downloader.log
```

### Buscar Erros nos Logs
```bash
grep ERROR downloader.log
```

### Ver Progresso
```bash
# Os downloads mostram progresso visual automaticamente
./build/downloader-music download all
```

## 🧪 Testes

### Testar um Clipe Específico
```bash
# Teste com um clipe simples primeiro
./build/downloader-music download title "A melhor vida"
```

### Verificar se Tudo Está Funcionando
```bash
# 1. Verificar conectividade
./build/downloader-music check

# 2. Testar configuração
./build/downloader-music config output-dir "/tmp/teste"

# 3. Voltar configuração
./build/downloader-music config output-dir "~/Downloads/ClipesJW"
```

## 🔄 Fluxo Recomendado

1. **Primeira execução**:
   ```bash
   make build
   ./build/downloader-music check
   ```

2. **Configurar diretório** (opcional):
   ```bash
   ./build/downloader-music config output-dir "~/MeusClipes"
   ```

3. **Download inicial**:
   ```bash
   ./build/downloader-music download all
   ```

4. **Verificações futuras**:
   ```bash
   ./build/downloader-music check
   ```

## 🆘 Solução de Problemas

### Erro de Conectividade
```bash
# Testar conectividade manual
curl -I https://www.jw.org/pt/biblioteca/musica-canticos/clipes-musicais/
```

### Downloads Falhando
```bash
# Verificar logs
tail -20 downloader.log

# Tentar com menos workers
# Editar config.yaml: concurrent_workers: 3
```

### Clipe Não Encontrado
```bash
# Verificar se o título está correto
./build/downloader-music check | grep "nome_do_clipe"
```

---

**💡 Dica**: Comece sempre com `./build/downloader-music check` para ver quantos clipes estão disponíveis!
