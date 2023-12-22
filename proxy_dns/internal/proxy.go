package internal

import (
	"fmt"
	"net"
	"regexp"

	"github.com/miekg/dns"
)

type DNSProxy struct {
	Cache         *Cache
	domains       map[string]interface{}
	servers       map[string]interface{}
	defaultServer string
}

func NewDNSProxy(cache *Cache, domains map[string]interface{}, servers map[string]interface{}, defaultServer string) *DNSProxy {
	return &DNSProxy{
		Cache:         cache,
		domains:       domains,
		servers:       servers,
		defaultServer: defaultServer,
	}
}

func (d *DNSProxy) GetResponse(requestMsg *dns.Msg) (*dns.Msg, error) {
	responseMsg := new(dns.Msg)
	if len(requestMsg.Question) > 0 {
		question := requestMsg.Question[0]

		dnsServer := d.getIPFromConfigs(question.Name, d.servers)
		if dnsServer == "" {
			dnsServer = d.defaultServer
		}

		switch question.Qtype {
		case dns.TypeA:
			answer, err := d.ProcessTypeA(dnsServer, &question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, answer.Answer...)
			responseMsg.Ns = append(responseMsg.Ns, answer.Ns...)
			responseMsg.Extra = append(responseMsg.Extra, answer.Extra...)
		default:
			answer, err := d.ProcessOtherTypes(dnsServer, &question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			responseMsg.Answer = append(responseMsg.Answer, answer.Answer...)
			responseMsg.Ns = append(responseMsg.Ns, answer.Ns...)
			responseMsg.Extra = append(responseMsg.Extra, answer.Extra...)
		}
	}

	return responseMsg, nil
}

func (d *DNSProxy) ProcessOtherTypes(dnsServer string, q *dns.Question, requestMsg *dns.Msg) (*dns.Msg, error) {
	queryMsg := new(dns.Msg)
	requestMsg.CopyTo(queryMsg)
	queryMsg.Question = []dns.Question{*q}

	msg, err := lookup(dnsServer, queryMsg)
	if err != nil {
		return nil, err
	}

	if len(msg.Answer) > 0 {
		return msg, nil
	}
	return nil, fmt.Errorf("not found")
}

func (d *DNSProxy) ProcessTypeA(dnsServer string, q *dns.Question, requestMsg *dns.Msg) (*dns.Msg, error) {
	ip := d.getIPFromConfigs(q.Name, d.domains)
	cacheMsg, found := d.Cache.Get(q.Name)

	if ip == "" && !found {
		queryMsg := new(dns.Msg)
		requestMsg.CopyTo(queryMsg)
		queryMsg.Question = []dns.Question{*q}

		msg, err := lookup(dnsServer, queryMsg)
		if err != nil {
			return nil, err
		}

		if len(msg.Answer) > 0 {
			d.Cache.Set(q.Name, msg)
			return msg, nil
		}
	} else if found {
		return cacheMsg.(*dns.Msg), nil
	} else if ip != "" {
		answer, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
		if err != nil {
			return nil, err
		}

		queryMsg := new(dns.Msg)
		requestMsg.CopyTo(queryMsg)
		queryMsg.Question = []dns.Question{*q}
		queryMsg.Answer = []dns.RR{answer}
		return queryMsg, nil
	}
	return nil, fmt.Errorf("not found")
}

func (d *DNSProxy) getIPFromConfigs(domain string, configs map[string]interface{}) string {
	for k, v := range configs {
		match, _ := regexp.MatchString(k+"\\.", domain)
		if match {
			return v.(string)
		}
	}
	return ""
}

func GetOutboundIP() (net.IP, error) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func lookup(server string, m *dns.Msg) (*dns.Msg, error) {
	dnsClient := new(dns.Client)
	dnsClient.Net = "udp"

	response, _, err := dnsClient.Exchange(m, server)
	if err != nil {
		return nil, err
	}

	return response, nil
}
