package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:     "discover",
	Aliases: []string{"d", "find"},
	Short:   "Discover DNS server information",
	Long:    `Discover the DNS server's IP address and help find domains to query`,
	RunE:    runDiscover,
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

func runDiscover(cmd *cobra.Command, args []string) error {
	server, _ := cmd.Flags().GetString("server")
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Println("🔍 GoDNS Server Discovery")
	fmt.Println("==========================")
	fmt.Println()

	// Show configured server
	fmt.Printf("Configured Server: %s\n\n", server)

	// Try to resolve the hostname if it's not just an IP
	host, port, err := net.SplitHostPort(server)
	if err != nil {
		host = server
		port = "53"
	}

	// Try to get server info
	fmt.Println("📍 Server Information:")
	fmt.Println("----------------------")

	// If it's localhost, show local IPs
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		fmt.Println("Type:     Local DNS server")
		fmt.Printf("Address:  %s:%s\n", host, port)

		// Show local network interfaces
		if verbose {
			fmt.Println("\n🌐 Local Network Interfaces:")
			showLocalIPs()
		}
	} else {
		// Try to resolve the hostname
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Printf("Type:     Remote DNS server (unresolved)\n")
			fmt.Printf("Address:  %s:%s\n", host, port)
		} else {
			fmt.Printf("Type:     Remote DNS server\n")
			fmt.Printf("Address:  %s:%s\n", host, port)
			if verbose && len(ips) > 0 {
				fmt.Println("\nResolved IPs:")
				for _, ip := range ips {
					fmt.Printf("  - %s\n", ip)
				}
			}
		}
	}

	// Show how to find what to query
	fmt.Println("\n📝 Finding Domains to Query:")
	fmt.Println("----------------------------")
	fmt.Println("GoDNS stores zones in Valkey (Redis). To find available domains:")
	fmt.Println()
	fmt.Println("1. Connect to Valkey:")
	fmt.Println("   docker-compose exec valkey valkey-cli")
	fmt.Println()
	fmt.Println("2. Authenticate (if password is set):")
	fmt.Println("   AUTH default mysecretpassword")
	fmt.Println()
	fmt.Println("3. List all DNS zones:")
	fmt.Println("   KEYS dns:zone:*")
	fmt.Println()
	fmt.Println("4. View a specific zone:")
	fmt.Println("   GET dns:zone:example.lan.")
	fmt.Println()

	// Show network discovery tips
	fmt.Println("🔎 Network Discovery Tips:")
	fmt.Println("--------------------------")
	fmt.Println("When on an unknown network:")
	fmt.Println()
	fmt.Println("1. Find your gateway (router):")
	fmt.Println("   # macOS/Linux")
	fmt.Println("   netstat -rn | grep default")
	fmt.Println("   # or")
	fmt.Println("   ip route | grep default")
	fmt.Println()
	fmt.Println("2. Find DNS servers in use:")
	fmt.Println("   # macOS")
	fmt.Println("   scutil --dns | grep 'nameserver'")
	fmt.Println("   # Linux")
	fmt.Println("   cat /etc/resolv.conf")
	fmt.Println()
	fmt.Println("3. Scan for DNS servers on local network:")
	fmt.Println("   nmap -p 53 192.168.1.0/24")
	fmt.Println()
	fmt.Println("4. Test if GoDNS is running locally:")
	fmt.Println("   ./bin/godnscli h")
	fmt.Println()

	// Show example queries
	fmt.Println("💡 Example Test Queries:")
	fmt.Println("------------------------")
	fmt.Println("Once you know a domain exists in GoDNS, try:")
	fmt.Println()
	fmt.Println("  # Query local server")
	fmt.Println("  ./bin/godnscli q example.lan")
	fmt.Println()
	fmt.Println("  # Query specific server on network")
	fmt.Println("  ./bin/godnscli q example.lan -s 192.168.1.100:53")
	fmt.Println()
	fmt.Println("  # Test external resolution (should work if upstream is configured)")
	fmt.Println("  ./bin/godnscli q google.com")
	fmt.Println()

	// Show Docker-specific tips
	if verbose {
		fmt.Println("🐳 Docker-Specific Tips:")
		fmt.Println("------------------------")
		fmt.Println("If running GoDNS in Docker:")
		fmt.Println()
		fmt.Println("1. Find the container IP:")
		fmt.Println("   docker inspect godns-godns-1 | grep IPAddress")
		fmt.Println()
		fmt.Println("2. Check if port 53 is exposed:")
		fmt.Println("   docker-compose ps")
		fmt.Println()
		fmt.Println("3. View logs for any DNS queries:")
		fmt.Println("   docker-compose logs -f godns")
		fmt.Println()
	}

	// Show quick reference for creating test zones
	fmt.Println("🛠️  Creating Test Zones:")
	fmt.Println("------------------------")
	fmt.Println("To add test domains to query, you'll need to add zones to Valkey.")
	fmt.Println("See the documentation for how to create DNS zones.")
	fmt.Println()

	return nil
}

func showLocalIPs() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("  Error listing interfaces: %v\n", err)
		return
	}

	for _, iface := range ifaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		if len(addrs) > 0 {
			fmt.Printf("  %s:\n", iface.Name)
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				// Only show IPv4 and global IPv6 addresses
				if ip != nil && (ip.To4() != nil || !ip.IsLinkLocalUnicast()) {
					fmt.Printf("    - %s\n", ip)
				}
			}
		}
	}
}
