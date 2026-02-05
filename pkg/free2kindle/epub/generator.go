// Package epub provides EPUB file generation functionality.
package epub

import (
	"fmt"
	"os"

	"github.com/go-shiori/go-epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
)

const (
	// fileModeReadWrite is the file permission for EPUB files (readable by user).
	fileModeReadWrite = 0o644
)

// Generator handles EPUB file generation from article content.
type Generator struct{}

// NewGenerator creates a new EPUB generator instance.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate creates an EPUB file from the given article and returns its bytes.
func (g *Generator) Generate(article *content.Article) ([]byte, error) {
	e, err := epub.NewEpub(article.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create EPUB: %w", err)
	}

	if article.Title != "" {
		e.SetTitle(article.Title)
	}

	if article.Author != "" {
		e.SetAuthor(article.Author)
	}

	if article.Excerpt != "" {
		e.SetDescription(article.Excerpt)
	}

	e.SetLang("en")

	_, err = e.AddSection(article.Content, "Chapter 1", "chapter1.xhtml", "")
	if err != nil {
		return nil, fmt.Errorf("failed to add chapter: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "*.epub")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if writeErr := e.Write(tmpFile.Name()); writeErr != nil {
		return nil, fmt.Errorf("failed to write EPUB: %w", writeErr)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read EPUB data: %w", err)
	}

	return data, nil
}

// GenerateAndWrite generates an EPUB file and writes it to the specified path.
func (g *Generator) GenerateAndWrite(article *content.Article, outputPath string) error {
	data, err := g.Generate(article)
	if err != nil {
		return err
	}

	// #nosec G306 - EPUB files need to be readable by user
	if writeErr := os.WriteFile(outputPath, data, fileModeReadWrite); writeErr != nil {
		return fmt.Errorf("failed to write EPUB file: %w", writeErr)
	}

	return nil
}
