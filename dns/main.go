package main

import (
	"build-your-own/dns/data"
	"log"
	"net"

	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
)

func main() {
	// Initialize server address
	address := net.UDPAddr{
		Port: 8090,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// Start the UDP server using the server address
	log.Print("Starting DNS server at [:8090]")
	udp, err := net.ListenUDP("udp", &address)
	if err != nil {
		log.Fatal("Failed to start DNS server [:8090]")
	}

	for {
		// Allocate memory for byte array
		udpData := make([]byte, 1024)
		// Read UDP connection
		_, clientAddress, err := udp.ReadFrom(udpData)
		if err != nil {
			log.Print("Failed reading from UDP")
			continue
		}
		// Init DNS packet using UDP data
		packet := gopacket.NewPacket(udpData, layers.LayerTypeDNS, gopacket.Default)
		// The dns question (the "request") will be in the first layer of the created packet
		dnsQuestion := packet.Layer(layers.LayerTypeDNS)
		// Serve the DNS answer based on the question back to the client address using UDP
		serveDNS(udp, clientAddress, dnsQuestion.(*layers.DNS))
	}
}

func serveDNS(u *net.UDPConn, clientAddr net.Addr, request *layers.DNS) {
	// Resolve and parse IP
	domain := request.Questions[0].Name
	IP := data.GetIP(domain)
	parsedIP := parseIP(IP)

	// Create DNS answer
	dnsAnswer := layers.DNSResourceRecord{
		Type:  layers.DNSTypeA,
		IP:    parsedIP,
		Name:  []byte(domain),
		Class: layers.DNSClassIN,
	}

	// Create server reply, serialize and write back to client
	reply := createReply(request, dnsAnswer)
	buffer := serializeReply(reply)
	u.WriteTo(buffer.Bytes(), clientAddr)
}

func createReply(request *layers.DNS, dnsAnswer layers.DNSResourceRecord) *layers.DNS {
	reply := request
	reply.QR = true
	reply.ANCount = 1
	reply.OpCode = layers.DNSOpCodeNotify
	reply.AA = true
	reply.Answers = append(reply.Answers, dnsAnswer)
	reply.ResponseCode = layers.DNSResponseCodeNoErr
	return reply
}

func serializeReply(reply *layers.DNS) gopacket.SerializeBuffer {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	err := reply.SerializeTo(buf, opts)
	if err != nil {
		log.Panicf("Error while serializing DNS answer: [%+v]", err)
	}
	return buf
}

func parseIP(resolvedIP string) net.IP {
	parsedIP, _, err := net.ParseCIDR(resolvedIP + "/24")
	if err != nil {
		log.Fatalf("Failed parsing resolved IP: %s", parsedIP)
	}
	return parsedIP
}
