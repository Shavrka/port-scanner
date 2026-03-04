#!/usr/bin/env zsh
#
# ping_sweep.zsh - Network Discovery Script za Fedora/GNOME
# Koristi se za otkrivanje aktivnih hostova na mreži
#
# Korištenje:
#   ./ping_sweep.zsh <CIDR>
#   ./ping_sweep.zsh 192.168.1.0/24
#   ./ping_sweep.zsh 192.168.1.0/24 -q         # Quiet mode (samo IP adrese)
#   ./ping_sweep.zsh 192.168.1.0/24 -c 10      # Konkurentnost
#   ./ping_sweep.zsh 192.168.1.0/24 -o hosts.txt  # Output u fajl
#

set -euo pipefail

# Boje za terminal
autoload -U colors && colors
RED="%{$fg[red]%}"
GREEN="%{$fg[green]%}"
YELLOW="%{$fg[yellow]%}"
BLUE="%{$fg[blue]%}"
RESET="%{$reset_color%}"

# Defaults
CONCURRENCY=50
QUIET=false
OUTPUT_FILE=""

# Funkcija za pomoć
show_help() {
    cat << EOF
ping_sweep.zsh - Network Discovery Script

KORIŠTENJE:
    ./ping_sweep.zsh <CIDR> [opcije]

OPCIJE:
    -c, --concurrency <n>  Broj istovremenih pingova (default: 50)
    -q, --quiet            Tihi mod - samo IP adrese
    -o, --output <file>    Sačuvaj rezultate u fajl
    -h, --help             Prikaži ovu pomoć

PRIMJERI:
    ./ping_sweep.zsh 192.168.1.0/24
    ./ping_sweep.zsh 10.0.0.0/24 -c 100 -o active_hosts.txt
    ./ping_sweep.zsh 172.16.0.0/16 -q | xargs -I {} port-scanner scan -target {}

EOF
}

# Parsing argumenata
CIDR=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--concurrency)
            CONCURRENCY="$2"
            shift 2
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        -*)
            echo "Nepoznata opcija: $1"
            show_help
            exit 1
            ;;
        *)
            CIDR="$1"
            shift
            ;;
    esac
done

# Provjera da je CIDR naveden
if [[ -z "$CIDR" ]]; then
    echo "Greška: CIDR range je obavezan"
    show_help
    exit 1
fi

# Funkcija za generisanje IP adresa iz CIDR-a
generate_ips() {
    local cidr=$1
    local base_ip=${cidr%/*}
    local prefix=${cidr#*/}
    
    # Izračunaj broj hostova
    local num_hosts=$((2 ** (32 - prefix) - 2))
    
    # Parse base IP
    local IFS='.'
    read -r a b c d <<< "$base_ip"
    
    # Izračunaj network address
    local mask=$((0xFFFFFFFF << (32 - prefix)))
    local ip_int=$(( (a << 24) + (b << 16) + (c << 8) + d ))
    local network_int=$((ip_int & mask))
    
    # Generiši sve IP adrese (osim network i broadcast)
    for ((i = 1; i <= num_hosts; i++)); do
        local current_int=$((network_int + i))
        local o1=$(( (current_int >> 24) & 255 ))
        local o2=$(( (current_int >> 16) & 255 ))
        local o3=$(( (current_int >> 8) & 255 ))
        local o4=$(( current_int & 255 ))
        echo "${o1}.${o2}.${o3}.${o4}"
    done
}

# Funkcija za ping jednog hosta
ping_host() {
    local ip=$1
    if ping -c 1 -W 1 "$ip" &>/dev/null; then
        echo "$ip"
    fi
}

# Export funkcije za parallel execution
export -f ping_host 2>/dev/null || true

# Main
if [[ "$QUIET" == false ]]; then
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "  PING SWEEP - Network Discovery"
    echo "════════════════════════════════════════════════════════════"
    echo "  Target:        $CIDR"
    echo "  Concurrency:   $CONCURRENCY"
    echo "════════════════════════════════════════════════════════════"
    echo ""
    echo "[*] Generišem IP adrese..."
fi

# Generiši IP adrese
ip_list=$(generate_ips "$CIDR")
total_ips=$(echo "$ip_list" | wc -l)

if [[ "$QUIET" == false ]]; then
    echo "[*] Skeniram $total_ips hostova..."
    echo ""
fi

# Izvršavanje ping sweep-a
active_hosts=""
start_time=$(date +%s)

# Koristi xargs za paralelno izvršavanje
if command -v parallel &>/dev/null; then
    # GNU Parallel je dostupan
    active_hosts=$(echo "$ip_list" | parallel -j "$CONCURRENCY" --timeout 2 "ping -c 1 -W 1 {} &>/dev/null && echo {}" 2>/dev/null)
else
    # Fallback na xargs
    active_hosts=$(echo "$ip_list" | xargs -P "$CONCURRENCY" -I {} sh -c 'ping -c 1 -W 1 "$1" &>/dev/null && echo "$1"' _ {})
fi

end_time=$(date +%s)
duration=$((end_time - start_time))

# Broj aktivnih hostova
active_count=$(echo "$active_hosts" | grep -c . || echo 0)

# Output
if [[ "$QUIET" == true ]]; then
    echo "$active_hosts" | sort -t. -k1,1n -k2,2n -k3,3n -k4,4n
else
    echo "════════════════════════════════════════════════════════════"
    echo "  REZULTATI"
    echo "════════════════════════════════════════════════════════════"
    echo ""
    
    if [[ -n "$active_hosts" ]]; then
        echo "  AKTIVNI HOSTOVI ($active_count):"
        echo "  ────────────────────────────────────"
        echo "$active_hosts" | sort -t. -k1,1n -k2,2n -k3,3n -k4,4n | while read -r ip; do
            echo "  ✓ $ip"
        done
    else
        echo "  Nisu pronađeni aktivni hostovi."
    fi
    
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "  Skenirano:     $total_ips hostova"
    echo "  Aktivnih:      $active_count"
    echo "  Trajanje:      ${duration}s"
    echo "════════════════════════════════════════════════════════════"
fi

# Sačuvaj u fajl ako je navedeno
if [[ -n "$OUTPUT_FILE" ]]; then
    echo "$active_hosts" | sort -t. -k1,1n -k2,2n -k3,3n -k4,4n > "$OUTPUT_FILE"
    if [[ "$QUIET" == false ]]; then
        echo ""
        echo "[+] Rezultati sačuvani u: $OUTPUT_FILE"
    fi
fi

# Desktop notifikacija (GNOME)
if [[ "$QUIET" == false ]] && command -v notify-send &>/dev/null; then
    notify-send "Ping Sweep Završen" "Pronađeno $active_count aktivnih hostova na $CIDR" 2>/dev/null || true
fi

exit 0
