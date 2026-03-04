package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"port-scanner/internal/network"
)

type HostStatus struct {
	IP    string
	Alive bool
	RTT   time.Duration
}

func RunDiscover(args []string) {
	fs := flag.NewFlagSet("discover", flag.ExitOnError)

	target := fs.String("target", "", "CIDR range za discovery (npr. 192.168.1.0/24)")
	concurrency := fs.Int("concurrency", 50, "Broj istovremenih ping zahtjeva")
	timeout := fs.Duration("timeout", 2*time.Second, "Timeout za ping")
	outputFormat := fs.String("output", "table", "Format outputa: table, list")

	fs.Parse(args)

	if *target == "" {
		fmt.Println("Greška: -target je obavezan parametar")
		fs.PrintDefaults()
		os.Exit(1)
	}

	// Parse CIDR
	hosts, err := network.ParseTarget(*target)
	if err != nil {
		fmt.Printf("Greška pri parsiranju CIDR-a: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n[*] Pokrećem network discovery...\n")
	fmt.Printf("[*] Range: %s\n", *target)
	fmt.Printf("[*] Hostova za skeniranje: %d\n", len(hosts))
	fmt.Printf("[*] Timeout: %v\n", *timeout)
	fmt.Println()

	results := discoverHosts(hosts, *concurrency, *timeout)

	// Sort by IP
	sort.Slice(results, func(i, j int) bool {
		return compareIPs(results[i].IP, results[j].IP)
	})

	// Count alive hosts
	alive := 0
	for _, r := range results {
		if r.Alive {
			alive++
		}
	}

	switch strings.ToLower(*outputFormat) {
	case "list":
		outputList(results)
	default:
		outputDiscoveryTable(results, len(hosts), alive)
	}
}

func discoverHosts(hosts []string, concurrency int, timeout time.Duration) []HostStatus {
	var results []HostStatus
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, concurrency)

	for _, host := range hosts {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			alive, rtt := pingHost(ip, timeout)

			mu.Lock()
			results = append(results, HostStatus{
				IP:    ip,
				Alive: alive,
				RTT:   rtt,
			})
			mu.Unlock()
		}(host)
	}

	wg.Wait()
	return results
}

func pingHost(ip string, timeout time.Duration) (bool, time.Duration) {
	start := time.Now()

	// Try TCP connect to common ports first (faster than ICMP)
	commonPorts := []int{80, 443, 22, 445, 139}

	for _, port := range commonPorts {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout/time.Duration(len(commonPorts)))
		if err == nil {
			conn.Close()
			return true, time.Since(start)
		}
	}

	// Fallback to system ping
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()

	return err == nil, time.Since(start)
}

func compareIPs(ip1, ip2 string) bool {
	parts1 := strings.Split(ip1, ".")
	parts2 := strings.Split(ip2, ".")

	for i := 0; i < 4; i++ {
		var n1, n2 int
		fmt.Sscanf(parts1[i], "%d", &n1)
		fmt.Sscanf(parts2[i], "%d", &n2)
		if n1 != n2 {
			return n1 < n2
		}
	}
	return false
}

func outputDiscoveryTable(results []HostStatus, total, alive int) {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("NETWORK DISCOVERY REZULTATI")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Ukupno hostova: %d\n", total)
	fmt.Printf("Aktivnih hostova: %d\n", alive)
	fmt.Println(strings.Repeat("-", 50))

	if alive == 0 {
		fmt.Println("Nisu pronađeni aktivni hostovi.")
		return
	}

	fmt.Printf("%-20s %-10s %s\n", "IP ADRESA", "STATUS", "RTT")
	fmt.Println(strings.Repeat("-", 50))

	for _, r := range results {
		if r.Alive {
			fmt.Printf("%-20s %-10s %v\n", r.IP, "UP", r.RTT.Round(time.Millisecond))
		}
	}

	fmt.Println(strings.Repeat("=", 50))
}

func outputList(results []HostStatus) {
	for _, r := range results {
		if r.Alive {
			fmt.Println(r.IP)
		}
	}
}
