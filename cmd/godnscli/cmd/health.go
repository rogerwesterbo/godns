package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	livenessPort  string
	readinessPort string
)

var healthCmd = &cobra.Command{
	Use:     "health",
	Aliases: []string{"h"},
	Short:   "Check DNS server health",
	Long:    `Check the health status of the GoDNS server (liveness and readiness probes)`,
	RunE:    runHealth,
}

func initHealthCommand() {
	rootCmd.AddCommand(healthCmd)

	healthCmd.Flags().StringVar(&livenessPort, "liveness-port", "14003", "Liveness probe port")
	healthCmd.Flags().StringVar(&readinessPort, "readiness-port", "14004", "Readiness probe port")
}

func runHealth(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Extract host from server flag
	server, _ := cmd.Flags().GetString("server")
	host := "localhost"
	if server != "" && server != "localhost:53" {
		// Extract host part (remove port)
		if idx := len(server) - 1; idx > 0 {
			for i := len(server) - 1; i >= 0; i-- {
				if server[i] == ':' {
					host = server[:i]
					break
				}
			}
		}
	}

	// Check liveness
	livenessURL := fmt.Sprintf("http://%s:%s/health/live", host, livenessPort)
	if verbose {
		fmt.Printf("Checking liveness: %s\n", livenessURL)
	}

	livenessStatus, livenessErr := checkEndpoint(livenessURL)

	// Check readiness
	readinessURL := fmt.Sprintf("http://%s:%s/health/ready", host, readinessPort)
	if verbose {
		fmt.Printf("Checking readiness: %s\n", readinessURL)
	}

	readinessStatus, readinessErr := checkEndpoint(readinessURL)

	// Display results
	fmt.Println("\n=== Health Check Results ===")
	fmt.Printf("Liveness:  %s\n", formatStatus(livenessStatus, livenessErr))
	fmt.Printf("Readiness: %s\n", formatStatus(readinessStatus, readinessErr))

	if livenessErr != nil || readinessErr != nil {
		return fmt.Errorf("health check failed")
	}

	fmt.Println("\n✓ Server is healthy and ready")
	return nil
}

func checkEndpoint(url string) (int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp.StatusCode, nil
}

func formatStatus(statusCode int, err error) string {
	if err != nil {
		return fmt.Sprintf("❌ FAILED (%v)", err)
	}

	if statusCode == http.StatusOK {
		return fmt.Sprintf("✓ OK (%d)", statusCode)
	}

	return fmt.Sprintf("❌ ERROR (%d)", statusCode)
}
