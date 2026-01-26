package epub

import (
	"fmt"
	"os"

	"github.com/go-shiori/go-epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

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

	if err := e.Write(tmpFile.Name()); err != nil {
		return nil, fmt.Errorf("failed to write EPUB: %w", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read EPUB data: %w", err)
	}

	return data, nil
}

func (g *Generator) GenerateAndWrite(article *content.Article, outputPath string) error {
	data, err := g.Generate(article)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write EPUB file: %w", err)
	}

	return nil
}
