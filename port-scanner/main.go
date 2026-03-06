package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"port-scanner/cmd"
)

// Version - set at compile time via -ldflags
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// If arguments provided, use classic CLI mode
	if len(os.Args) >= 2 {
		command := os.Args[1]
		switch command {
		case "scan":
			cmd.RunScan(os.Args[2:])
		case "discover":
			cmd.RunDiscover(os.Args[2:])
		case "version", "-v", "--version":
			printVersion()
		case "help", "-h", "--help":
			printUsage()
		default:
			fmt.Printf("Unknown command: %s\n", command)
			printUsage()
			os.Exit(1)
		}
		return
	}

	// Interactive mode
	runInteractiveMode()
}

func printVersion() {
	fmt.Printf("Port Scanner v%s\n", Version)
	fmt.Printf("Build time: %s\n", BuildTime)
	fmt.Printf("Git commit: %s\n", GitCommit)
}

func runInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printBanner()
		printMenu()

		fmt.Print("\n> Select option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			runInteractiveScan(reader)
		case "2":
			runInteractiveDiscover(reader)
		case "3":
			runQuickScan(reader)
		case "4":
			printUsage()
			fmt.Print("\nPress ENTER to continue...")
			reader.ReadString('\n')
		case "0":
			fmt.Println("\nGoodbye!")
			os.Exit(0)
		default:
			fmt.Println("\nUnknown option. Please try again.")
			fmt.Print("Press ENTER to continue...")
			reader.ReadString('\n')
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner() {
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║   ██████╗  ██████╗ ██████╗ ████████╗    ███████╗ ██████╗ █████╗ ███╗   ██╗   ║
║   ██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝    ██╔════╝██╔════╝██╔══██╗████╗  ██║   ║
║   ██████╔╝██║   ██║██████╔╝   ██║       ███████╗██║     ███████║██╔██╗ ██║   ║
║   ██╔═══╝ ██║   ██║██╔══██╗   ██║       ╚════██║██║     ██╔══██║██║╚██╗██║   ║
║   ██║     ╚██████╔╝██║  ██║   ██║       ███████║╚██████╗██║  ██║██║ ╚████║   ║
║   ╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝       ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝   ║
║                                                                              ║`)
	versionLine := fmt.Sprintf("Fast Concurrent Port Scanner v%s", Version)
	padding := 78 - len(versionLine)
	leftPad := padding / 2
	rightPad := padding - leftPad
	fmt.Printf("║%s%s%s║\n", strings.Repeat(" ", leftPad), versionLine, strings.Repeat(" ", rightPad))
	fmt.Println(`╚══════════════════════════════════════════════════════════════════════════════╝`)
}

func printMenu() {
	fmt.Println(`
┌──────────────────────────────────────────────────────────────────────────────┐
│                               MAIN MENU                                      │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   [1]  Port Scan         - Scan ports on a host/network                      │
│   [2]  Network Discovery - Discover active hosts (ping sweep)                │
│   [3]  Quick Scan        - Fast scan of common ports                         │
│   [4]  Help              - Show usage instructions                           │
│   [0]  Exit              - Close program                                     │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘`)
}

func runInteractiveScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                               PORT SCAN                                      ║
╚══════════════════════════════════════════════════════════════════════════════╝
`)

	// Target
	fmt.Print("> Enter target address (IP, hostname or CIDR): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)
	if target == "" {
		fmt.Println("Error: Target address is required!")
		fmt.Print("Press ENTER to continue...")
		reader.ReadString('\n')
		return
	}

	// Ports
	fmt.Print("> Enter ports [default: 1-1024] (e.g. 22,80,443 or 1-1000): ")
	ports, _ := reader.ReadString('\n')
	ports = strings.TrimSpace(ports)
	if ports == "" {
		ports = "1-1024"
	}

	// Concurrency
	fmt.Print("> Number of concurrent connections [default: 100]: ")
	concurrency, _ := reader.ReadString('\n')
	concurrency = strings.TrimSpace(concurrency)
	if concurrency == "" {
		concurrency = "100"
	}

	// Timeout
	fmt.Print("> Timeout in seconds [default: 2]: ")
	timeout, _ := reader.ReadString('\n')
	timeout = strings.TrimSpace(timeout)
	if timeout == "" {
		timeout = "2s"
	} else {
		timeout = timeout + "s"
	}

	// Output format
	fmt.Print("> Output format (table/json) [default: table]: ")
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)
	if output == "" {
		output = "table"
	}

	// Banner grabbing
	fmt.Print("> Banner grabbing? (yes/no) [default: yes]: ")
	bannerInput, _ := reader.ReadString('\n')
	bannerInput = strings.TrimSpace(strings.ToLower(bannerInput))
	banner := "true"
	if bannerInput == "no" || bannerInput == "n" {
		banner = "false"
	}

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Starting scan...")
	fmt.Println(strings.Repeat("─", 80))

	// Run scan
	args := []string{
		"-target", target,
		"-ports", ports,
		"-concurrency", concurrency,
		"-timeout", timeout,
		"-output", output,
		"-banner=" + banner,
	}
	cmd.RunScan(args)

	fmt.Print("\nScan complete. Press ENTER to continue...")
	reader.ReadString('\n')
}

func runInteractiveDiscover(reader *bufio.Reader) {
	clearScreen()
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                           NETWORK DISCOVERY                                  ║
╚══════════════════════════════════════════════════════════════════════════════╝
`)

	// Target
	fmt.Print("> Enter CIDR range (e.g. 192.168.1.0/24): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)
	if target == "" {
		fmt.Println("Error: CIDR range is required!")
		fmt.Print("Press ENTER to continue...")
		reader.ReadString('\n')
		return
	}

	// Concurrency
	fmt.Print("> Number of concurrent ping requests [default: 50]: ")
	concurrency, _ := reader.ReadString('\n')
	concurrency = strings.TrimSpace(concurrency)
	if concurrency == "" {
		concurrency = "50"
	}

	// Timeout
	fmt.Print("> Timeout in seconds [default: 2]: ")
	timeout, _ := reader.ReadString('\n')
	timeout = strings.TrimSpace(timeout)
	if timeout == "" {
		timeout = "2s"
	} else {
		timeout = timeout + "s"
	}

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Starting network discovery...")
	fmt.Println(strings.Repeat("─", 80))

	args := []string{
		"-target", target,
		"-concurrency", concurrency,
		"-timeout", timeout,
	}
	cmd.RunDiscover(args)

	fmt.Print("\nDiscovery complete. Press ENTER to continue...")
	reader.ReadString('\n')
}

func runQuickScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                              QUICK SCAN                                      ║
╚══════════════════════════════════════════════════════════════════════════════╝

Quick Scan checks the most commonly used ports:
  22 (SSH), 80 (HTTP), 443 (HTTPS), 21 (FTP), 25 (SMTP),
  3306 (MySQL), 5432 (PostgreSQL), 8080 (HTTP-Proxy), 3389 (RDP)
`)

	fmt.Print("> Enter target address (IP or hostname): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)
	if target == "" {
		fmt.Println("Error: Target address is required!")
		fmt.Print("Press ENTER to continue...")
		reader.ReadString('\n')
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Starting quick scan...")
	fmt.Println(strings.Repeat("─", 80))

	args := []string{
		"-target", target,
		"-ports", "21,22,23,25,53,80,110,143,443,445,993,995,1433,3306,3389,5432,5900,8080,8443,27017",
		"-concurrency", "50",
		"-timeout", "1s",
		"-banner=true",
	}
	cmd.RunScan(args)

	fmt.Print("\nQuick scan complete. Press ENTER to continue...")
	reader.ReadString('\n')
}

func printUsage() {
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                                 HELP                                         ║
╚══════════════════════════════════════════════════════════════════════════════╝

INTERACTIVE MODE:
    Run the program without arguments for interactive menu.

COMMAND-LINE MODE:
    port-scanner scan -target <address> [options]
    port-scanner discover -target <CIDR> [options]

SCAN OPTIONS:
    -target      Target address (IP, hostname, CIDR)
    -ports       Ports (e.g. 22,80,443 or 1-1000)
    -concurrency Number of concurrent connections
    -timeout     Timeout per connection
    -output      Format: table, json
    -banner      Banner grabbing (true/false)
    -rate-limit  Rate limit (packets/s)

DISCOVER OPTIONS:
    -target      CIDR range
    -concurrency Number of concurrent pings
    -timeout     Ping timeout

EXAMPLES:
    port-scanner scan -target 192.168.1.1 -ports 1-1000
    port-scanner scan -target 192.168.1.0/24 -ports 22,80,443 -output json
    port-scanner discover -target 192.168.1.0/24
`)
}
