# Port Scanner

**Brzi concurrent port scanner napisan u Go programskom jeziku**

Port Scanner je alat za mreznu sigurnost koji omogucava skeniranje TCP portova na pojedinacnim hostovima ili cijelim mreznim rangovima. Program koristi Go-ove goroutine za paralelno skeniranje stotina portova istovremeno, sto ga cini izuzetno brzim u odnosu na sekvencijalne skenere.

---

## Sadrzaj

- [Svrha programa](#svrha-programa)
- [Karakteristike](#karakteristike)
- [Sistemski zahtjevi](#sistemski-zahtjevi)
- [Instalacija](#instalacija)
- [Azuriranje](#azuriranje-update)
- [Koriscenje](#koriscenje)
- [Komande i opcije](#komande-i-opcije)
- [Primjeri](#primjeri)
- [Struktura projekta](#struktura-projekta)
- [Tehnicka implementacija](#tehnicka-implementacija)
- [Sigurnosne napomene](#sigurnosne-napomene)

---

## Svrha programa

Port Scanner je namijenjen za:

1. **Network Security Auditing** - Provjera otvorenih portova na serverima i mreznim uredjajima
2. **Penetration Testing** - Identifikacija potencijalnih ulaznih tacaka u mrezu
3. **System Administration** - Verifikacija da su samo potrebni portovi otvoreni
4. **Service Discovery** - Otkrivanje koji servisi rade na mreznoj infrastrukturi
5. **Network Inventory** - Mapiranje aktivnih hostova i servisa na mrezi

Program je posebno koristan za:
- DevOps inzenjere koji trebaju provjeriti firewall konfiguracije
- Sigurnosne analiticare koji vrse penetracijske testove
- Sistem administratore koji odrzavaju mreznu infrastrukturu

---

## Karakteristike

### Osnovne mogucnosti

- **TCP Connect Scan** - Potpuna TCP konekcija za pouzdanu detekciju otvorenih portova
- **CIDR Range Scanning** - Skeniranje citavih podmreza (npr. 192.168.1.0/24)
- **Banner Grabbing** - Automatska identifikacija servisa i njihovih verzija
- **Service Detection** - Prepoznavanje 50+ poznatih servisa po portu

### Performanse

- **Goroutine Pool** - Semaphore pattern za kontrolisanu konkurentnost
- **Konfigurisana konkurentnost** - Od 1 do 1000+ istovremenih konekcija
- **Rate Limiting** - Kontrola brzine skeniranja za izbjegavanje IDS detekcije
- **Retry Logic** - Automatsko ponavljanje neuspjesnih konekcija

### Izlazni formati

- **Table** - Pregledna tabelarna forma za terminal
- **JSON** - Strukturirani format za integraciju sa drugim alatima

### Korisnicko iskustvo

- **Interaktivni mod** - Meni-baziran interfejs za lakse koriscenje
- **CLI mod** - Direktne komande za skriptiranje i automatizaciju
- **Quick Scan** - Brzo skeniranje najcescih portova

---

## Sistemski zahtjevi

### Obavezni

| Komponenta  | Verzija         |
|-------------|-----------------|
| Go          | 1.21 ili noviji |
| Linux       | Kernel 4.x+     |
| Arhitektura | x86_64, ARM64   |

### Testirano na

- Fedora 39, 40, 41
- Ubuntu 22.04, 24.04
- Debian 12
- Arch Linux

### Opcionalno

- `ping` komanda (za network discovery)
- Root privilegije (za neke tipove skeniranja)

---

## Instalacija

### Metoda 1: Build iz izvornog koda

```bash
# 1. Kloniraj ili kopiraj projekat
cd /putanja/do/port-scanner

# 2. Inicijaliziraj Go modul (ako vec nije)
go mod init port-scanner

# 3. Kompajliraj program
go build -o port-scanner .

# 4. Testiraj instalaciju
./port-scanner help
```

### Metoda 2: Sistemska instalacija (preporuceno)

```bash
# 1. Kompajliraj program
go build -o port-scanner .

# 2. Kopiraj u sistemski PATH
sudo cp port-scanner /usr/local/bin/

# 3. Postavi izvrsne dozvole
sudo chmod +x /usr/local/bin/port-scanner

# 4. Verificiraj instalaciju
port-scanner help
```

### Metoda 3: Instalacija ping-sweep skripte

```bash
# Kopiraj Zsh skriptu za network discovery
sudo cp scripts/ping_sweep.zsh /usr/local/bin/ping-sweep
sudo chmod +x /usr/local/bin/ping-sweep
```

### Metoda 4: Koriscenje Makefile (preporuceno)

```bash
# Build programa
make build

# Instalacija na sistem
sudo make install

# Provjeri verziju
port-scanner version
```

### Azuriranje (Update)

Kada izadje nova verzija, azuriraj program na sljedeci nacin:

#### Automatsko azuriranje (preporuceno)

```bash
# Pozicioniraj se u direktorij projekta
cd /putanja/do/port-scanner

# Povuci najnovije izmjene i reinstaliraj
make update
```

`make update` automatski izvrsava:
1. `git pull` - povlaci posljednje izmjene iz repozitorija
2. `make build` - kompajlira novu verziju
3. `sudo make install` - instalira na sistem

#### Rucno azuriranje

```bash
# 1. Pozicioniraj se u direktorij projekta
cd /putanja/do/port-scanner

# 2. Povuci najnovije izmjene
git pull origin main

# 3. Kompajliraj novu verziju sa verzionisanjem
make build

# 4. Instaliraj na sistem
sudo make install

# 5. Provjeri novu verziju
port-scanner version
```

#### Azuriranje bez Makefile

```bash
cd /putanja/do/port-scanner
git pull origin main
go build -o port-scanner .
sudo cp port-scanner /usr/local/bin/
port-scanner version
```

### Deinstalacija

```bash
# Koriscenjem Makefile
sudo make uninstall

# Ili rucno
sudo rm /usr/local/bin/port-scanner
sudo rm /usr/local/bin/ping-sweep  # ako je instalirana
```

---

## Koriscenje

### Interaktivni mod

Pokreni program bez argumenata za interaktivni meni:

```bash
port-scanner
```

Prikazat ce se meni:

```
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
└──────────────────────────────────────────────────────────────────────────────┘

> Select option:
```

Izaberi opciju unosom broja (0-4) i pritisni ENTER.

### CLI mod

Za skriptiranje i automatizaciju koristi direktne komande:

```bash
port-scanner <komanda> [opcije]
```

Dostupne komande:
- `scan` - Skeniranje portova
- `discover` - Network discovery (ping sweep)
- `help` - Prikaz pomoci

---

## Komande i opcije

### Komanda: scan

Skenira TCP portove na zadatoj adresi ili mreznom rangu.

```bash
port-scanner scan [opcije]
```

#### Opcije

| Opcija | Opis | Podrazumijevano | Primjer |
|--------|------|-----------------|---------|
| `-target` | Ciljna IP adresa, hostname ili CIDR | **obavezno** | `192.168.1.1`, `google.com`, `10.0.0.0/24` |
| `-ports` | Portovi za skeniranje | `1-1024` | `80,443`, `1-65535`, `22,80,443,3306` |
| `-concurrency` | Broj istovremenih konekcija | `100` | `50`, `200`, `500` |
| `-timeout` | Timeout po konekciji | `2s` | `1s`, `5s`, `500ms` |
| `-rate-limit` | Max konekcija po sekundi (0=bez limita) | `0` | `100`, `500` |
| `-output` | Format izlaza | `table` | `table`, `json` |
| `-banner` | Ukljuci banner grabbing | `true` | `true`, `false` |
| `-retries` | Broj ponovnih pokusaja | `1` | `2`, `3` |

### Komanda: discover

Otkriva aktivne hostove na mrezi koristeci ping sweep.

```bash
port-scanner discover [opcije]
```

#### Opcije

| Opcija | Opis | Podrazumijevano | Primjer |
|--------|------|-----------------|---------|
| `-target` | CIDR range za skeniranje | **obavezno** | `192.168.1.0/24`, `10.0.0.0/16` |
| `-concurrency` | Broj istovremenih ping zahtjeva | `50` | `20`, `100` |
| `-timeout` | Timeout za ping | `2s` | `1s`, `3s` |

### Komanda: help

Prikazuje pomoc i uputstva za koriscenje.

```bash
port-scanner help
```

---

## Primjeri

### Osnovno skeniranje

```bash
# Skeniranje jednog hosta, portovi 1-1024
port-scanner scan -target 192.168.1.1

# Skeniranje specificnih portova
port-scanner scan -target 192.168.1.1 -ports 22,80,443,3306,5432

# Skeniranje opsega portova
port-scanner scan -target 192.168.1.1 -ports 1-1000
```

### Skeniranje mreze (CIDR)

```bash
# Skeniranje citave /24 podmreze
port-scanner scan -target 192.168.1.0/24 -ports 22,80,443

# Skeniranje veceg ranga
port-scanner scan -target 10.0.0.0/16 -ports 80,443 -concurrency 200
```

### Napredne opcije

```bash
# Brzo skeniranje sa vecim brojem konekcija
port-scanner scan -target 192.168.1.1 -ports 1-65535 -concurrency 500 -timeout 1s

# Tiho skeniranje sa rate limitingom (izbjegavanje IDS)
port-scanner scan -target 192.168.1.1 -ports 1-1000 -rate-limit 50 -concurrency 10

# JSON izlaz za dalju obradu
port-scanner scan -target 192.168.1.1 -ports 80,443 -output json > rezultati.json

# Skeniranje bez banner grabbinga (brze)
port-scanner scan -target 192.168.1.1 -ports 1-1000 -banner=false
```

### Network Discovery

```bash
# Otkrivanje aktivnih hostova na mrezi
port-scanner discover -target 192.168.1.0/24

# Sa vecim timeout-om za spore mreze
port-scanner discover -target 10.0.0.0/24 -timeout 3s -concurrency 30
```

### Kombinovani workflow

```bash
# 1. Prvo otkrij aktivne hostove
port-scanner discover -target 192.168.1.0/24

# 2. Zatim skeniraj pronadjene hostove
port-scanner scan -target 192.168.1.1 -ports 1-1000
port-scanner scan -target 192.168.1.10 -ports 1-1000
```

---

## Primjer izlaza

### Tabelarni format

```
================================================================================
REZULTATI SKENIRANJA
================================================================================
Vrijeme pocetka: 2026-03-04 15:30:00
Vrijeme zavrsetka: 2026-03-04 15:30:05
Trajanje: 5.123s
Skenirano hostova: 1
Skenirano portova: 1024
Otvorenih portova: 4
--------------------------------------------------------------------------------
HOST                 PORT     STATE      SERVICE         BANNER
--------------------------------------------------------------------------------
192.168.1.1          22       open       ssh             OpenSSH_8.9p1
192.168.1.1          80       open       http            nginx/1.24.0
192.168.1.1          443      open       https           
192.168.1.1          3306     open       mysql           MySQL 8.0.35
================================================================================
```

### JSON format

```json
{
  "start_time": "2026-03-04T15:30:00Z",
  "end_time": "2026-03-04T15:30:05Z",
  "total_hosts": 1,
  "total_ports": 1024,
  "open_ports": 4,
  "results": [
    {
      "host": "192.168.1.1",
      "port": 22,
      "state": "open",
      "service": "ssh",
      "banner": "OpenSSH_8.9p1"
    },
    {
      "host": "192.168.1.1",
      "port": 80,
      "state": "open",
      "service": "http",
      "banner": "nginx/1.24.0"
    }
  ]
}
```

---

## Struktura projekta

```
port-scanner/
├── main.go                      # Ulazna tacka, interaktivni meni
├── go.mod                       # Go module definicija
├── cmd/
│   ├── scan.go                  # Implementacija scan komande
│   └── discover.go              # Implementacija discover komande
├── internal/
│   ├── scanner/
│   │   ├── tcp.go               # TCP skener sa goroutine poolom
│   │   └── banner.go            # Banner grabbing logika
│   └── network/
│       └── cidr.go              # CIDR parsiranje i IP iteracija
├── scripts/
│   └── ping_sweep.zsh           # Zsh skripta za ping sweep
└── README.md                    # Dokumentacija
```

### Opis komponenti

| Datoteka | Opis |
|----------|------|
| `main.go` | Glavni ulaz, CLI parsing, interaktivni meni |
| `cmd/scan.go` | Logika za port skeniranje, output formatiranje |
| `cmd/discover.go` | Network discovery implementacija |
| `internal/scanner/tcp.go` | TCP konekcije, goroutine pool, rate limiting |
| `internal/scanner/banner.go` | Protokol-specifican banner grabbing |
| `internal/network/cidr.go` | CIDR parsing, IP generacija |

---

## Tehnicka implementacija

### Goroutine Pool Pattern

Program koristi semaphore pattern za kontrolu konkurentnosti:

```go
semaphore := make(chan struct{}, concurrency)

for _, port := range ports {
    go func(p int) {
        semaphore <- struct{}{}        // Zauzmi slot
        defer func() { <-semaphore }() // Oslobodi slot
        
        // Skeniranje porta
    }(port)
}
```

### Podrzani servisi za banner grabbing

| Port | Servis | Metoda |
|------|--------|--------|
| 21 | FTP | Citanje welcome bannera |
| 22 | SSH | Citanje verzije protokola |
| 25, 587 | SMTP | Citanje SMTP bannera |
| 80, 8080 | HTTP | HEAD zahtjev, Server header |
| 110 | POP3 | Citanje +OK bannera |
| 143 | IMAP | Citanje capability bannera |
| 443 | HTTPS | Detekcija SSL/TLS |
| 3306 | MySQL | Citanje handshake paketa |
| 5432 | PostgreSQL | Detekcija |
| 6379 | Redis | INFO server komanda |

### Mapa poznatih portova

Program prepoznaje 50+ standardnih servisa ukljucujuci:
- Web servise (HTTP, HTTPS, HTTP-Proxy)
- Baze podataka (MySQL, PostgreSQL, MongoDB, Redis)
- Mail servise (SMTP, POP3, IMAP)
- Remote access (SSH, RDP, VNC, Telnet)
- File sharing (FTP, SMB, NFS)
- I mnoge druge...

---

## Sigurnosne napomene

### UPOZORENJE

**Ovaj alat je namijenjen iskljucivo za legitimno testiranje sigurnosti mreza za koje imate eksplicitnu dozvolu.**

### Legalni aspekti

1. **Uvijek dobijte pisanu dozvolu** prije skeniranja bilo koje mreze
2. **Neovlasceno skeniranje** moze biti ilegalno u mnogim jurisdikcijama
3. **Koristite samo na vlastitoj infrastrukturi** ili uz pisanu autorizaciju

### Preporuke za koriscenje

1. **Rate limiting** - Koristite `-rate-limit` za izbjegavanje preopterecenja mreze
2. **Konkurentnost** - Pocnite sa manjim brojem (`-concurrency 50`) pa povecavajte
3. **Timeout** - Prilagodite timeout sprorosti mreze
4. **Logging** - Cuvajte JSON izlaz kao dokaz obavljenog testiranja

### IDS/IPS izbjegavanje

Za diskretnije skeniranje:

```bash
# Sporo skeniranje koje je teze detektovati
port-scanner scan -target 192.168.1.1 -ports 22,80,443 -rate-limit 10 -concurrency 5 -timeout 5s
```

---

## Cestia pitanja (FAQ)

### Zasto program ne pronalazi otvorene portove?

1. **Firewall** - Host moze imati firewall koji blokira konekcije
2. **Timeout** - Povecajte timeout sa `-timeout 5s`
3. **Mrezna latencija** - Smanjite konkurentnost za stabilnije rezultate

### Kako skenirati sve portove?

```bash
port-scanner scan -target 192.168.1.1 -ports 1-65535 -concurrency 500
```

### Mogu li sacuvati rezultate?

```bash
# JSON format u datoteku
port-scanner scan -target 192.168.1.1 -output json > scan_results.json
```

### Trebam li root privilegije?

Za vecinu operacija **nije potreban root**. Root je potreban samo ako zelite koristiti ICMP ping u discover modu.

---

## Licenca

MIT License

---

## Autor

Projekat kreiran kao portfolio projekat za demonstraciju:
- Go network programiranja (`net` paket)
- Goroutine pooling i concurrency patterns
- Network security koncepata
- CLI aplikacija sa interaktivnim interfejsom
