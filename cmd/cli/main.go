package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DEFAULT_TIMEOUT_SECONDS = 30

var (
	outputPath string
	timeout    time.Duration
	verbose    bool
	extractor  *content.Extractor
	generator  *epub.Generator

	sendEmail    bool
	emailSubject string
)

var rootCmd = &cobra.Command{
	Use:   "free2kindle",
	Short: "Convert web articles to EPUB format",
	Long:  `A CLI tool to fetch web articles and convert them to EPUB format for Kindle devices.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("F2K")
		viper.AutomaticEnv()
		viper.BindEnv("kindle-email", "F2K_KINDLE_EMAIL")
		viper.BindEnv("sender-email", "F2K_SENDER_EMAIL")
		viper.BindEnv("api-key", "MAILJET_API_KEY", "MJ_APIKEY_PUBLIC")
		viper.BindEnv("api-secret", "MAILJET_API_SECRET", "MJ_APIKEY_PRIVATE")

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}
		return nil
	},
}

func stringJoin(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for i := 1; i < len(items); i++ {
		result += sep + items[i]
	}
	return result
}

var convertCmd = &cobra.Command{
	Use:   "convert [url]",
	Short: "Convert a URL to EPUB",
	Long: `Fetch a web article from the given URL and convert it to EPUB format.
Use --send to skip local EPUB generation and send the converted EPUB to your Kindle.`,
	Args:    cobra.ExactArgs(1),
	PreRunE: validateEmailConfig,
	RunE:    runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	url := args[0]

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

	if sendEmail {
		return sendToKindle(ctx, article)
	}

	if outputPath == "" {
		outputPath = email.GenerateFilename(article)
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

func main() {
	extractor = content.NewExtractor()
	generator = epub.NewGenerator()

	convertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
	convertCmd.Flags().DurationVarP(&timeout, "timeout", "t", DEFAULT_TIMEOUT_SECONDS*time.Second, "Timeout for HTTP requests")
	convertCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show extracted HTML content")

	convertCmd.Flags().BoolVar(&sendEmail, "send", false, "Send EPUB to Kindle via email instead of saving locally")
	convertCmd.Flags().StringVar(&emailSubject, "email-subject", "", "Email subject (defaults to article title)")

	convertCmd.Flags().String("kindle-email", "", "Kindle email address (env: F2K_KINDLE_EMAIL)")
	convertCmd.Flags().String("sender-email", "", "Sender email address (env: F2K_SENDER_EMAIL)")
	convertCmd.Flags().String("api-key", "", "Mailjet API key (env: MAILJET_API_KEY)")
	convertCmd.Flags().String("api-secret", "", "Mailjet API secret (env: MAILJET_API_SECRET)")

	rootCmd.AddCommand(convertCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
