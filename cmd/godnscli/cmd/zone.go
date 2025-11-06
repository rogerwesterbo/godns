package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/rogerwesterbo/godns/cmd/godnscli/settings"
	"github.com/spf13/cobra"
)

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage DNS zones",
	Long:  `Create, list, view, update, and delete DNS zones via the GoDNS HTTP API.`,
}

var zoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS zones",
	Long:  `List all DNS zones from the GoDNS HTTP API.`,
	RunE:  runZoneList,
}

var zoneGetCmd = &cobra.Command{
	Use:   "get [domain]",
	Short: "Get a specific DNS zone",
	Long:  `Get details of a specific DNS zone including all its records.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runZoneGet,
}

var zoneDeleteCmd = &cobra.Command{
	Use:   "delete [domain]",
	Short: "Delete a DNS zone",
	Long:  `Delete a DNS zone and all its records.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runZoneDelete,
}

func init() {
	rootCmd.AddCommand(zoneCmd)
	zoneCmd.AddCommand(zoneListCmd)
	zoneCmd.AddCommand(zoneGetCmd)
	zoneCmd.AddCommand(zoneDeleteCmd)

	// Add API URL flag to zone commands
	zoneCmd.PersistentFlags().String("api-url", "", "GoDNS API URL (default from config)")
}

func getAPIURL(cmd *cobra.Command) string {
	apiURL, _ := cmd.Flags().GetString("api-url")
	if apiURL == "" {
		apiURL = settings.GetAPIURL()
	}
	return apiURL
}

func getAuthToken() (string, error) {
	return settings.GetAccessToken()
}

func makeAPIRequest(method, url string, body io.Reader) (*http.Response, error) {
	token, err := getAuthToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func runZoneList(cmd *cobra.Command, args []string) error {
	apiURL := getAPIURL(cmd)
	url := fmt.Sprintf("%s/api/v1/zones", apiURL)

	resp, err := makeAPIRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	var zones []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&zones); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(zones) == 0 {
		fmt.Println("No zones found")
		return nil
	}

	// Display zones in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "DOMAIN\tRECORDS")
	for _, zone := range zones {
		domain := zone["domain"]
		recordCount := 0
		if records, ok := zone["records"].([]interface{}); ok {
			recordCount = len(records)
		}
		_, _ = fmt.Fprintf(w, "%s\t%d\n", domain, recordCount)
	}
	_ = w.Flush()

	return nil
}

func runZoneGet(cmd *cobra.Command, args []string) error {
	domain := args[0]
	apiURL := getAPIURL(cmd)
	url := fmt.Sprintf("%s/api/v1/zones/%s", apiURL, domain)

	resp, err := makeAPIRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	var zone map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&zone); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Pretty print the zone
	data, err := json.MarshalIndent(zone, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format response: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func runZoneDelete(cmd *cobra.Command, args []string) error {
	domain := args[0]
	apiURL := getAPIURL(cmd)
	url := fmt.Sprintf("%s/api/v1/zones/%s", apiURL, domain)

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete zone '%s'? (yes/no): ", domain)
	var confirm string
	_, _ = fmt.Scanln(&confirm)
	if confirm != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	resp, err := makeAPIRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✓ Zone '%s' deleted successfully\n", domain)
	return nil
}
