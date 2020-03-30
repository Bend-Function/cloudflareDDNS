//Author: BendFunction
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var ipArray []string

func updateDNSRecord(zonesID, domainID, email, apiKey, subDomain, IP string) string {
	//Payload is the payload
	type Payload struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
		Proxied bool   `json:"proxied"`
	}

	data := Payload{
		// fill struct
	}

	data.Type = "A"
	data.Name = subDomain
	data.Content = IP
	data.TTL = 120
	data.Proxied = false
	// fmt.Println(data)
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", "https://api.cloudflare.com/client/v4/zones/"+zonesID+"/dns_records/"+domainID+"/", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)
	res := string(buf.Bytes())
	state := strings.Contains(res, `"success":true`)
	if state == true {
		rtn := data.Name + "  updateDNSRecord---SUCCEESS!"
		return rtn
	} else {
		fmt.Println(res)
		return "updateDNSRecord---ERROR!"
	}
}

func getZonesID(email, apiKey, mainDomain string) string {
	data := ""
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones?name="+mainDomain, body)
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)
	res := string(buf.Bytes())
	state := strings.Contains(res, `"success":true`)
	// return res
	if state == true {
		startIndex := strings.Index(res, "id") + 5
		rtn := res[startIndex : startIndex+32]
		return string(rtn)
	} else {
		fmt.Println("getZoneID--ERROR!")
		return "getZoneID--ERROR!"
	}

}

func getDonmainID(email, apiKey, zoneID, subDomain string) string {
	data := ""
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records?type=A&name="+subDomain, body)
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)
	res := string(buf.Bytes())
	state := strings.Contains(res, `"success":true`)
	// return res
	if state == true {
		startIndex := strings.Index(res, "id") + 5
		rtn := res[startIndex : startIndex+32]
		return string(rtn)
	} else {
		fmt.Println("getDonmainID--ERROR!")
		return "getDonmainID--ERROR!"
	}
}

func search(key string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == key {
			return true
		}
	}
	return false
}

func getip(detectURL string) string {
	resp, err := http.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)

	ren := string(buf.Bytes())
	ipArray = append(ipArray, strings.Replace(ren, "\n", "", -1))
	return "0"
}

func getAllIPs(detectURL string) []string {
	for i := 0; i < 30; i++ {
		go getip(detectURL)
	}
	time.Sleep(time.Second * 9)
	var cutIP []string

	for i := 0; i < len(ipArray); i++ {
		if search(ipArray[i], cutIP) == false {
			cutIP = append(cutIP, ipArray[i])
		}
	}

	return cutIP
}

func main() {
	var subDomainArray []string
	// make config struct
	type conf struct {
		Email           string   `json:"email"`
		APIKey          string   `json:"apiKey"`
		MainDomain      string   `json:"mainDomain"`
		SubDomainArray  []string `json:"subDomainArray"`
		IPdetectAddress string   `json:"IPdetectAddress"`
	}
	// get arguments
	args := os.Args
	//set default config path
	confpath := "src/config/conf.json"
	if len(args) == 3 && args[1] == "-c" {
		confpath = args[2]
	} else if len(args) >= 2 {
		fmt.Println("Usage: [-c configPath] ")
		fmt.Println("-c: config file path")
		fmt.Println("repo url: https://github.com/Bend-Function/cloudflareDDNS")
		os.Exit(1)
	} else {
		fmt.Println("Program will use default config path")
		fmt.Println("src/config/conf.json")
	}
	// read config file
	file, err := os.Open(confpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := conf{}
	err = decoder.Decode(&config)
	if err != nil {
		// handle err
		fmt.Println(err)
	}
	// make old version value equle to new version
	email := config.Email
	apiKey := config.APIKey
	mainDomain := config.MainDomain
	subDomainArray = config.SubDomainArray

	// print domain that will be update
	fmt.Println("These domains will be update")
	for i := 0; i < len(subDomainArray); i++ {
		fmt.Println(subDomainArray[i] + "." + mainDomain)
	}

	// get zonesID
	zonesID := getZonesID(email, apiKey, mainDomain)

	// get local ip
	var IP []string
	// get no more than 3 times
	for j := 0; j < 3; j++ {
		IP = getAllIPs(config.IPdetectAddress)
		// if get ips == domains break
		if len(IP) == len(subDomainArray) {
			break
		}
	}
	if len(IP) == len(subDomainArray) {
		fmt.Println("Number of ip is equal to number of domains")
		for j := 0; j < len(IP); j++ {
			fmt.Println(IP[j] + "--->" + subDomainArray[j] + "." + mainDomain)
		}
	} else {
		fmt.Println("Number of ip ISN'T equal to number of domains!!")
		var tempNum int
		if len(subDomainArray) < len(IP) {
			tempNum = len(subDomainArray)
		} else {
			tempNum = len(IP)
		}
		for j := 0; j < tempNum; j++ {
			fmt.Println(IP[j] + "--->" + subDomainArray[j] + "." + mainDomain)
		}
	}

	// updateDNSRecord
	for i := 0; i < len(subDomainArray); i++ {
		domainID := getDonmainID(email, apiKey, zonesID, subDomainArray[i]+"."+mainDomain)
		fmt.Println(updateDNSRecord(zonesID, domainID, email, apiKey, subDomainArray[i], IP[i]))
	}
}
