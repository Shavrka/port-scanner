package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"port-scanner/cmd"
)

func main() {
	// Ako ima argumente, koristi klasični CLI mod
	if len(os.Args) >= 2 {
		command := os.Args[1]
		switch command {
		case "scan":
			cmd.RunScan(os.Args[2:])
		case "discover":
			cmd.RunDiscover(os.Args[2:])
		case "help", "-h", "--help":
			printUsage()
		default:
			fmt.Printf("Nepoznata komanda: %s\n", command)
			printUsage()
			os.Exit(1)
		}
		return
	}

	// Interaktivni mod
	runInteractiveMode()
}

func runInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printBanner()
		printMenu()

		fmt.Print("\n> Izaberi opciju: ")
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
			fmt.Print("\nPritisni ENTER za nastavak...")
			reader.ReadString('\n')
		case "0":
			fmt.Println("\nDovidjenja!")
			os.Exit(0)
		default:
			fmt.Println("\nNepoznata opcija. Pokusaj ponovo.")
			fmt.Print("Pritisni ENTER za nastavak...")
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
║                                                                              ║
║                    Brzi Concurrent Port Scanner v1.0                         ║
╚══════════════════════════════════════════════════════════════════════════════╝`)
}

func printMenu() {
	fmt.Println(`
┌──────────────────────────────────────────────────────────────────────────────┐
│                              GLAVNI MENI                                     │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   [1]  Port Scan         - Skeniraj portove na hostu/mreži                   │
│   [2]  Network Discovery - Otkrij aktivne hostove (ping sweep)               │
│   [3]  Quick Scan        - Brzo skeniranje cestih portova                    │
│   [4]  Pomoc             - Prikazi uputstva za koristenje                    │
│   [0]  Izlaz             - Zatvori program                                   │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘`)
}

func runInteractiveScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                           PORT SCAN                                          ║
╚══════════════════════════════════════════════════════════════════════════════╝
`)

	// Target
	fmt.Print("> Unesi ciljnu adresu (IP, hostname ili CIDR): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)
	if target == "" {
		fmt.Println("Greska: Ciljna adresa je obavezna!")
		fmt.Print("Pritisni ENTER za nastavak...")
		reader.ReadString('\n')
		return
	}

	// Ports
	fmt.Print("> Unesi portove [default: 1-1024] (npr. 22,80,443 ili 1-1000): ")
	ports, _ := reader.ReadString('\n')
	ports = strings.TrimSpace(ports)
	if ports == "" {
		ports = "1-1024"
	}

	// Concurrency
	fmt.Print("> Broj istovremenih konekcija [default: 100]: ")
	concurrency, _ := reader.ReadString('\n')
	concurrency = strings.TrimSpace(concurrency)
	if concurrency == "" {
		concurrency = "100"
	}

	// Timeout
	fmt.Print("> Timeout u sekundama [default: 2]: ")
	timeout, _ := reader.ReadString('\n')
	timeout = strings.TrimSpace(timeout)
	if timeout == "" {
		timeout = "2s"
	} else {
		timeout = timeout + "s"
	}

	// Output format
	fmt.Print("> Format izlaza (table/json) [default: table]: ")
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)
	if output == "" {
		output = "table"
	}

	// Banner grabbing
	fmt.Print("> Banner grabbing? (da/ne) [default: da]: ")
	bannerInput, _ := reader.ReadString('\n')
	bannerInput = strings.TrimSpace(strings.ToLower(bannerInput))
	banner := "true"
	if bannerInput == "ne" || bannerInput == "n" {
		banner = "false"
	}

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Pokrecem skeniranje...")
	fmt.Println(strings.Repeat("─", 80))

	// Pokreni scan
	args := []string{
		"-target", target,
		"-ports", ports,
		"-concurrency", concurrency,
		"-timeout", timeout,
		"-output", output,
		"-banner=" + banner,
	}
	cmd.RunScan(args)

	fmt.Print("\nSkeniranje zavrseno. Pritisni ENTER za nastavak...")
	reader.ReadString('\n')
}

func runInteractiveDiscover(reader *bufio.Reader) {
	clearScreen()
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                        NETWORK DISCOVERY                                     ║
╚══════════════════════════════════════════════════════════════════════════════╝
`)

	// Target
	fmt.Print("> Unesi CIDR range (npr. 192.168.1.0/24): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)
	if target == "" {
		fmt.Println("Greska: CIDR range je obavezan!")
		fmt.Print("Pritisni ENTER za nastavak...")
		reader.ReadString('\n')
		return
	}

	// Concurrency
	fmt.Print("> Broj istovremenih ping zahtjeva [default: 50]: ")
	concurrency, _ := reader.ReadString('\n')
	concurrency = strings.TrimSpace(concurrency)
	if concurrency == "" {
		concurrency = "50"
	}

	// Timeout
	fmt.Print("> Timeout u sekundama [default: 2]: ")
	timeout, _ := reader.ReadString('\n')
	timeout = strings.TrimSpace(timeout)
	if timeout == "" {
		timeout = "2s"
	} else {
		timeout = timeout + "s"
	}

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Pokrecem network discovery...")
	fmt.Println(strings.Repeat("─", 80))

	args := []string{
		"-target", target,
		"-concurrency", concurrency,
		"-timeout", timeout,
	}
	cmd.RunDiscover(args)

	fmt.Print("\nDiscovery zavrsen. Pritisni ENTER za nastavak...")
	reader.ReadString('\n')
}

func runQuickScan(reader *bufio.Reader) {
	clearScreen()
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                          QUICK SCAN                                          ║
╚══════════════════════════════════════════════════════════════════════════════╝

Quick Scan skenira najčešće korištene portove:
  22 (SSH), 80 (HTTP), 443 (HTTPS), 21 (FTP), 25 (SMTP), 
  3306 (MySQL), 5432 (PostgreSQL), 8080 (HTTP-Proxy), 3389 (RDP)
`)

	fmt.Print("> Unesi ciljnu adresu (IP ili hostname): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)
	if target == "" {
		fmt.Println("Greska: Ciljna adresa je obavezna!")
		fmt.Print("Pritisni ENTER za nastavak...")
		reader.ReadString('\n')
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("Pokrecem quick scan...")
	fmt.Println(strings.Repeat("─", 80))

	args := []string{
		"-target", target,
		"-ports", "21,22,23,25,53,80,110,143,443,445,993,995,1433,3306,3389,5432,5900,8080,8443,27017",
		"-concurrency", "50",
		"-timeout", "1s",
		"-banner=true",
	}
	cmd.RunScan(args)

	fmt.Print("\nQuick scan zavrsen. Pritisni ENTER za nastavak...")
	reader.ReadString('\n')
}

func printUsage() {
	fmt.Println(`
╔══════════════════════════════════════════════════════════════════════════════╗
║                              POMOĆ                                           ║
╚══════════════════════════════════════════════════════════════════════════════╝

INTERAKTIVNI MOD:
    Pokreni program bez argumenata za interaktivni meni.
    
COMMAND-LINE MOD:
    port-scanner scan -target <adresa> [opcije]
    port-scanner discover -target <CIDR> [opcije]

SCAN OPCIJE:
    -target      Ciljna adresa (IP, hostname, CIDR)
    -ports       Portovi (npr. 22,80,443 ili 1-1000)
    -concurrency Broj istovremenih konekcija
    -timeout     Timeout po konekciji
    -output      Format: table, json
    -banner      Banner grabbing (true/false)
    -rate-limit  Rate limit (paketa/s)

DISCOVER OPCIJE:
    -target      CIDR range
    -concurrency Broj istovremenih pingova
    -timeout     Timeout za ping

PRIMJERI:
    port-scanner scan -target 192.168.1.1 -ports 1-1000
    port-scanner scan -target 192.168.1.0/24 -ports 22,80,443 -output json
    port-scanner discover -target 192.168.1.0/24
`)
}
