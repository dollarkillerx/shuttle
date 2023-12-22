package faker

import (
	"sync"
	"time"
)

type defMapDumper struct {
	realIPMap map[string]*FakeMap
	fakeIPMap map[string]*FakeMap
	domainMap map[string]*FakeMap

	mutex          sync.RWMutex
	expirationTime int64
}

func newDefMapDumper() *defMapDumper {
	return &defMapDumper{
		realIPMap:      make(map[string]*FakeMap, 1024),
		fakeIPMap:      make(map[string]*FakeMap, 1024),
		domainMap:      make(map[string]*FakeMap, 1024),
		mutex:          sync.RWMutex{},
		expirationTime: 0,
	}
}

// Load try decide string k should be ip or domain, then try fetch FakeMap and check expiration
func (dumper *defMapDumper) Load(k string) (res *FakeMap) {
	// 过期检测
	defer func() {
		if res != nil {
			if dumper.elementAvailable(res) {
				return
			}
			dumper.DelByFakeMap(res)
		}
	}()

	dumper.mutex.RLock()
	defer dumper.mutex.RUnlock()

	if _IP_REGEX.MatchString(k) {
		// 先尝试 fake ip
		if res, ok := dumper.fakeIPMap[k]; ok {
			return res
		}

		// 再尝试 real ip
		return dumper.realIPMap[k]
	}

	// 最后域名
	return dumper.domainMap[k]
}

// Save store FakeMap or update a stored FakeMap that part match the input FakeMap
func (dumper *defMapDumper) Save(fmap *FakeMap) {
	// 检测是否有已经存储的, 如果存在则只更新值
	orig := dumper.LoadByFakeMap(fmap)
	if orig != nil {
		dumper.mutex.Lock()
		defer dumper.mutex.Unlock()

		if fmap.domain == "" {
			orig.domain = fmap.domain
		}
		if fmap.realIP == "" {
			orig.realIP = fmap.realIP
		}
		if fmap.fakeIP == "" {
			orig.fakeIP = fmap.fakeIP
		}
		return
	}

	dumper.mutex.Lock()
	defer dumper.mutex.Unlock()

	// NOTE: key 不能为 ""
	if fmap.domain != "" {
		dumper.domainMap[fmap.domain] = fmap
	}
	if fmap.realIP != "" {
		dumper.realIPMap[fmap.realIP] = fmap
	}
	if fmap.fakeIP != "" {
		dumper.fakeIPMap[fmap.fakeIP] = fmap
	}
}

// LoadByFakeMap get full stored FakeMap by matching part of input FakeMap
func (dumper *defMapDumper) LoadByFakeMap(fmap *FakeMap) (res *FakeMap) {
	// 过期检测
	defer func() {
		if res != nil {
			if dumper.elementAvailable(res) {
				return
			}
			dumper.DelByFakeMap(res)
		}
	}()

	dumper.mutex.RLock()
	defer dumper.mutex.RUnlock()

	if fmap.domain != "" {
		if m, ok := dumper.domainMap[fmap.domain]; ok {
			return m
		}
	}

	if fmap.fakeIP != "" {
		if m, ok := dumper.fakeIPMap[fmap.fakeIP]; ok {
			return m
		}
	}

	if fmap.realIP != "" {
		if m, ok := dumper.realIPMap[fmap.realIP]; ok {
			return m
		}
	}

	return nil
}

// DelByFakeMap delete all key value binding in stored maps
func (dumper *defMapDumper) DelByFakeMap(fmap *FakeMap) {
	dumper.mutex.Lock()
	defer dumper.mutex.Unlock()

	if fmap.domain != "" {
		delete(dumper.domainMap, fmap.domain)
	}
	if fmap.realIP != "" {
		delete(dumper.realIPMap, fmap.realIP)
	}
	if fmap.fakeIP != "" {
		delete(dumper.fakeIPMap, fmap.fakeIP)
	}
}

// Del try decide string k should be ip or domain, then try delete FakeMap
func (dumper *defMapDumper) Del(key string) {
	if fmap := dumper.Load(key); fmap != nil {
		dumper.DelByFakeMap(fmap)
	}
}

// elementAvailable check FakeMap expiration
func (dumper *defMapDumper) elementAvailable(fmap *FakeMap) bool {
	return dumper.expirationTime > 0 && (time.Now().UnixNano()-fmap.TimeAdded) > dumper.expirationTime
}

// Dump add time to FakeMap and then save it
func (dumper *defMapDumper) Dump(fmap *FakeMap) error {
	fmap.TimeAdded = time.Now().UnixNano()
	dumper.Save(fmap)
	return nil
}

// LoadByDomain get FakeMap by domain key string
func (dumper *defMapDumper) LoadByDomain(domain string) (res *FakeMap) {
	// 过期检测
	defer func() {
		if res != nil {
			if dumper.elementAvailable(res) {
				return
			}
			dumper.DelByFakeMap(res)
		}
	}()

	dumper.mutex.RLock()
	defer dumper.mutex.RUnlock()

	if res, ok := dumper.domainMap[domain]; ok {
		return res
	}
	return nil
}

// LoadByFakeIP get FakeMap by fake ip key string
func (dumper *defMapDumper) LoadByFakeIP(ip string) (res *FakeMap) {
	// 过期检测
	defer func() {
		if res != nil {
			if dumper.elementAvailable(res) {
				return
			}
			dumper.DelByFakeMap(res)
		}
	}()

	dumper.mutex.RLock()
	defer dumper.mutex.RUnlock()

	if res, ok := dumper.fakeIPMap[ip]; ok {
		return res
	}
	return nil
}

// LoadByRealIP iget FakeMap by real ip key string
func (dumper *defMapDumper) LoadByRealIP(ip string) (res *FakeMap) {
	// 过期检测
	defer func() {
		if res != nil {
			if dumper.elementAvailable(res) {
				return
			}
			dumper.DelByFakeMap(res)
		}
	}()

	dumper.mutex.RLock()
	defer dumper.mutex.RUnlock()

	if res, ok := dumper.realIPMap[ip]; ok {
		return res
	}
	return nil
}

// Clean starts a goroutine to clean expirated FakeMap after every duration
func (dumper *defMapDumper) Clean(minutes time.Duration) error {
	go func() {
		for {
			time.Sleep(time.Minute * minutes)

			for _, fmap := range dumper.fakeIPMap {
				if dumper.elementAvailable(fmap) {
					continue
				}
				dumper.DelByFakeMap(fmap)
			}
		}
	}()
	return nil
}
