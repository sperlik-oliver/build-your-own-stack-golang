package data

import "log"

var records map[string]string = map[string]string{
	"google.com": "216.58.196.142",
	"amazon.com": "176.32.103.205",
}

func GetIP(domain []byte) string {
	ip, isSuccess := records[string(domain)]
	if !isSuccess {
		log.Fatalf("Couldnt find IP based on record")
	}

	return ip
}
