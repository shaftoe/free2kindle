package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shaftoe/free2kindle/pkg/free2kindle/content"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/email/mailjet"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/epub"
	"github.com/shaftoe/free2kindle/pkg/free2kindle/service"
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
		if err := viper.BindEnv("kindle-email", "F2K_KINDLE_EMAIL"); err != nil {
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
		var missing []string
		if viper.GetString("kindle-email") == "" {
			missing = append(missing, "--kindle-email or F2K_KINDLE_EMAIL")
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
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("Fetching article from: %s\n", url)

	mailjetConfig := &mailjet.Config{
		APIKey:      viper.GetString("api-key"),
		APISecret:   viper.GetString("api-secret"),
		SenderEmail: viper.GetString("sender-email"),
	}

	svcCfg := &service.Config{
		Extractor:    content.NewExtractor(),
		Generator:    generator,
		Sender:       mailjet.NewSender(mailjetConfig),
		SendEmail:    sendEmail,
		GenerateEPUB: sendEmail,
		KindleEmail:  viper.GetString("kindle-email"),
		SenderEmail:  viper.GetString("sender-email"),
		Subject:      emailSubject,
		OutputPath:   outputPath,
	}

	start := time.Now()
	result, err := service.Run(ctx, svcCfg, url)
	if err != nil {
		return err
	}
	fmt.Printf("Processed in %v\n", time.Since(start))

	if verbose {
		fmt.Println("\n--- Extracted Content (HTML) ---")
		fmt.Println(result.Article.Content)
		fmt.Println("--- End of Extracted Content ---")
		fmt.Println()
	}

	if !sendEmail {
		if outputPath == "" {
			outputPath = email.GenerateFilename(result.Article)
		}

		absPath, _ := filepath.Abs(outputPath)
		fmt.Printf("\n✓ EPUB saved to: %s\n", absPath)
	} else {
		senderEmail := viper.GetString("sender-email")
		fmt.Printf("\n✓ Article sent to Kindle (%s -> %s)\n", senderEmail, svcCfg.KindleEmail)
	}

	return nil
}

func main() {
	generator = epub.NewGenerator()

	convertCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
	convertCmd.Flags().DurationVarP(&timeout, "timeout", "t", defaultTimeoutSeconds*time.Second, "Timeout for HTTP requests")
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
