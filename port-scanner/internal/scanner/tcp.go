package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortResult predstavlja rezultat skeniranja jednog porta
type PortResult struct {
	Port     int
	State    string
	Service  string
	Banner   string
	ScanTime time.Duration
}

// TCPScanner je glavni skener
type TCPScanner struct {
	concurrency int
	timeout     time.Duration
	retries     int
	rateLimit   int
	rateLimiter <-chan time.Time
}

// NewTCPScanner kreira novi TCP skener
func NewTCPScanner(concurrency int, timeout time.Duration, retries int, rateLimit int) *TCPScanner {
	var rateLimiter <-chan time.Time
	if rateLimit > 0 {
		rateLimiter = time.Tick(time.Second / time.Duration(rateLimit))
	}

	return &TCPScanner{
		concurrency: concurrency,
		timeout:     timeout,
		retries:     retries,
		rateLimit:   rateLimit,
		rateLimiter: rateLimiter,
	}
}

// ScanPorts skenira listu portova na zadanom hostu
func (s *TCPScanner) ScanPorts(host string, ports []int, grabBanner bool) []PortResult {
	var results []PortResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Semaphore za kontrolu konkurentnosti
	semaphore := make(chan struct{}, s.concurrency)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			// Rate limiting
			if s.rateLimiter != nil {
				<-s.rateLimiter
			}

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := s.scanPort(host, p, grabBanner)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

func (s *TCPScanner) scanPort(host string, port int, grabBanner bool) PortResult {
	start := time.Now()
	result := PortResult{
		Port:    port,
		State:   "closed",
		Service: GetServiceName(port),
	}

	var conn net.Conn
	var err error

	// Retry logic
	for i := 0; i <= s.retries; i++ {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), s.timeout)
		if err == nil {
			break
		}
	}

	if err != nil {
		result.ScanTime = time.Since(start)
		return result
	}
	defer conn.Close()

	result.State = "open"

	// Banner grabbing
	if grabBanner {
		result.Banner = GrabBanner(conn, port, s.timeout)
	}

	result.ScanTime = time.Since(start)
	return result
}

// GetServiceName vraća ime servisa za poznate portove
func GetServiceName(port int) string {
	services := map[int]string{
		20:    "ftp-data",
		21:    "ftp",
		22:    "ssh",
		23:    "telnet",
		25:    "smtp",
		53:    "dns",
		67:    "dhcp",
		68:    "dhcp",
		69:    "tftp",
		80:    "http",
		110:   "pop3",
		111:   "rpcbind",
		119:   "nntp",
		123:   "ntp",
		135:   "msrpc",
		137:   "netbios-ns",
		138:   "netbios-dgm",
		139:   "netbios-ssn",
		143:   "imap",
		161:   "snmp",
		162:   "snmptrap",
		389:   "ldap",
		443:   "https",
		445:   "microsoft-ds",
		465:   "smtps",
		514:   "syslog",
		515:   "printer",
		587:   "submission",
		631:   "ipp",
		636:   "ldaps",
		873:   "rsync",
		993:   "imaps",
		995:   "pop3s",
		1080:  "socks",
		1433:  "mssql",
		1434:  "mssql-m",
		1521:  "oracle",
		1723:  "pptp",
		2049:  "nfs",
		2082:  "cpanel",
		2083:  "cpanel-ssl",
		2181:  "zookeeper",
		2375:  "docker",
		2376:  "docker-ssl",
		3000:  "grafana",
		3306:  "mysql",
		3389:  "rdp",
		3690:  "svn",
		4443:  "https-alt",
		5000:  "upnp",
		5432:  "postgresql",
		5672:  "amqp",
		5900:  "vnc",
		5984:  "couchdb",
		6379:  "redis",
		6443:  "kubernetes",
		6667:  "irc",
		7001:  "weblogic",
		8000:  "http-alt",
		8008:  "http-alt",
		8080:  "http-proxy",
		8443:  "https-alt",
		8888:  "http-alt",
		9000:  "php-fpm",
		9090:  "prometheus",
		9200:  "elasticsearch",
		9300:  "elasticsearch",
		9418:  "git",
		10000: "webmin",
		11211: "memcached",
		27017: "mongodb",
		27018: "mongodb",
		28015: "rethinkdb",
		50000: "db2",
	}

	if service, ok := services[port]; ok {
		return service
	}
	return "unknown"
}
