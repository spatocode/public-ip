package publicip

import (
	"errors"
	"net"
	"time"

	"github.com/miekg/dns"
)

type dnsServer struct {
	servers  []string
	question question
}

type question struct {
	name      string
	qType     uint16
}

type dnsData struct {
	dnsServers   []dnsServer
}

func initDnsData(version string) dnsData {
	googleServers := dnsServer{
		servers: v4googlednsServers,
		question: question{
			name: googlednsServerName,
			qType: googlednsServerType,
		},
	}

	opendnsServers := dnsServer{
		servers: v4opendnsServers,
		question: question{
			name: opendnsServerName,
			qType: v4opendnsServerType,
		},
	}

	data := dnsData{
		dnsServers: []dnsServer{googleServers, opendnsServers},
	}
	return data
}

const (
	dnsTimeout  = 3 * time.Second
	opendnsServerName = "myip.opendns.com."
	v4opendnsServerType = dns.TypeA
	googlednsServerName = "o-o.myaddr.l.google.com."
	googlednsServerType = dns.TypeTXT
)

var (
	v4opendnsServers = []string{
		"208.67.222.222",
		"208.67.220.220",
		"208.67.222.220",
		"208.67.220.222",
	}
	v4googlednsServers = []string{
		"216.239.32.10",
		"216.239.34.10",
		"216.239.36.10",
		"216.239.38.10",
	}
)

func queryDNS(version string) (net.IP, error) {
	var err error
	var ip net.IP

	data := initDnsData(version)
	client := &dns.Client{Timeout: dnsTimeout}

	for _, dnsServer := range data.dnsServers {
		for _, server := range dnsServer.servers {
			if ip != nil {
				break
			}
			message := &dns.Msg{}
			message.SetQuestion(dnsServer.question.name, dnsServer.question.qType)
			msg, _, err := client.Exchange(message, net.JoinHostPort(server, "53"))
			if err != nil {
				return nil, err
			}

			switch dnsServer.question.qType {
			case dns.TypeA:
				for _, ans := range msg.Answer {
					if a, ok := ans.(*dns.A); ok {
						ip = a.A
						break
					}
				}
			case dns.TypeTXT:
				for _, ans := range msg.Answer {
					if txt, ok := ans.(*dns.TXT); ok {
						for _, s := range txt.Txt {
							if ip = net.ParseIP(s); ip != nil {
								break
							}
						}
					}
				}
			default:
				return nil, errors.New("invalid record type")
			}
		}
	}

	return ip, err
}

func V4() (net.IP, error) {
	return queryDNS("v4")
}
