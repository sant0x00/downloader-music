package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sant0x00/downloader-music/internal/application"
	"github.com/sant0x00/downloader-music/internal/domain"
	"github.com/sant0x00/downloader-music/internal/infrastructure/config"
	"github.com/sant0x00/downloader-music/internal/infrastructure/download"
	"github.com/sant0x00/downloader-music/internal/infrastructure/storage"
	"github.com/sant0x00/downloader-music/internal/infrastructure/web"
	"github.com/sant0x00/downloader-music/pkg/logger"
	"github.com/spf13/cobra"
)

func showBanner() {
	banner := `
 ██████╗ ██╗     ██╗██████╗ ███████╗███████╗        ██╗██╗    ██╗
██╔════╝ ██║     ██║██╔══██╗██╔════╝██╔════╝        ██║██║    ██║
██║      ██║     ██║██████╔╝█████╗  ███████╗        ██║██║ █╗ ██║
██║      ██║     ██║██╔═══╝ ██╔══╝  ╚════██║   ██   ██║██║███╗██║
╚██████╗ ███████╗██║██║     ███████╗███████║   ╚█████╔╝╚███╔███╔╝
 ╚═════╝ ╚══════╝╚═╝╚═╝     ╚══════╝╚══════╝    ╚════╝  ╚══╝╚══╝ 

🎵 Downloader de Clipes Musicais JW.org 🎵
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`
	fmt.Print(banner)
}

func showSmallBanner() {
	fmt.Println("🎵 Clipes JW.org Downloader 🎵")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

type CLI struct {
	config          *config.Config
	logger          domain.Logger
	downloadService *application.DownloadService
}

func NewCLI() (*CLI, error) {
	configPath := filepath.Join("configs", "config.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	log, err := logger.NewSimpleLogger(cfg.Logging.Level, cfg.Logging.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar logger: %w", err)
	}

	repository := storage.NewFileSystemRepository(cfg.Download.OutputDirectory, log)
	scraper := web.NewJWScraper(cfg.Scraping.UserAgent, cfg.Scraping.DelayBetweenRequests, log)
	downloader := download.NewHTTPDownloader(
		repository,
		log,
		cfg.Download.ConcurrentWorkers,
		cfg.Download.RetryAttempts,
		cfg.Download.TimeoutSeconds,
	)

	downloadService := application.NewDownloadService(scraper, downloader, repository, log)

	return &CLI{
		config:          cfg,
		logger:          log,
		downloadService: downloadService,
	}, nil
}

func (c *CLI) Execute() error {
	rootCmd := &cobra.Command{
		Use:   "downloader-music",
		Short: "Downloader de clipes musicais do jw.org",
		Long: `🎵 Clipes JW.org Downloader 🎵
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Uma aplicação para baixar automaticamente todos os clipes musicais 
em português do site oficial das Testemunhas de Jeová (jw.org).

A aplicação organiza os arquivos em pastas por ano e mantém os 
nomes de arquivos limpos e organizados.`,
		Run: func(cmd *cobra.Command, args []string) {
			showBanner()
			fmt.Println("Use 'downloader-music --help' para ver os comandos disponíveis.")
			fmt.Println()
		},
	}

	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Baixa clipes musicais",
		Long:  "Baixa clipes musicais do site jw.org",
	}

	downloadAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Baixa todos os clipes disponíveis",
		Long:  "Baixa todos os clipes musicais disponíveis na página de clipes do jw.org",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.downloadAll()
		},
	}

	downloadTitleCmd := &cobra.Command{
		Use:   "title [título]",
		Short: "Baixa um clipe específico por título",
		Long:  "Baixa um clipe musical específico procurando pelo título exato",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.downloadSpecific(args[0])
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Verifica novos clipes disponíveis",
		Long:  "Verifica se há novos clipes disponíveis sem fazer download",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.checkNewClipes()
		},
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Gerencia configurações",
		Long:  "Gerencia as configurações da aplicação",
	}

	configOutputCmd := &cobra.Command{
		Use:   "output-dir [caminho]",
		Short: "Define o diretório de saída",
		Long:  "Define o diretório onde os clipes serão salvos",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.setOutputDirectory(args[0])
		},
	}

	downloadAllCmd.Flags().BoolP("verbose", "v", false, "Modo verboso")
	downloadTitleCmd.Flags().BoolP("verbose", "v", false, "Modo verboso")
	checkCmd.Flags().Bool("dry-run", true, "Apenas verificar sem baixar (sempre ativo neste comando)")

	downloadCmd.AddCommand(downloadAllCmd, downloadTitleCmd)
	configCmd.AddCommand(configOutputCmd)
	rootCmd.AddCommand(downloadCmd, checkCmd, configCmd)

	return rootCmd.Execute()
}

func (c *CLI) downloadAll() error {
	showSmallBanner()
	fmt.Println("🎵 Iniciando download de todos os clipes musicais...")
	fmt.Printf("📁 Diretório de saída: %s\n", c.config.Download.OutputDirectory)
	fmt.Printf("👥 Workers concorrentes: %d\n", c.config.Download.ConcurrentWorkers)
	fmt.Println()

	err := c.downloadService.DownloadAllClipes(c.config.Scraping.BaseURL)
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return err
	}

	fmt.Println("✅ Download concluído com sucesso!")
	return nil
}

func (c *CLI) downloadSpecific(titulo string) error {
	showSmallBanner()
	fmt.Printf("🎵 Procurando clipe: %s\n", titulo)
	fmt.Printf("📁 Diretório de saída: %s\n", c.config.Download.OutputDirectory)
	fmt.Println()

	err := c.downloadService.DownloadSpecificClipe(c.config.Scraping.BaseURL, titulo)
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return err
	}

	fmt.Println("✅ Download concluído!")
	return nil
}

func (c *CLI) checkNewClipes() error {
	showSmallBanner()
	fmt.Println("🔍 Verificando novos clipes disponíveis...")
	fmt.Println()

	novosClipes, err := c.downloadService.CheckForNewClipes(c.config.Scraping.BaseURL)
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return err
	}

	if len(novosClipes) == 0 {
		fmt.Println("✅ Nenhum clipe novo encontrado. Todos os clipes já foram baixados!")
		return nil
	}

	fmt.Printf("📋 Encontrados %d novos clipes:\n\n", len(novosClipes))
	for i, clipe := range novosClipes {
		fmt.Printf("%d. %s\n", i+1, clipe.Titulo)
		if clipe.Ano > 0 {
			fmt.Printf("   Ano: %d\n", clipe.Ano)
		}
		fmt.Printf("   URL: %s\n", clipe.URL)
		fmt.Println()
	}

	fmt.Printf("💡 Execute 'downloader-music download all' para baixar os novos clipes.\n")
	return nil
}

func (c *CLI) setOutputDirectory(dir string) error {
	if dir[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("erro ao obter diretório home: %w", err)
		}
		dir = filepath.Join(homeDir, dir[1:])
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	c.config.Download.OutputDirectory = dir

	configPath := filepath.Join("configs", "config.yaml")
	err = config.SaveConfig(c.config, configPath)
	if err != nil {
		return fmt.Errorf("erro ao salvar configuração: %w", err)
	}

	fmt.Printf("✅ Diretório de saída configurado: %s\n", dir)
	return nil
}
