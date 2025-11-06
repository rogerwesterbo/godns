package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage DNS records",
	Long:  `Create, list, view, update, and delete DNS records via the GoDNS HTTP API.`,
}

var recordListCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "List DNS records in a zone",
	Long:  `List all DNS records in a specific zone.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRecordList,
}

var recordGetCmd = &cobra.Command{
	Use:   "get [domain] [name] [type]",
	Short: "Get a specific DNS record",
	Long:  `Get details of a specific DNS record.`,
	Args:  cobra.ExactArgs(3),
	RunE:  runRecordGet,
}

var recordCreateCmd = &cobra.Command{
	Use:   "create [domain]",
	Short: "Create a DNS record",
	Long: `Create a new DNS record in a zone.

Examples:
  # Create an A record
  godnscli record create example.lan --name www.example.lan. --type A --value 192.168.1.100 --ttl 300

  # Create an MX record with structured fields
  godnscli record create example.lan --name example.lan. --type MX --mx-priority 10 --mx-host mail.example.lan. --ttl 300

  # Create an SRV record
  godnscli record create example.lan --name _http._tcp.example.lan. --type SRV --srv-priority 10 --srv-weight 60 --srv-port 80 --srv-target web.example.lan. --ttl 300

  # Create a CAA record
  godnscli record create example.lan --name example.lan. --type CAA --caa-flags 0 --caa-tag issue --caa-value letsencrypt.org --ttl 300`,
	Args: cobra.ExactArgs(1),
	RunE: runRecordCreate,
}

var recordUpdateCmd = &cobra.Command{
	Use:   "update [domain] [name] [type]",
	Short: "Update a DNS record",
	Long:  `Update an existing DNS record.`,
	Args:  cobra.ExactArgs(3),
	RunE:  runRecordUpdate,
}

var recordDeleteCmd = &cobra.Command{
	Use:   "delete [domain] [name] [type]",
	Short: "Delete a DNS record",
	Long:  `Delete a DNS record from a zone.`,
	Args:  cobra.ExactArgs(3),
	RunE:  runRecordDelete,
}

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.AddCommand(recordListCmd)
	recordCmd.AddCommand(recordGetCmd)
	recordCmd.AddCommand(recordCreateCmd)
	recordCmd.AddCommand(recordUpdateCmd)
	recordCmd.AddCommand(recordDeleteCmd)

	// Common flags for create/update
	for _, cmd := range []*cobra.Command{recordCreateCmd, recordUpdateCmd} {
		cmd.Flags().String("name", "", "Record name (FQDN) (required)")
		cmd.Flags().String("type", "", "Record type: A, AAAA, CNAME, ALIAS, MX, NS, TXT, PTR, SRV, SOA, CAA (required)")
		cmd.Flags().Int("ttl", 300, "Time to live in seconds")
		cmd.Flags().String("value", "", "Record value (for simple record types)")

		// MX record flags
		cmd.Flags().Int("mx-priority", 0, "MX record priority (0-65535)")
		cmd.Flags().String("mx-host", "", "MX record mail host")

		// SRV record flags
		cmd.Flags().Int("srv-priority", 0, "SRV record priority (0-65535)")
		cmd.Flags().Int("srv-weight", 0, "SRV record weight (0-65535)")
		cmd.Flags().Int("srv-port", 0, "SRV record port (0-65535)")
		cmd.Flags().String("srv-target", "", "SRV record target hostname")

		// SOA record flags
		cmd.Flags().String("soa-mname", "", "SOA primary nameserver")
		cmd.Flags().String("soa-rname", "", "SOA admin email (@ replaced with .)")
		cmd.Flags().Uint32("soa-serial", 0, "SOA serial number")
		cmd.Flags().Uint32("soa-refresh", 3600, "SOA refresh interval (seconds)")
		cmd.Flags().Uint32("soa-retry", 1800, "SOA retry interval (seconds)")
		cmd.Flags().Uint32("soa-expire", 604800, "SOA expire time (seconds)")
		cmd.Flags().Uint32("soa-minimum", 300, "SOA minimum TTL (seconds)")

		// CAA record flags
		cmd.Flags().Int("caa-flags", 0, "CAA flags (0 or 128)")
		cmd.Flags().String("caa-tag", "", "CAA tag: issue, issuewild, iodef")
		cmd.Flags().String("caa-value", "", "CAA value (CA domain or URL)")

		_ = cmd.MarkFlagRequired("name")
		_ = cmd.MarkFlagRequired("type")
	}

	// Add API URL flag
	recordCmd.PersistentFlags().String("api-url", "", "GoDNS API URL (default from config)")

	// List filter flags
	recordListCmd.Flags().String("type-filter", "", "Filter by record type")
}

func buildRecordJSON(cmd *cobra.Command) (map[string]interface{}, error) {
	name, _ := cmd.Flags().GetString("name")
	recordType, _ := cmd.Flags().GetString("type")
	ttl, _ := cmd.Flags().GetInt("ttl")
	value, _ := cmd.Flags().GetString("value")

	record := map[string]interface{}{
		"name": name,
		"type": recordType,
		"ttl":  ttl,
	}

	// Add type-specific fields based on record type
	switch recordType {
	case "MX":
		mxPriority, _ := cmd.Flags().GetInt("mx-priority")
		mxHost, _ := cmd.Flags().GetString("mx-host")
		if mxHost == "" {
			return nil, fmt.Errorf("--mx-host is required for MX records")
		}
		record["mx_priority"] = mxPriority
		record["mx_host"] = mxHost

	case "SRV":
		srvPriority, _ := cmd.Flags().GetInt("srv-priority")
		srvWeight, _ := cmd.Flags().GetInt("srv-weight")
		srvPort, _ := cmd.Flags().GetInt("srv-port")
		srvTarget, _ := cmd.Flags().GetString("srv-target")
		if srvTarget == "" {
			return nil, fmt.Errorf("--srv-target is required for SRV records")
		}
		record["srv_priority"] = srvPriority
		record["srv_weight"] = srvWeight
		record["srv_port"] = srvPort
		record["srv_target"] = srvTarget

	case "SOA":
		soaMname, _ := cmd.Flags().GetString("soa-mname")
		soaRname, _ := cmd.Flags().GetString("soa-rname")
		soaSerial, _ := cmd.Flags().GetUint32("soa-serial")
		soaRefresh, _ := cmd.Flags().GetUint32("soa-refresh")
		soaRetry, _ := cmd.Flags().GetUint32("soa-retry")
		soaExpire, _ := cmd.Flags().GetUint32("soa-expire")
		soaMinimum, _ := cmd.Flags().GetUint32("soa-minimum")

		if soaMname == "" || soaRname == "" {
			return nil, fmt.Errorf("--soa-mname and --soa-rname are required for SOA records")
		}
		if soaSerial == 0 {
			// Auto-generate serial from current date (YYYYMMDDnn format)
			now := time.Now().UTC()
			year := now.Year()
			if year < 0 {
				return nil, fmt.Errorf("invalid system year: %d", year)
			}
			month := int(now.Month())
			day := now.Day()

			toUint64 := func(label string, v int) (uint64, error) {
				if v < 0 {
					return 0, fmt.Errorf("%s cannot be negative: %d", label, v)
				}
				return uint64(v), nil
			}

			year64, err := toUint64("year", year)
			if err != nil {
				return nil, err
			}
			month64, err := toUint64("month", month)
			if err != nil {
				return nil, err
			}
			day64, err := toUint64("day", day)
			if err != nil {
				return nil, err
			}

			serialBase := year64*1_000_000 + month64*10_000 + day64*100
			proposed := serialBase + 1
			if proposed > uint64(math.MaxUint32) {
				soaSerial = math.MaxUint32
			} else {
				soaSerial = uint32(proposed)
			}
		}

		record["soa_mname"] = soaMname
		record["soa_rname"] = soaRname
		record["soa_serial"] = soaSerial
		record["soa_refresh"] = soaRefresh
		record["soa_retry"] = soaRetry
		record["soa_expire"] = soaExpire
		record["soa_minimum"] = soaMinimum

	case "CAA":
		caaFlags, _ := cmd.Flags().GetInt("caa-flags")
		caaTag, _ := cmd.Flags().GetString("caa-tag")
		caaValue, _ := cmd.Flags().GetString("caa-value")
		if caaTag == "" || caaValue == "" {
			return nil, fmt.Errorf("--caa-tag and --caa-value are required for CAA records")
		}
		record["caa_flags"] = caaFlags
		record["caa_tag"] = caaTag
		record["caa_value"] = caaValue

	default:
		// Simple record types: A, AAAA, CNAME, ALIAS, NS, TXT, PTR
		if value == "" {
			return nil, fmt.Errorf("--value is required for %s records", recordType)
		}
		record["value"] = value
	}

	return record, nil
}

func runRecordList(cmd *cobra.Command, args []string) error {
	domain := args[0]
	apiURL := getAPIURL(cmd)
	reqURL := fmt.Sprintf("%s/api/v1/zones/%s", apiURL, url.PathEscape(domain))

	resp, err := makeAPIRequest("GET", reqURL, nil)
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

	records, ok := zone["records"].([]interface{})
	if !ok || len(records) == 0 {
		fmt.Println("No records found")
		return nil
	}

	// Apply type filter if specified
	typeFilter, _ := cmd.Flags().GetString("type-filter")
	if typeFilter != "" {
		var filtered []interface{}
		for _, r := range records {
			rec := r.(map[string]interface{})
			if rec["type"] == typeFilter {
				filtered = append(filtered, r)
			}
		}
		records = filtered
	}

	// Display records in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tTYPE\tVALUE\tTTL")
	for _, r := range records {
		rec := r.(map[string]interface{})
		name := rec["name"]
		recordType := rec["type"]
		ttl := rec["ttl"]

		// Format value based on record type
		var valueStr string
		switch recordType {
		case "MX":
			if mxPriority, ok := rec["mx_priority"]; ok {
				valueStr = fmt.Sprintf("%v → %v", mxPriority, rec["mx_host"])
			} else {
				valueStr = fmt.Sprint(rec["value"])
			}
		case "SRV":
			if _, ok := rec["srv_priority"]; ok {
				valueStr = fmt.Sprintf("Pri:%v Wgt:%v Port:%v → %v",
					rec["srv_priority"], rec["srv_weight"], rec["srv_port"], rec["srv_target"])
			} else {
				valueStr = fmt.Sprint(rec["value"])
			}
		case "SOA":
			if soaMname, ok := rec["soa_mname"]; ok {
				valueStr = fmt.Sprintf("%v %v (Serial: %v)",
					soaMname, rec["soa_rname"], rec["soa_serial"])
			} else {
				valueStr = fmt.Sprint(rec["value"])
			}
		case "CAA":
			if caaTag, ok := rec["caa_tag"]; ok {
				valueStr = fmt.Sprintf("[%v] %v: %v", rec["caa_flags"], caaTag, rec["caa_value"])
			} else {
				valueStr = fmt.Sprint(rec["value"])
			}
		default:
			valueStr = fmt.Sprint(rec["value"])
		}

		// Truncate long values
		if len(valueStr) > 60 {
			valueStr = valueStr[:57] + "..."
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", name, recordType, valueStr, ttl)
	}
	_ = w.Flush()

	return nil
}

func runRecordGet(cmd *cobra.Command, args []string) error {
	domain := args[0]
	name := args[1]
	recordType := args[2]
	apiURL := getAPIURL(cmd)
	reqURL := fmt.Sprintf("%s/api/v1/zones/%s/records/%s/%s",
		apiURL, url.PathEscape(domain), url.PathEscape(name), url.PathEscape(recordType))

	resp, err := makeAPIRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	var record map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Pretty print the record
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format response: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func runRecordCreate(cmd *cobra.Command, args []string) error {
	domain := args[0]
	apiURL := getAPIURL(cmd)
	reqURL := fmt.Sprintf("%s/api/v1/zones/%s/records", apiURL, url.PathEscape(domain))

	record, err := buildRecordJSON(cmd)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to encode record: %w", err)
	}

	resp, err := makeAPIRequest("POST", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✓ Record created: %s %s\n", record["name"], record["type"])
	return nil
}

func runRecordUpdate(cmd *cobra.Command, args []string) error {
	domain := args[0]
	name := args[1]
	recordType := args[2]
	apiURL := getAPIURL(cmd)
	reqURL := fmt.Sprintf("%s/api/v1/zones/%s/records/%s/%s",
		apiURL, url.PathEscape(domain), url.PathEscape(name), url.PathEscape(recordType))

	// Override name and type from args
	_ = cmd.Flags().Set("name", name)
	_ = cmd.Flags().Set("type", recordType)

	record, err := buildRecordJSON(cmd)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to encode record: %w", err)
	}

	resp, err := makeAPIRequest("PUT", reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✓ Record updated: %s %s\n", name, recordType)
	return nil
}

func runRecordDelete(cmd *cobra.Command, args []string) error {
	domain := args[0]
	name := args[1]
	recordType := args[2]
	apiURL := getAPIURL(cmd)
	reqURL := fmt.Sprintf("%s/api/v1/zones/%s/records/%s/%s",
		apiURL, url.PathEscape(domain), url.PathEscape(name), url.PathEscape(recordType))

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete record '%s %s' from zone '%s'? (yes/no): ", name, recordType, domain)
	var confirm string
	_, _ = fmt.Scanln(&confirm)
	if confirm != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	resp, err := makeAPIRequest("DELETE", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✓ Record deleted: %s %s\n", name, recordType)
	return nil
}
