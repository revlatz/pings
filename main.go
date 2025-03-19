package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// clearConsole löscht den Konsoleninhalt
func clearConsole() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// pingHost führt einen Ping zu einer IP aus und gibt true zurück, wenn sie online ist
func pingHost(ip string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}
	return cmd.Run() == nil
}

// loadIPsFromFile lädt die IP-Adressen aus einer Datei
func loadIPsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Fehler: Datei %s nicht gefunden!", filename)
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			ips = append(ips, line)
		}
	}
	return ips, scanner.Err()
}

// printHelp gibt die Hilfemeldung aus
func printHelp() {
	fmt.Println("Verwendung: pings [IP1 IP2 ...] oder pings [Datei] [Intervall in Sekunden]")
	fmt.Println("Beispiel:")
	fmt.Println("  pings 192.168.1.1 192.168.1.2 10")
	fmt.Println("  pings ip-list.txt 5")
	fmt.Println("Falls kein Intervall angegeben wird, beträgt der Standardwert 5 Sekunden.")
}

func main() {
	var ips []string
	interval := 5 // Standardintervall in Sekunden
	pingCounts := make(map[string]int)  // Erfolgreiche Pings
	failedCounts := make(map[string]int) // Fehlgeschlagene Pings

	// Prüfen, ob Parameter übergeben wurden
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-h" || arg == "--help" {
			printHelp()
			return
		}

		// Prüfen, ob das Argument eine Datei ist
		if _, err := os.Stat(arg); err == nil {
			// Falls Datei existiert, lade IPs aus Datei
			ips, err = loadIPsFromFile(arg)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Falls noch ein dritter Parameter (Intervall) vorhanden ist, interpretieren wir ihn als Zeit
			if len(os.Args) > 2 {
				if i, err := strconv.Atoi(os.Args[2]); err == nil {
					interval = i
				}
			}
		} else {
			// Falls keine Datei, dann als einzelne IPs behandeln
			ips = os.Args[1:]
			// Prüfen, ob das letzte Argument eine Zahl ist, um es als Intervall zu verwenden
			if len(ips) > 1 {
				if i, err := strconv.Atoi(ips[len(ips)-1]); err == nil {
					interval = i
					ips = ips[:len(ips)-1] // Entferne das Intervall aus der IP-Liste
				}
			}
		}
	} else {
		fmt.Println("Fehler: Bitte eine IP-Liste oder eine Datei als Parameter angeben!")
		printHelp()
		return
	}

	for {
		var output strings.Builder
		output.WriteString("\n--- IP Status Check ---\n")
		for _, ip := range ips {
			status := "Offline"
			if pingHost(ip) {
				pingCounts[ip]++ // Erfolgreicher Ping erhöhen
				status = "Online"
			} else {
				failedCounts[ip]++ // Fehlgeschlagener Ping erhöhen
			}

			// Farbige Statusausgabe
			statusColor := "\033[91m" // Rot für Offline
			if status == "Online" {
				statusColor = "\033[92m" // Grün für Online
			}
			counterFormat := fmt.Sprintf("(%d ok / %d !ok)", pingCounts[ip], failedCounts[ip])

			// Falls !ok > 0 ist, den Counter gelb markieren
			if failedCounts[ip] > 0 {
				counterFormat = fmt.Sprintf("\033[93m%s\033[0m", counterFormat)
			}

			// Einheitliche Formatierung mit fester Spaltenbreite
			output.WriteString(fmt.Sprintf("%-18s  %s%-8s\033[0m  %s\n", ip, statusColor, status, counterFormat))
		}

		clearConsole()
		fmt.Print(output.String())
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
