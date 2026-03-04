package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"port-scanner/internal/network"
	"port-scanner/internal/scanner"
)

type ScanResult struct {
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	State    string        `json:"state"`
	Service  string        `json:"service"`
	Banner   string        `json:"banner,omitempty"`
	ScanTime time.Duration `json:"scan_time"`
}

type ScanReport struct {
	StartTime  time.Time    `json:"start_time"`
	EndTime    time.Time    `json:"end_time"`
	TotalHosts int          `json:"total_hosts"`
	TotalPorts int          `json:"total_ports"`
	OpenPorts  int          `json:"open_ports"`
	Results    []ScanResult `json:"results"`
}

func RunScan(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)

	target := fs.String("target", "", "Ciljna IP adresa ili CIDR range (npr. 192.168.1.1 ili 192.168.1.0/24)")
	ports := fs.String("ports", "1-1024", "Portovi za skeniranje (npr. 22,80,443 ili 1-1000)")
	concurrency := fs.Int("concurrency", 100, "Broj istovremenih konekcija")
	timeout := fs.Duration("timeout", 2*time.Second, "Timeout za konekciju")
	rateLimit := fs.Int("rate-limit", 0, "Rate limit (paketa po sekundi, 0 = bez limita)")
	outputFormat := fs.String("output", "table", "Format outputa: table, json")
	grabBanner := fs.Bool("banner", true, "Pokušaj banner grabbing")
	retries := fs.Int("retries", 1, "Broj ponovnih pokušaja")

	fs.Parse(args)

	if *target == "" {
		fmt.Println("Greška: -target je obavezan parametar")
		fs.PrintDefaults()
		os.Exit(1)
	}

	// Parse target hosts
	hosts, err := network.ParseTarget(*target)
	if err != nil {
		fmt.Printf("Greška pri parsiranju cilja: %v\n", err)
		os.Exit(1)
	}

	// Parse ports
	portList, err := network.ParsePorts(*ports)
	if err != nil {
		fmt.Printf("Greška pri parsiranju portova: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n[*] Pokrećem skeniranje...\n")
	fmt.Printf("[*] Ciljevi: %d host(ova)\n", len(hosts))
	fmt.Printf("[*] Portovi: %d\n", len(portList))
	fmt.Printf("[*] Konkurentnost: %d\n", *concurrency)
	fmt.Printf("[*] Timeout: %v\n", *timeout)
	if *rateLimit > 0 {
		fmt.Printf("[*] Rate limit: %d paketa/s\n", *rateLimit)
	}
	fmt.Println()

	report := ScanReport{
		StartTime:  time.Now(),
		TotalHosts: len(hosts),
		TotalPorts: len(portList),
		Results:    []ScanResult{},
	}

	// Create scanner
	s := scanner.NewTCPScanner(*concurrency, *timeout, *retries, *rateLimit)

	// Scan all hosts
	for _, host := range hosts {
		results := s.ScanPorts(host, portList, *grabBanner)

		for _, r := range results {
			if r.State == "open" {
				report.OpenPorts++
				result := ScanResult{
					Host:     host,
					Port:     r.Port,
					State:    r.State,
					Service:  r.Service,
					Banner:   r.Banner,
					ScanTime: r.ScanTime,
				}
				report.Results = append(report.Results, result)
			}
		}
	}

	report.EndTime = time.Now()

	// Sort results by host and port
	sort.Slice(report.Results, func(i, j int) bool {
		if report.Results[i].Host != report.Results[j].Host {
			return report.Results[i].Host < report.Results[j].Host
		}
		return report.Results[i].Port < report.Results[j].Port
	})

	// Output results
	switch strings.ToLower(*outputFormat) {
	case "json":
		outputJSON(report)
	default:
		outputTable(report)
	}
}

func outputTable(report ScanReport) {
	duration := report.EndTime.Sub(report.StartTime)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("REZULTATI SKENIRANJA")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Vrijeme početka: %s\n", report.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Vrijeme završetka: %s\n", report.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Trajanje: %v\n", duration.Round(time.Millisecond))
	fmt.Printf("Skenirano hostova: %d\n", report.TotalHosts)
	fmt.Printf("Skenirano portova: %d\n", report.TotalPorts)
	fmt.Printf("Otvorenih portova: %d\n", report.OpenPorts)
	fmt.Println(strings.Repeat("-", 80))

	if len(report.Results) == 0 {
		fmt.Println("Nisu pronađeni otvoreni portovi.")
		return
	}

	// Table header
	fmt.Printf("%-20s %-8s %-10s %-15s %s\n", "HOST", "PORT", "STATE", "SERVICE", "BANNER")
	fmt.Println(strings.Repeat("-", 80))

	currentHost := ""
	for _, r := range report.Results {
		host := r.Host
		if host == currentHost {
			host = ""
		} else {
			currentHost = r.Host
		}

		banner := r.Banner
		if len(banner) > 30 {
			banner = banner[:27] + "..."
		}
		banner = strings.ReplaceAll(banner, "\n", " ")
		banner = strings.ReplaceAll(banner, "\r", "")

		fmt.Printf("%-20s %-8d %-10s %-15s %s\n", host, r.Port, r.State, r.Service, banner)
	}

	fmt.Println(strings.Repeat("=", 80))
}

func outputJSON(report ScanReport) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Printf("Greška pri generisanju JSON-a: %v\n", err)
		return
	}
	fmt.Println(string(data))
}
