package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ParseTarget parsira target string i vraća listu IP adresa
// Podržava: pojedinačnu IP, CIDR notaciju, ili range (192.168.1.1-10)
func ParseTarget(target string) ([]string, error) {
	// Provjeri da li je CIDR
	if strings.Contains(target, "/") {
		return parseCIDR(target)
	}

	// Provjeri da li je range
	if strings.Contains(target, "-") {
		return parseRange(target)
	}

	// Pojedinačna IP adresa
	ip := net.ParseIP(target)
	if ip == nil {
		// Možda je hostname
		ips, err := net.LookupIP(target)
		if err != nil {
			return nil, fmt.Errorf("nevažeća IP adresa ili hostname: %s", target)
		}
		var result []string
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				result = append(result, ipv4.String())
			}
		}
		return result, nil
	}

	return []string{target}, nil
}

// parseCIDR parsira CIDR notaciju i vraća sve IP adrese
func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("nevažeći CIDR: %s", cidr)
	}

	var ips []string

	// Iteriraj kroz sve IP adrese u rangeu
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	// Ukloni network i broadcast adresu za /24 i manje
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}

	return ips, nil
}

// incrementIP inkrementira IP adresu za 1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// parseRange parsira IP range (npr. 192.168.1.1-10)
func parseRange(target string) ([]string, error) {
	parts := strings.Split(target, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("nevažeći range format: %s", target)
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	if startIP == nil {
		return nil, fmt.Errorf("nevažeća početna IP: %s", parts[0])
	}

	// Provjeri da li je kraj samo broj ili puna IP
	endPart := strings.TrimSpace(parts[1])
	var endIP net.IP

	if strings.Contains(endPart, ".") {
		// Puna IP adresa
		endIP = net.ParseIP(endPart)
	} else {
		// Samo zadnji oktet
		lastOctet, err := strconv.Atoi(endPart)
		if err != nil || lastOctet < 0 || lastOctet > 255 {
			return nil, fmt.Errorf("nevažeći zadnji oktet: %s", endPart)
		}
		// Kopiraj početnu IP i izmijeni zadnji oktet
		endIP = make(net.IP, len(startIP))
		copy(endIP, startIP)
		endIP = endIP.To4()
		endIP[3] = byte(lastOctet)
	}

	if endIP == nil {
		return nil, fmt.Errorf("nevažeća krajnja IP: %s", endPart)
	}

	var ips []string
	startIP = startIP.To4()
	endIP = endIP.To4()

	for ip := startIP; compareIPs(ip, endIP) <= 0; incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

// compareIPs uporedi dvije IP adrese
func compareIPs(ip1, ip2 net.IP) int {
	for i := 0; i < 4; i++ {
		if ip1[i] < ip2[i] {
			return -1
		}
		if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

// ParsePorts parsira port string i vraća listu portova
// Podržava: pojedinačni port, range (1-1000), ili listu (22,80,443)
func ParsePorts(portStr string) ([]int, error) {
	var ports []int
	seen := make(map[int]bool)

	// Split by comma
	parts := strings.Split(portStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			// Range
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("nevažeći port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("nevažeći početni port: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("nevažeći krajnji port: %s", rangeParts[1])
			}

			if start < 1 || start > 65535 || end < 1 || end > 65535 {
				return nil, fmt.Errorf("port mora biti između 1 i 65535")
			}

			if start > end {
				start, end = end, start
			}

			for p := start; p <= end; p++ {
				if !seen[p] {
					ports = append(ports, p)
					seen[p] = true
				}
			}
		} else {
			// Single port
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("nevažeći port: %s", part)
			}

			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port mora biti između 1 i 65535")
			}

			if !seen[port] {
				ports = append(ports, port)
				seen[port] = true
			}
		}
	}

	return ports, nil
}
