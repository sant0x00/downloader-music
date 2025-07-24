package main

import (
	"fmt"
	"os"

	"github.com/sant0x00/downloader-music/internal/interfaces/cli"
)

func main() {
	cliApp, err := cli.NewCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar aplicação: %v\n", err)
		os.Exit(1)
	}

	if err := cliApp.Execute(); err != nil {
		os.Exit(1)
	}
}
