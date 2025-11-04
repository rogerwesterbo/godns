package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [domain]",
	Short: "Export DNS zones to different formats",
	Long: `Export DNS zones to different DNS provider formats.

Supported formats:
  - bind      : Standard BIND zone file format (default)
  - coredns   : CoreDNS configuration format
  - powerdns  : PowerDNS JSON API format
  - zonefile  : Generic zone file (same as bind)

Examples:
	# Export all zones in BIND format
	godnscli export --api-url http://localhost:14082

	# Export all zones in CoreDNS format
	godnscli export --format coredns --api-url http://localhost:14082

	# Export a specific zone in PowerDNS format
	godnscli export example.lan --format powerdns --api-url http://localhost:14082

  # Export to file
  godnscli export example.lan --format bind --output example.lan.zone`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

var (
	exportFormat string
	exportOutput string
	apiURL       string
)

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "bind", "Export format (bind, coredns, powerdns, zonefile)")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: stdout)")
	exportCmd.Flags().StringVar(&apiURL, "api-url", "http://localhost:14082", "GoDNS API URL")
}

func runExport(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Validate format
	validFormats := map[string]bool{
		"bind":     true,
		"coredns":  true,
		"powerdns": true,
		"zonefile": true,
	}

	if !validFormats[exportFormat] {
		return fmt.Errorf("invalid format: %s. Supported formats: bind, coredns, powerdns, zonefile", exportFormat)
	}

	// Build URL
	var url string
	if len(args) == 0 {
		// Export all zones
		url = fmt.Sprintf("%s/api/v1/export?format=%s", apiURL, exportFormat)
		if verbose {
			fmt.Fprintf(os.Stderr, "Exporting all zones in %s format...\n", exportFormat)
		}
	} else {
		// Export specific zone
		domain := args[0]
		url = fmt.Sprintf("%s/api/v1/export/%s?format=%s", apiURL, domain, exportFormat)
		if verbose {
			fmt.Fprintf(os.Stderr, "Exporting zone %s in %s format...\n", domain, exportFormat)
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "URL: %s\n", url)
	}

	// Make HTTP request
	// #nosec G107 - URL is constructed from user-provided flags (--api-url and domain argument)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Write output
	if exportOutput != "" {
		// Write to file
		// #nosec G306 - Output file permissions 0644 are appropriate for DNS zone files
		err = os.WriteFile(exportOutput, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Exported to: %s\n", exportOutput)
		} else {
			fmt.Fprintf(os.Stderr, "Successfully exported to %s\n", exportOutput)
		}
	} else {
		// Write to stdout
		fmt.Print(string(data))
	}

	return nil
}
