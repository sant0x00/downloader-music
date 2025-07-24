package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sant0x00/downloader-music/internal/domain"
)

type FileSystemRepository struct {
	outputDirectory string
	logger          domain.Logger
}

func NewFileSystemRepository(outputDirectory string, logger domain.Logger) *FileSystemRepository {
	return &FileSystemRepository{
		outputDirectory: outputDirectory,
		logger:          logger,
	}
}

// TODO: Implement FindAll method to list all clips in the filesystem
func (r *FileSystemRepository) FindAll() ([]domain.ClipeMusical, error) {
	return []domain.ClipeMusical{}, nil
}

func (r *FileSystemRepository) Save(clipe domain.ClipeMusical) error {
	r.logger.Debug("Salvando clipe", "titulo", clipe.Titulo, "arquivo", clipe.GetSanitizedFilename())
	return nil
}

func (r *FileSystemRepository) Exists(filename string) bool {
	dirs := []string{"2025", "2024", "2023", "2022", "outros"}

	for _, dir := range dirs {
		fullPath := filepath.Join(r.outputDirectory, dir, filename)
		if _, err := os.Stat(fullPath); err == nil {
			r.logger.Debug("Arquivo já existe", "path", fullPath)
			return true
		}
	}

	fullPath := filepath.Join(r.outputDirectory, filename)
	if _, err := os.Stat(fullPath); err == nil {
		r.logger.Debug("Arquivo já existe", "path", fullPath)
		return true
	}

	return false
}

func (r *FileSystemRepository) GetOutputDirectory() string {
	return r.outputDirectory
}

func (r *FileSystemRepository) CreateDirectoryStructure(clipe domain.ClipeMusical) error {
	var targetDir string

	if clipe.Ano > 0 {
		targetDir = filepath.Join(r.outputDirectory, strconv.Itoa(clipe.Ano))
	} else {
		targetDir = filepath.Join(r.outputDirectory, "outros")
	}

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		r.logger.Error("Erro ao criar diretório", err, "path", targetDir)
		return fmt.Errorf("erro ao criar diretório %s: %w", targetDir, err)
	}

	r.logger.Debug("Diretório criado/verificado", "path", targetDir)
	return nil
}

func (r *FileSystemRepository) GetClipeFilePath(clipe domain.ClipeMusical) string {
	var targetDir string

	if clipe.Ano > 0 {
		targetDir = filepath.Join(r.outputDirectory, strconv.Itoa(clipe.Ano))
	} else {
		targetDir = filepath.Join(r.outputDirectory, "outros")
	}

	return filepath.Join(targetDir, clipe.GetSanitizedFilename())
}
