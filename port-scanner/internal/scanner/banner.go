package scanner

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// GrabBanner pokušava izvući banner sa otvorenog porta
func GrabBanner(conn net.Conn, port int, timeout time.Duration) string {
	conn.SetReadDeadline(time.Now().Add(timeout))
	conn.SetWriteDeadline(time.Now().Add(timeout))

	switch port {
	case 80, 8080, 8000, 8888, 3000:
		return grabHTTPBanner(conn)
	case 443, 8443:
		return "HTTPS (SSL/TLS)"
	case 21:
		return grabFTPBanner(conn)
	case 22:
		return grabSSHBanner(conn)
	case 25, 587:
		return grabSMTPBanner(conn)
	case 110:
		return grabPOP3Banner(conn)
	case 143:
		return grabIMAPBanner(conn)
	case 3306:
		return grabMySQLBanner(conn)
	case 6379:
		return grabRedisBanner(conn)
	case 27017:
		return "MongoDB"
	default:
		return grabGenericBanner(conn)
	}
}

func grabHTTPBanner(conn net.Conn) string {
	// Pošalji HTTP HEAD zahtjev
	request := "HEAD / HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n"
	conn.Write([]byte(request))

	scanner := bufio.NewScanner(conn)
	var lines []string
	lineCount := 0

	for scanner.Scan() && lineCount < 5 {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
		lineCount++
	}

	if len(lines) > 0 {
		// Traži Server header
		for _, line := range lines {
			if strings.HasPrefix(strings.ToLower(line), "server:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Server:"))
			}
		}
		return lines[0]
	}
	return ""
}

func grabFTPBanner(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		banner := scanner.Text()
		// Remove FTP status code
		if len(banner) > 4 && banner[3] == ' ' {
			return strings.TrimSpace(banner[4:])
		}
		return banner
	}
	return ""
}

func grabSSHBanner(conn net.Conn) string {
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}
	banner := strings.TrimSpace(string(buffer[:n]))
	// SSH banner format: SSH-2.0-OpenSSH_8.9
	return banner
}

func grabSMTPBanner(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		banner := scanner.Text()
		// Remove SMTP status code
		if len(banner) > 4 && banner[3] == ' ' {
			return strings.TrimSpace(banner[4:])
		}
		return banner
	}
	return ""
}

func grabPOP3Banner(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		banner := scanner.Text()
		if strings.HasPrefix(banner, "+OK") {
			return strings.TrimSpace(strings.TrimPrefix(banner, "+OK"))
		}
		return banner
	}
	return ""
}

func grabIMAPBanner(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func grabMySQLBanner(conn net.Conn) string {
	buffer := make([]byte, 256)
	n, err := conn.Read(buffer)
	if err != nil || n < 5 {
		return "MySQL"
	}

	// MySQL handshake packet sadrži verziju
	// Packet format: [length:3][seq:1][protocol:1][version:null-terminated]
	if n > 5 {
		// Pronađi verziju string (null-terminated)
		start := 5
		end := start
		for end < n && buffer[end] != 0 {
			end++
		}
		if end > start {
			return fmt.Sprintf("MySQL %s", string(buffer[start:end]))
		}
	}
	return "MySQL"
}

func grabRedisBanner(conn net.Conn) string {
	// Pošalji INFO command
	conn.Write([]byte("INFO server\r\n"))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "Redis"
	}

	response := string(buffer[:n])
	lines := strings.Split(response, "\r\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "redis_version:") {
			return fmt.Sprintf("Redis %s", strings.TrimPrefix(line, "redis_version:"))
		}
	}
	return "Redis"
}

func grabGenericBanner(conn net.Conn) string {
	// Prvo pokušaj čitati (neki servisi šalju banner automatski)
	buffer := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	n, err := conn.Read(buffer)
	if err == nil && n > 0 {
		banner := strings.TrimSpace(string(buffer[:n]))
		// Ukloni neprintabilne karaktere
		banner = cleanBanner(banner)
		if len(banner) > 0 {
			return banner
		}
	}

	return ""
}

func cleanBanner(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r >= 32 && r < 127 || r == '\n' || r == '\r' || r == '\t' {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}
