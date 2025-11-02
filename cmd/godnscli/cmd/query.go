package cmd

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var (
	queryType    string
	queryTimeout int
)

var queryCmd = &cobra.Command{
	Use:     "query [domain]",
	Aliases: []string{"q"},
	Short:   "Query DNS records",
	Long:    `Query DNS records for a specific domain name`,
	Args:    cobra.ExactArgs(1),
	RunE:    runQuery,
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&queryType, "type", "t", "A", "Query type (A, AAAA, MX, NS, TXT, etc.)")
	queryCmd.Flags().IntVar(&queryTimeout, "timeout", 5, "Query timeout in seconds")
}

func runQuery(cmd *cobra.Command, args []string) error {
	domain := args[0]
	server, _ := cmd.Flags().GetString("server")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Ensure domain is FQDN
	if domain[len(domain)-1] != '.' {
		domain = domain + "."
	}

	// Parse query type
	qtype, ok := dns.StringToType[queryType]
	if !ok {
		return fmt.Errorf("invalid query type: %s", queryType)
	}

	// Create DNS message
	m := new(dns.Msg)
	m.SetQuestion(domain, qtype)
	m.RecursionDesired = true

	// Create DNS client
	c := new(dns.Client)
	c.Timeout = time.Duration(queryTimeout) * time.Second

	if verbose {
		fmt.Printf("Querying %s for %s record of %s\n", server, queryType, domain)
	}

	// Send query
	r, rtt, err := c.Exchange(m, server)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Display results
	fmt.Printf("\n;; Query time: %v\n", rtt)
	fmt.Printf(";; SERVER: %s\n", server)
	fmt.Printf(";; WHEN: %s\n", time.Now().Format(time.RFC1123))
	fmt.Printf(";; MSG SIZE rcvd: %d\n\n", r.Len())

	if len(r.Answer) == 0 {
		fmt.Println(";; ANSWER SECTION:")
		fmt.Println(";; No answers found")
		return nil
	}

	fmt.Println(";; ANSWER SECTION:")
	for _, ans := range r.Answer {
		fmt.Println(ans.String())
	}

	if verbose && len(r.Ns) > 0 {
		fmt.Println("\n;; AUTHORITY SECTION:")
		for _, ns := range r.Ns {
			fmt.Println(ns.String())
		}
	}

	if verbose && len(r.Extra) > 0 {
		fmt.Println("\n;; ADDITIONAL SECTION:")
		for _, extra := range r.Extra {
			fmt.Println(extra.String())
		}
	}

	return nil
}
