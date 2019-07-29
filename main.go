package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

const ipListUrl = "https://raw.githubusercontent.com/ejrv/VPNs/master/vpn-ipv4.txt"

var ips []*net.IPNet

func main() {
	//get ip list from github
	ips = getIps()

	//server
	http.HandleFunc("/check", handler)
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["ip"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintf(w,"Url Param 'ip' is missing")
		return
	}

	ip := keys[0]
	fmt.Fprintf(w, "%v", checkIp(net.ParseIP(ip)))
}

func getIps() []*net.IPNet {
	resp, err := http.Get(ipListUrl)
	if err != nil {
		log.Fatalln("Cannot get ip's to block, err:", err)
	}
	defer resp.Body.Close()

	ips := make([]*net.IPNet, 0, 35000)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		ip, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			ip = net.ParseIP(line)
			ipnet = &net.IPNet{ip, net.IPv4Mask(255, 255, 255, 255)}
		}
		if ip == nil  {
			log.Printf("Unable to parse line: \"%v\" as ip.\n", line)
			continue
		}

		ips = append(ips, ipnet)
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading response body:", err)
	}
	return ips
}

func checkIp(ip net.IP) bool {
	for _, ipnet := range ips {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}
