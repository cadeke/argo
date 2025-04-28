package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func ip2domain(baseUrl string) {
	ips := []string{
		// Existing records
		"8.8.8.8",
		"1.1.1.1",
		"9.9.9.9",
		"208.67.222.222",
		"192.168.1.1",
		"93.184.216.34",
		"142.250.190.46",
		"142.251.33.229",
		"172.217.10.238",
		"172.217.9.206",
		"104.16.132.229",
		"162.159.135.232",
		"140.82.114.3",
		"151.101.129.69",
		"98.138.219.231",
		"204.79.197.200",
		"35.186.224.25",
		"192.168.0.1",
		"203.0.113.1",
		"198.51.100.1",
		"192.0.2.1",
		"192.88.99.1",
		"203.0.113.99",
		"198.51.100.99",
		"208.80.154.224",
		"176.32.103.205",
		"52.250.42.157",
		"52.89.124.206",
		"23.246.0.5",
		"17.253.144.10",
		"104.215.148.63",
		"192.147.130.1",

		// Non-existing records
		"10.9.8.1",
		"10.9.8.2",
		"10.9.8.3",
		"10.9.8.4",
		"10.9.8.5",
		"44.44.44.44",
		"11.11.11.11",
		"99.99.99.99",
		"66.66.66.66",
		"77.77.77.77",
		"81.12.34.56",
		"50.50.50.50",
		"20.20.20.20",
		"78.19.92.3",
		"10.0.0.111",
		"10.0.0.222",
		"10.0.0.252",
		"127.123.45.67",
		"192.168.2.222",
	}

	ip := ips[rand.Intn(len(ips))]
	resp, err := http.Get(fmt.Sprintf("%s/api/ip2domain?query=%s", baseUrl, ip))
	if err != nil {
		fmt.Println("Error calling API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("API Response for IP", ip, ":", string(body))
}

func domain2ip(baseUrl string) {
	domains := []string{
		// Real domains
		"google.com", "facebook.com", "twitter.com", "github.com", "example.com",
		"wikipedia.org", "amazon.com", "netflix.com", "apple.com", "microsoft.com",
		"youtube.com", "linkedin.com", "instagram.com", "reddit.com", "yahoo.com",
		"bing.com", "ebay.com", "paypal.com", "stackoverflow.com", "quora.com",
		"cloudflare.com", "opendns.com", "quad9.net", "customdns.net",
		"example.net", "discord.com", "medium.com", "notion.so", "maps.google.com",
		"news.google.com", "drive.google.com", "play.google.com", "meet.google.com",
		"photos.google.com", "pinterest.com", "tiktok.com", "spotify.com",
		"localnetwork.com", "mywebsite.com", "testdomain.com", "exampledomain.net",
		"fictionalsite.com", "randompage.net", "privatehost.net", "experimental.com",
		"unusedsite.com", "espn.com", "cnn.com", "bbc.com",

		// Non-existent domains
		"nonexistentdomain.xyz", "randomsite.org", "fakeurl.io", "thisisnotreal.com",
		"totallylegitdomain.com", "anotherfakedomain.net", "notarealwebsite.org",
		"fakesite123.com", "doesnotexist.io", "imaginarydomain.net",
		"unrealwebsite.org", "madeupdomain.xyz", "notrealwebpage.com",
		"ghostsite.net", "phantomurl.org",
	}

	domain := domains[rand.Intn(len(domains))]
	resp, err := http.Get(fmt.Sprintf("%s/api/domain2ip?query=%s", baseUrl, domain))
	if err != nil {
		log.Fatal("Error calling API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return
	}

	fmt.Println("API Response for Domain", domain, ":", string(body))
}

func main() {
	host := os.Getenv("API_HOST")
	port := os.Getenv("API_PORT")

	// Consul service mesh
	// host = "127.0.0.1"
	// port = "8080"

	baseUrl := fmt.Sprintf("http://%s:%s", host, port)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ip2domain(baseUrl)
		domain2ip(baseUrl)
	}
}
