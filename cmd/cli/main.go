package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/spf13/cobra"
)

const DEFAULT_TIMEOUT_SECONDS = 30

var (
	outputPath string
	timeout    time.Duration
	verbose    bool
	extractor  *content.Extractor
	generator  *epub.Generator
)

func init() {
	extractor = content.NewExtractor()
	generator = epub.NewGenerator()
}

var rootCmd = &cobra.Command{
	Use:   "free2kindle",
	Short: "Convert web articles to EPUB format",
	Long:  `A CLI tool to fetch web articles and convert them to EPUB format for Kindle devices.`,
}

var convertCmd = &cobra.Command{
	Use:   "convert [url]",
	Short: "Convert a URL to EPUB",
	Long: `Fetch a web article from the given URL and convert it to EPUB format.
Use -v flag to see the extracted HTML content before conversion.`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	url := args[0]

	if timeout == 0 {
		timeout = DEFAULT_TIMEOUT_SECONDS * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Fetching article from: %s\n", url)

	start := time.Now()
	article, err := extractor.ExtractFromURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to extract article: %w", err)
	}
	fmt.Printf("Extracted in %v\n", time.Since(start))

	if article.Title == "" {
		article.Title = "Untitled"
	}

	fmt.Printf("Title: %s\n", article.Title)
	if article.Author != "" {
		fmt.Printf("Author: %s\n", article.Author)
	}
	if !article.PublishedAt.IsZero() {
		fmt.Printf("Published: %s\n", article.PublishedAt.Format("2006-01-02"))
	}

	if verbose {
		fmt.Println("\n--- Extracted Content (HTML) ---")
		fmt.Println(article.Content)
		fmt.Println("--- End of Extracted Content ---")
		fmt.Println()
	}

	if outputPath == "" {
		outputPath = sanitizeFilename(article.Title) + ".epub"
	}

	fmt.Printf("Generating EPUB: %s\n", outputPath)

	start = time.Now()
	if err := generator.GenerateAndWrite(article, outputPath); err != nil {
		return fmt.Errorf("failed to generate EPUB: %w", err)
	}
	fmt.Printf("Generated in %v\n", time.Since(start))

	absPath, _ := filepath.Abs(outputPath)
	fmt.Printf("\n✓ EPUB saved to: %s\n", absPath)

	return nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	sanitized := replacer.Replace(name)
	sanitized = strings.TrimSpace(sanitized)

	if len(sanitized) > 100 {
		sanitized = sanitized[:100]
	}

	if sanitized == "" {
		sanitized = "article"
	}

	return sanitized
}

func main() {
	convertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
	convertCmd.Flags().DurationVarP(&timeout, "timeout", "t", 30*time.Second, "Timeout for HTTP requests")
	convertCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show extracted HTML content")
	rootCmd.AddCommand(convertCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
