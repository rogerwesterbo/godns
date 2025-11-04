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
	Long:    `Run a comprehensive test suite against the DNS server`,
	RunE:    runTest,
}

func initTestCommand() {
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
	server, _ := cmd.Flags().GetString("server")
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Printf("Testing DNS server at %s\n\n", server)

	tests := []struct {
		name   string
		domain string
		qtype  uint16
	}{
		{"A Record (IPv4)", "example.lan.", dns.TypeA},
		{"AAAA Record (IPv6)", "example.lan.", dns.TypeAAAA},
		{"MX Record", "example.lan.", dns.TypeMX},
		{"NS Record", "example.lan.", dns.TypeNS},
		{"External Resolution (google.com)", "google.com.", dns.TypeA},
	}

	passed := 0
	failed := 0

	for i, test := range tests {
		fmt.Printf("[%d/%d] Testing %s... ", i+1, len(tests), test.name)

		m := new(dns.Msg)
		m.SetQuestion(test.domain, test.qtype)
		m.RecursionDesired = true

		c := new(dns.Client)
		c.Timeout = 5 * time.Second

		r, rtt, err := c.Exchange(m, server)
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			if verbose {
				fmt.Printf("   Error: %v\n", err)
			}
			failed++
			continue
		}

		if len(r.Answer) == 0 {
			fmt.Printf("⚠️  NO ANSWER (%.2fms)\n", float64(rtt.Microseconds())/1000.0)
			if verbose {
				fmt.Printf("   Response code: %s\n", dns.RcodeToString[r.Rcode])
			}
			failed++
			continue
		}

		fmt.Printf("✓ PASSED (%.2fms, %d answer(s))\n", float64(rtt.Microseconds())/1000.0, len(r.Answer))
		if verbose {
			for _, ans := range r.Answer {
				fmt.Printf("   %s\n", ans.String())
			}
		}
		passed++
	}

	fmt.Printf("\n=== Test Results ===\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total:  %d\n", len(tests))

	if failed > 0 {
		return fmt.Errorf("some tests failed")
	}

	fmt.Println("\n✓ All tests passed!")
	return nil
}
