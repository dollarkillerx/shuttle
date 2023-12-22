package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dollarkillerx/urllib"
	"github.com/oschwald/geoip2-golang"
	"github.com/sourcegraph/conc/pool"
	"github.com/txthinking/socks5"
	"google.dev/google/socks5_discovery/internal/conf"
	"google.dev/google/socks5_discovery/proto"
)

type Service struct {
	proto.UnimplementedSocks5DiscoveryServer

	conf     conf.S5DiscoveryConfig
	muGeoip2 sync.Mutex
	geoip2   *geoip2.Reader

	muProxiesToVerify sync.Mutex
	proxiesToVerify   []*proto.Socks5
}

func NewService(conf conf.S5DiscoveryConfig) *Service {
	ser := &Service{conf: conf}

	db, err := geoip2.Open("./resource/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	ser.geoip2 = db

	var s5Chan = make(chan *proto.Socks5, 100)
	go ser.core(s5Chan)
	go ser.toClean(s5Chan)
	go ser.keepAlive()

	return ser
}

func (s *Service) getIOSCountry(ip string) string {
	s.muGeoip2.Lock()
	defer s.muGeoip2.Unlock()

	index := strings.Index(ip, ":")
	if index != -1 {
		ip = ip[:index]
	}

	ip2 := net.ParseIP(ip)
	record, err := s.geoip2.City(ip2)
	if err != nil {
		return "ad"
	}
	return record.Country.IsoCode
}

func (s *Service) core(s5Chan chan *proto.Socks5) {
	for {
		uris := []string{
			"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks5.txt",
			"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks5.txt",
			"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS5_RAW.txt",
			"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
			"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt",
		}

		var proxy = func(uri string) {
			code, byt, err := urllib.Get(uri).Byte()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second)
				return
			}

			if code != 200 {
				log.Println("core code : ", code)
				time.Sleep(time.Second)
				return
			}

			split := strings.Split(string(byt), "\n")
			for _, v := range split {
				v = strings.TrimSpace(v)
				if v != "" {
					s5Chan <- &proto.Socks5{
						Country:  s.getIOSCountry(v),
						Address:  v,
						Username: "",
						Password: "",
						Delay:    0,
					}
				}
			}
		}

		for _, v := range uris {
			proxy(v)
		}

		time.Sleep(time.Minute)
	}
}

func (s *Service) toClean(s5Chan chan *proto.Socks5) {
	//var getProxiesToVerify = func() []*proto.Socks5 {
	//	s.muProxiesToVerify.Lock()
	//	defer s.muProxiesToVerify.Unlock()
	//
	//	return s.proxiesToVerify
	//}

	p := pool.New().WithMaxGoroutines(50)
	for {
		select {
		case item := <-s5Chan:
			//verify := getProxiesToVerify()
			//for _,v := range verify {
			//	if v.Address == item.Address {
			//		continue
			//	}
			//}

			p.Go(func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from panic:", r)
						debug.PrintStack()
					}
				}()

				delay, err := s.checkOk(item)
				if err != nil {
					log.Println(err)
					return
				}

				if delay == 0 {
					return
				}
				item.Delay = delay

				s.muProxiesToVerify.Lock()
				defer s.muProxiesToVerify.Unlock()

				s.proxiesToVerify = append(s.proxiesToVerify, item)
				newS5 := make([]*proto.Socks5, 0)
				newS5Map := map[string]struct{}{}
				for i, v := range s.proxiesToVerify {
					idx := i
					_, ex := newS5Map[v.Address]
					if !ex {
						newS5 = append(newS5, s.proxiesToVerify[idx])
						newS5Map[v.Address] = struct{}{}
					}
				}

				s.proxiesToVerify = newS5
			})
		}
	}
}

func (s *Service) keepAlive() {

	var getProxiesToVerify = func() []*proto.Socks5 {
		s.muProxiesToVerify.Lock()
		defer s.muProxiesToVerify.Unlock()

		return s.proxiesToVerify
	}

	var mu sync.Mutex
	var s5s []*proto.Socks5

	for {
		proxies := getProxiesToVerify()
		if len(proxies) == 0 {
			time.Sleep(time.Second)
			continue
		}

		p := pool.New().WithMaxGoroutines(50)

		for i := range proxies {
			idx := i
			p.Go(func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from panic:", r)
						debug.PrintStack()
					}
				}()

				//client, err := socks5.NewClient(proxies[idx].Address, proxies[idx].Username, proxies[idx].Password, 10, 60)
				//if err != nil {
				//	log.Println(err)
				//	return
				//}
				//defer client.Close()

				delay, err := s.checkOk(proxies[idx])
				if err != nil {
					log.Println(err)
					return
				}
				if delay == 0 {
					return
				}

				var setS5 = func() {
					mu.Lock()
					defer mu.Unlock()

					s5s = append(s5s, &proto.Socks5{
						Country:  proxies[idx].Country,
						Address:  proxies[idx].Address,
						Username: proxies[idx].Username,
						Password: proxies[idx].Password,
						Delay:    delay,
					})
				}

				setS5()
			})
		}

		p.Wait()

		var setPPP = func() {
			s.muProxiesToVerify.Lock()
			defer s.muProxiesToVerify.Unlock()

			var newS5 []*proto.Socks5
			var fMap = map[string]struct{}{}
			for i, v := range s5s {
				_, ex := fMap[v.Address]
				if !ex {
					newS5 = append(newS5, s5s[i])
					fMap[v.Address] = struct{}{}
				}
			}

			s.proxiesToVerify = newS5
			s5s = []*proto.Socks5{}
		}

		setPPP()

		time.Sleep(time.Minute * 2)
	}
}

func (s *Service) checkOk(proxies *proto.Socks5) (delay int64, err error) {
	testUri := "https://dns.google/"

	u, _ := url.Parse(testUri)

	req := &http.Request{
		URL:        u,
		Method:     "get",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
	}

	transport := &http.Transport{
		TLSHandshakeTimeout: time.Second * 11,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 100,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			client, err := socks5.NewClient(proxies.Address, proxies.Username, proxies.Password, 10, 60)
			if err != nil {
				return nil, err
			}

			return client.Dial(network, addr)
		},
	}

	httpClient := http.Client{
		Transport: transport,
		Timeout:   time.Second * 11,
	}
	now := time.Now().UnixMilli()
	do, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	do.Body.Close()
	rNow := time.Now().UnixMilli()

	return rNow - now, nil
}

func (s *Service) Discovery(ctx context.Context, request *proto.DiscoveryRequest) (*proto.DiscoveryResponse, error) {
	s.muProxiesToVerify.Lock()
	defer s.muProxiesToVerify.Unlock()

	var newS5 []*proto.Socks5
	var fMap = map[string]struct{}{}
	for i, v := range s.proxiesToVerify {
		_, ex := fMap[v.Address]
		if !ex {
			newS5 = append(newS5, s.proxiesToVerify[i])
			fMap[v.Address] = struct{}{}
		}
	}

	s.proxiesToVerify = newS5

	st := &proto.Socks5List{
		Socks5: newS5,
	}
	sort.Sort(st)

	return &proto.DiscoveryResponse{
		Socks5S: st.Socks5,
	}, nil
}
