package faker

import (
	"fmt"
	"regexp"
	"time"
)

var _IP_REGEX = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))$`)

type MapDumper interface {
	Load(string) *FakeMap
	Dump(*FakeMap) error
	LoadByDomain(string) *FakeMap
	LoadByRealIP(string) *FakeMap
	LoadByFakeIP(string) *FakeMap
	Clean(time.Duration) error
}

var defaultMapDumper MapDumper = nil

func init() {
	defaultMapDumper = newDefMapDumper()
}

func SetDefaultMapDumper(dumper MapDumper) {
	defaultMapDumper = dumper
}

// FakeMap stores domain, real ip and fake ip
type FakeMap struct {
	domain    string
	realIP    string
	fakeIP    string
	TimeAdded int64
}

func EmptyFakeMap() *FakeMap {
	return &FakeMap{}
}

func LoadFakeMapByDomain(domain string) *FakeMap {
	return defaultMapDumper.LoadByDomain(domain)
}

func LoadFakeMapByReadIP(ip string) *FakeMap {
	return defaultMapDumper.LoadByRealIP(ip)
}

func LoadFakeMapByFakeIP(ip string) *FakeMap {
	return defaultMapDumper.LoadByFakeIP(ip)
}

func Load(k string) *FakeMap {
	if fmap := defaultMapDumper.Load(k); fmap != nil {
		return fmap
	}
	return nil
}

func Dump(m *FakeMap) error {
	return defaultMapDumper.Dump(m)
}

func (m *FakeMap) GetAll() (string, string, string) {
	return m.domain, m.realIP, m.fakeIP
}

func (m *FakeMap) GetDomain() string {
	return m.domain
}

func (m *FakeMap) SetDomain(domain string) {
	m.domain = domain
}

func (m *FakeMap) GetReadIP() string {
	return m.realIP
}

func (m *FakeMap) SetReadIP(ip string) error {
	if _IP_REGEX.MatchString(ip) {
		m.realIP = ip
		return nil
	}
	return fmt.Errorf("%q is not ip", ip)
}

func (m *FakeMap) GetFakeIP() string {
	return m.fakeIP
}

func (m *FakeMap) SetFakeIP(ip string) error {
	if _IP_REGEX.MatchString(ip) {
		m.fakeIP = ip
		return nil
	}
	return fmt.Errorf("%q is not ip", ip)
}
