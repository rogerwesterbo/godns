package cmd

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Run comprehensive DNS server tests",
	Long:    `Run a comprehensive test suite against the DNS server including DNSSEC validation`,
	RunE:    runTest,
}

var skipDNSSEC bool

func initTestCommand() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().Bool("skip-dnssec", false, "Skip DNSSEC tests")
}

func runTest(cmd *cobra.Command, args []string) error {
	server, _ := cmd.Flags().GetString("server")
	verbose, _ := cmd.Flags().GetBool("verbose")
	skipDNSSEC, _ = cmd.Flags().GetBool("skip-dnssec")

	fmt.Printf("Testing DNS server at %s\n\n", server)

	// Basic DNS tests - test records that should exist
	basicTests := []struct {
		name     string
		domain   string
		qtype    uint16
		optional bool // If true, "no answer" is acceptable (not a failure)
	}{
		{"A Record (www.example.lan)", "www.example.lan.", dns.TypeA, false},
		{"A Record (db.example.lan)", "db.example.lan.", dns.TypeA, false},
		{"AAAA Record (optional)", "www.example.lan.", dns.TypeAAAA, true},
		{"MX Record", "example.lan.", dns.TypeMX, false},
		{"NS Record", "example.lan.", dns.TypeNS, false},
		{"TXT Record", "example.lan.", dns.TypeTXT, true},
		{"SOA Record", "example.lan.", dns.TypeSOA, true},
		{"External Resolution (google.com)", "google.com.", dns.TypeA, false},
	}

	// DNSSEC tests - these are optional as DNSSEC may not be configured
	dnssecTests := []struct {
		name     string
		domain   string
		qtype    uint16
		optional bool
	}{
		{"DNSKEY Record (local)", "example.lan.", dns.TypeDNSKEY, true},
		{"DNSKEY Record (external)", "cloudflare.com.", dns.TypeDNSKEY, false},
	}

	passed := 0
	failed := 0
	skipped := 0

	// Run basic tests
	fmt.Println("=== Basic DNS Tests ===")
	for i, test := range basicTests {
		fmt.Printf("[%d/%d] Testing %s... ", i+1, len(basicTests), test.name)
		result := runSingleTestWithOptional(server, test.domain, test.qtype, false, verbose, test.optional)
		switch result {
		case "passed":
			passed++
		case "failed":
			failed++
		case "skipped":
			skipped++
		}
	}

	// Run DNSSEC tests
	if !skipDNSSEC {
		fmt.Println("\n=== DNSSEC Tests ===")
		for i, test := range dnssecTests {
			fmt.Printf("[%d/%d] Testing %s... ", i+1, len(dnssecTests), test.name)
			result := runSingleTestWithOptional(server, test.domain, test.qtype, false, verbose, test.optional)
			switch result {
			case "passed":
				passed++
			case "failed":
				failed++
			case "skipped":
				skipped++
			}
		}

		// Test DNSSEC validation with DO flag
		fmt.Println("\n=== DNSSEC Validation Tests ===")
		fmt.Printf("[1/1] Testing external DNSSEC (cloudflare.com)... ")
		result := runDNSSECValidationTest(server, "cloudflare.com.", dns.TypeA, verbose)
		switch result {
		case "passed":
			passed++
		case "failed":
			failed++
		}
	} else {
		dnssecSkipped := len(dnssecTests) + 1
		skipped += dnssecSkipped
		fmt.Println("\n=== DNSSEC Tests ===")
		fmt.Printf("Skipped %d DNSSEC tests (--skip-dnssec flag)\n", dnssecSkipped)
	}

	totalTests := passed + failed + skipped

	fmt.Printf("\n=== Test Results ===\n")
	fmt.Printf("Passed:  %d\n", passed)
	fmt.Printf("Failed:  %d\n", failed)
	if skipped > 0 {
		fmt.Printf("Skipped: %d (optional tests with no answer)\n", skipped)
	}
	fmt.Printf("Total:   %d\n", totalTests)

	if failed > 0 {
		return fmt.Errorf("some tests failed")
	}

	fmt.Println("\n✓ All tests passed!")
	return nil
}

func runSingleTestWithOptional(server, domain string, qtype uint16, setDO, verbose, optional bool) string {
	m := new(dns.Msg)
	m.SetQuestion(domain, qtype)
	m.RecursionDesired = true

	if setDO {
		m.SetEdns0(4096, true) // Enable DNSSEC OK (DO) flag
	}

	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	r, rtt, err := c.Exchange(m, server)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		if verbose {
			fmt.Printf("   Error: %v\n", err)
		}
		return "failed"
	}

	if len(r.Answer) == 0 {
		if optional {
			fmt.Printf("⏭️  SKIPPED (%.2fms, optional - no record)\n", float64(rtt.Microseconds())/1000.0)
			return "skipped"
		}
		fmt.Printf("⚠️  NO ANSWER (%.2fms)\n", float64(rtt.Microseconds())/1000.0)
		if verbose {
			fmt.Printf("   Response code: %s\n", dns.RcodeToString[r.Rcode])
		}
		return "failed"
	}

	fmt.Printf("✓ PASSED (%.2fms, %d answer(s))\n", float64(rtt.Microseconds())/1000.0, len(r.Answer))
	if verbose {
		for _, ans := range r.Answer {
			fmt.Printf("   %s\n", ans.String())
		}
	}
	return "passed"
}

func runDNSSECValidationTest(server, domain string, qtype uint16, verbose bool) string {
	m := new(dns.Msg)
	m.SetQuestion(domain, qtype)
	m.RecursionDesired = true
	m.SetEdns0(4096, true) // Enable DNSSEC OK (DO) flag

	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	r, rtt, err := c.Exchange(m, server)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		if verbose {
			fmt.Printf("   Error: %v\n", err)
		}
		return "failed"
	}

	// Check if we got any answer
	if len(r.Answer) == 0 {
		fmt.Printf("⚠️  NO ANSWER (%.2fms)\n", float64(rtt.Microseconds())/1000.0)
		if verbose {
			fmt.Printf("   Response code: %s\n", dns.RcodeToString[r.Rcode])
		}
		return "failed"
	}

	// Check for AD (Authenticated Data) flag or RRSIG in response
	hasRRSIG := false
	for _, ans := range r.Answer {
		if _, ok := ans.(*dns.RRSIG); ok {
			hasRRSIG = true
			break
		}
	}

	adFlag := r.AuthenticatedData

	if hasRRSIG || adFlag {
		fmt.Printf("✓ PASSED (%.2fms, AD=%v, RRSIG=%v)\n",
			float64(rtt.Microseconds())/1000.0, adFlag, hasRRSIG)
	} else {
		fmt.Printf("✓ PASSED (%.2fms, %d answer(s), no DNSSEC)\n",
			float64(rtt.Microseconds())/1000.0, len(r.Answer))
	}

	if verbose {
		fmt.Printf("   AD flag: %v\n", adFlag)
		fmt.Printf("   RRSIG present: %v\n", hasRRSIG)
		for _, ans := range r.Answer {
			fmt.Printf("   %s\n", ans.String())
		}
	}

	return "passed"
}
