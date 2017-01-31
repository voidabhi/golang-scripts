
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	devName 	string
	es_index 	string
	es_docType 	string
	es_server	string
	err      	error
	handle   	*pcap.Handle
	InetAddr 	string
	SrcIP    	string
	DstIP    	string
)

type DnsMsg struct {
	Timestamp       string
	SourceIP        string
	DestinationIP   string
	DnsQuery        string
	DnsAnswer       []string
	DnsAnswerTTL    []string
	NumberOfAnswers string
	DnsResponseCode string
	DnsOpCode       string
}

func sendToElastic(dnsMsg DnsMsg, wg *sync.WaitGroup) {
	defer wg.Done()

	var jsonMsg, jsonErr = json.Marshal(dnsMsg)
	if jsonErr != nil {
		panic(jsonErr)
	}

	// getting ready for elasticsearch
	request, reqErr := http.NewRequest("POST", "http://"+es_server+":9200/"+es_index+"/"+es_docType,
		bytes.NewBuffer(jsonMsg))
	if reqErr != nil {
		panic(reqErr)
	}

	client := &http.Client{}
	resp, elErr := client.Do(request)

	if elErr != nil {
		panic(elErr)
	}

	defer resp.Body.Close()

}

func main() {

	////// CONFIG SECTION: REVIEW THESE BEFORE USING
	
	// select a device to listen on
	//windows example
	devName = "\\Device\\NPF_{9CA25EBF-B3D8-4FD0-90A6-070A16A7F2B4}"
	//linux example
	//devName = "eth0"

	// define an elasticsearch server to send to
	es_server = "192.168.10.15"
	// define an elasticsearch index to send to 
	es_index = "dns_index"
	es_docType = "syslog"

	// END CONFIG SECTION


	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	var dns layers.DNS

	var payload gopacket.Payload

	wg := new(sync.WaitGroup)


	// Find all devices
	devices, devErr := pcap.FindAllDevs()
	if devErr != nil {
		log.Fatal(devErr)
	}

	// Print device information
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			if device.Name == devName {
				InetAddr = address.IP.String()
				break
			}
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}

	// // Create DNSQuery index
	// _, elErr = client.CreateIndex("dns_query").Do()
	// if elErr != nil {
	//     // Handle error
	//     panic(elErr)
	// }

	// Open device
	handle, err = pcap.OpenLive(devName, 1600, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	var filter string = "udp and port 53 and src host " + InetAddr
	fmt.Println("    Filter: ", filter)
	err := handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &payload)

	decodedLayers := make([]gopacket.LayerType, 0, 10)
	for {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			fmt.Println("Error reading packet data: ", err)
			continue
		}

		err = parser.DecodeLayers(data, &decodedLayers)
		for _, typ := range decodedLayers {
			switch typ {
			case layers.LayerTypeIPv4:
				SrcIP = ip4.SrcIP.String()
				DstIP = ip4.DstIP.String()
			case layers.LayerTypeIPv6:
				SrcIP = ip6.SrcIP.String()
				DstIP = ip6.DstIP.String()
			case layers.LayerTypeDNS:
				dnsOpCode := int(dns.OpCode)
				dnsResponseCode := int(dns.ResponseCode)
				dnsANCount := int(dns.ANCount)

				if (dnsANCount == 0 && dnsResponseCode > 0) || (dnsANCount > 0) {

					fmt.Println("------------------------")
					fmt.Println("    DNS Record Detected")

					for _, dnsQuestion := range dns.Questions {

						t := time.Now()
						timestamp := t.Format(time.RFC3339)

						// Add a document to the index
						d := DnsMsg{Timestamp: timestamp, SourceIP: SrcIP,
							DestinationIP:   DstIP,
							DnsQuery:        string(dnsQuestion.Name),
							DnsOpCode:       strconv.Itoa(dnsOpCode),
							DnsResponseCode: strconv.Itoa(dnsResponseCode),
							NumberOfAnswers: strconv.Itoa(dnsANCount)}
						fmt.Println("    DNS OpCode: ", strconv.Itoa(int(dns.OpCode)))
						fmt.Println("    DNS ResponseCode: ", dns.ResponseCode.String())
						fmt.Println("    DNS # Answers: ", strconv.Itoa(dnsANCount))
						fmt.Println("    DNS Question: ", string(dnsQuestion.Name))
						fmt.Println("    DNS Endpoints: ", SrcIP, DstIP)

						if dnsANCount > 0 {

							for _, dnsAnswer := range dns.Answers {
								d.DnsAnswerTTL = append(d.DnsAnswerTTL, fmt.Sprint(dnsAnswer.TTL))
								if dnsAnswer.IP.String() != "<nil>" {
									fmt.Println("    DNS Answer: ", dnsAnswer.IP.String())
									d.DnsAnswer = append(d.DnsAnswer, dnsAnswer.IP.String())
								}
							}

						}

						wg.Add(1)
						sendToElastic(d, wg)

					}
				}

			}
		}

		if err != nil {
			fmt.Println("  Error encountered:", err)
		}
	}
}
