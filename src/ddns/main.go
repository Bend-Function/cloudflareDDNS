package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"os"
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
	fmt.Println(data)
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
		return "updateDNSRecord---SUCCEESS!"
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

func getip() string {
	resp, err := http.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.ReadFrom(resp.Body)

	ren := string(buf.Bytes())
	ipArray = append(ipArray, ren)
	return "0"
}

func getAllIPs() []string {
	for i := 0; i < 60; i++ {
		go getip()
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
	type conf struct {
		Email           string `json:"email"`
		APIKey          string `json:"apiKey"`
		MainDomain      string `json:"mainDomain"`
		SubDomainArray  []string `json:"subDomainArray"`
		IPdetectAddress string `json:"IPdetectAddress"`
	}

	file, _ := os.Open("src/config/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := conf{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	email := config.Email
	apiKey := config.APIKey
	mainDomain := config.MainDomain
	subDomainArray = config.SubDomainArray

	fmt.Print(subDomainArray)
	// IP := "167.123.123.123"
	zonesID := getZonesID(email, apiKey, mainDomain)

	// fmt.Print(domainID)
	var IP []string
	for j := 0; j < 3 ; j ++{
		IP = getAllIPs()
		fmt.Println(IP)
		if len(IP) == len(subDomainArray) {
			break
		}
	}

	for i := 0; i < len(subDomainArray); i++ {
		domainID := getDonmainID(email, apiKey, zonesID, subDomainArray[i]+"."+mainDomain)
		fmt.Println(updateDNSRecord(zonesID, domainID, email, apiKey, subDomainArray[i], strings.Replace(IP[i], "\n", "", -1)))
	}
}