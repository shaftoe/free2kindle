// Package main implements the CLI application for converting web articles to EPUB.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shaftoe/free2kindle/internal/config"
	"github.com/shaftoe/free2kindle/internal/content"
	"github.com/shaftoe/free2kindle/internal/email"
	"github.com/shaftoe/free2kindle/internal/email/mailjet"
	"github.com/shaftoe/free2kindle/internal/epub"
	"github.com/shaftoe/free2kindle/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultTimeoutSeconds = 30

var (
	outputPath string
	timeout    time.Duration
	verbose    bool
	generator  *epub.Generator

	sendEmail    bool
	emailSubject string
)

var rootCmd = &cobra.Command{
	Use:   "free2kindle",
	Short: "Convert web articles to EPUB format",
	Long:  `A CLI tool to fetch web articles and convert them to EPUB format for Kindle devices.`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		viper.SetEnvPrefix("F2K")
		viper.AutomaticEnv()

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}
		return nil
	},
}

var convertCmd = &cobra.Command{
	Use:   "convert [url]",
	Short: "Convert a URL to EPUB",
	Long: `Fetch a web article from given URL and convert it to EPUB format.
 Use --send to skip local EPUB generation and send converted EPUB to your Kindle.`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func runConvert(_ *cobra.Command, args []string) error {
	url := args[0]

	cfg, err := config.Load(config.ModeCLI)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Fetching article from: %s\n", url)

	var sender email.Sender
	if sendEmail {
		sender = mailjet.NewSender(cfg.MailjetAPIKey, cfg.MailjetAPISecret, cfg.SenderEmail)
	}

	d := service.NewDeps(
		content.NewExtractor(),
		generator,
		sender,
	)

	svc := service.New(d)

	start := time.Now()
	result, err := svc.Process(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to process article: %w", err)
	}
	fmt.Printf("Processed in %v\n", time.Since(start))

	printVerboseOutput(result)

	var resp *email.SendEmailResponse
	if sendEmail {
		resp, err = svc.Send(ctx, cfg, result, emailSubject)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	} else if outputPath != "" {
		if writeErr := svc.WriteToFile(result, outputPath); writeErr != nil {
			return fmt.Errorf("failed to write EPUB: %w", writeErr)
		}
	}

	printResult(result, cfg, resp)

	return nil
}

func printVerboseOutput(result *service.ProcessResult) {
	if verbose {
		fmt.Println("\n--- Extracted Content (HTML) ---")
		fmt.Println(result.Article().Content)
		fmt.Println("--- End of Extracted Content ---")
		fmt.Println()
	}
}

func printResult(result *service.ProcessResult, cfg *config.Config, resp *email.SendEmailResponse) {
	if !sendEmail {
		if outputPath == "" {
			outputPath = email.GenerateFilename(result.Article())
		}
		absPath, _ := filepath.Abs(outputPath)
		fmt.Printf("\n✓ EPUB saved to: %s\n", absPath)
	} else {
		fmt.Printf("\n✓ Article sent to Kindle (%s -> %s, email ID: %s)\n", cfg.SenderEmail, cfg.DestEmail, resp.EmailUUID)
	}
}

func main() {
	generator = epub.NewGenerator()

	convertCmd.Flags().StringVarP(&outputPath, "output", "o", "article.epub", "Output file path")
	convertCmd.Flags().DurationVarP(&timeout, "timeout", "t",
		defaultTimeoutSeconds*time.Second, "Timeout for HTTP requests")
	convertCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show extracted HTML content")

	convertCmd.Flags().BoolVar(&sendEmail, "send", false, "Send EPUB to Kindle via email instead of saving locally")
	convertCmd.Flags().StringVar(&emailSubject, "email-subject", "", "Email subject (defaults to article title)")

	rootCmd.AddCommand(convertCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
