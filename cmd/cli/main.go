// Package main implements the CLI application for converting web articles to EPUB.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		if err := viper.BindEnv("destination-email", "F2K_DEST_EMAIL"); err != nil {
			return fmt.Errorf("failed to bind env: %w", err)
		}
		if err := viper.BindEnv("sender-email", "F2K_SENDER_EMAIL"); err != nil {
			return fmt.Errorf("failed to bind env: %w", err)
		}
		if err := viper.BindEnv("api-key", "MAILJET_API_KEY", "MJ_APIKEY_PUBLIC"); err != nil {
			return fmt.Errorf("failed to bind env: %w", err)
		}
		if err := viper.BindEnv("api-secret", "MAILJET_API_SECRET", "MJ_APIKEY_PRIVATE"); err != nil {
			return fmt.Errorf("failed to bind env: %w", err)
		}

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return fmt.Errorf("failed to bind flags: %w", err)
		}
		return nil
	},
}

func validateSendEmailConfig() error {
	var missing []string
	if viper.GetString("destination-email") == "" {
		missing = append(missing, "--destination-email or F2K_DEST_EMAIL")
	}
	if viper.GetString("sender-email") == "" {
		missing = append(missing, "--sender-email or F2K_SENDER_EMAIL")
	}
	if viper.GetString("api-key") == "" {
		missing = append(missing, "--api-key or MAILJET_API_KEY")
	}
	if viper.GetString("api-secret") == "" {
		missing = append(missing, "--api-secret or MAILJET_API_SECRET")
	}

	if len(missing) > 0 {
		return fmt.Errorf("when using --send flag, the following are required: %s",
			stringJoin(missing, ", "))
	}
	return nil
}

func stringJoin(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	var resultSb63 strings.Builder
	for i := 1; i < len(items); i++ {
		resultSb63.WriteString(sep + items[i])
	}
	result += resultSb63.String()
	return result
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

	if sendEmail {
		if err := validateSendEmailConfig(); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Fetching article from: %s\n", url)

	cfg := buildConfig()

	var sender email.Sender
	if sendEmail {
		sender = mailjet.NewSender(cfg.MailjetAPIKey, cfg.MailjetAPISecret, cfg.SenderEmail)
	}

	d := service.NewDeps(
		content.NewExtractor(),
		generator,
		sender,
	)

	opts := service.NewOptions(sendEmail, true, emailSubject, outputPath)

	start := time.Now()
	result, err := service.Run(ctx, d, cfg, opts, url)
	if err != nil {
		return fmt.Errorf("failed to process article: %w", err)
	}
	fmt.Printf("Processed in %v\n", time.Since(start))

	printVerboseOutput(result)
	printResult(result, cfg)

	return nil
}

func buildConfig() *config.Config {
	return &config.Config{
		DestEmail:        viper.GetString("destination-email"),
		SenderEmail:      viper.GetString("sender-email"),
		MailjetAPIKey:    viper.GetString("api-key"),
		MailjetAPISecret: viper.GetString("api-secret"),
		SendEnabled:      viper.GetBool("send-enabled"),
	}
}

func printVerboseOutput(result *service.Result) {
	if verbose {
		fmt.Println("\n--- Extracted Content (HTML) ---")
		fmt.Println(result.Article.Content)
		fmt.Println("--- End of Extracted Content ---")
		fmt.Println()
	}
}

func printResult(result *service.Result, cfg *config.Config) {
	if !sendEmail {
		if outputPath == "" {
			outputPath = email.GenerateFilename(result.Article)
		}
		absPath, _ := filepath.Abs(outputPath)
		fmt.Printf("\n✓ EPUB saved to: %s\n", absPath)
	} else {
		senderEmail := viper.GetString("sender-email")
		fmt.Printf("\n✓ Article sent to Kindle (%s -> %s)\n", senderEmail, cfg.DestEmail)
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

	convertCmd.Flags().String("destination-email", "", "Kindle email address (env: F2K_DEST_EMAIL)")
	convertCmd.Flags().String("sender-email", "", "Sender email address (env: F2K_SENDER_EMAIL)")
	convertCmd.Flags().String("api-key", "", "Mailjet API key (env: MAILJET_API_KEY)")
	convertCmd.Flags().String("api-secret", "", "Mailjet API secret (env: MAILJET_API_SECRET)")

	rootCmd.AddCommand(convertCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
