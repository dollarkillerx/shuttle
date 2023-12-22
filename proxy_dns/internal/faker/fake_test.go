package faker

import (
	"fmt"
	"sync"
	"testing"
)

func Test(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 10; i++ {
			fmap := EmptyFakeMap()

			fmap.SetDomain("goroutine1")
			fmap.SetFakeIP(fmt.Sprintf("127.0.0.%d", i))
			fmap.SetReadIP(fmt.Sprintf("127.0.0.%d", i))
			Dump(fmap)

			t.Logf("goroutine1: %#v", Load("goroutine1"))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 20; i++ {
			fmap := EmptyFakeMap()

			fmap.SetDomain("goroutine2")
			fmap.SetFakeIP(fmt.Sprintf("192.168.%d.10", i))
			fmap.SetReadIP(fmt.Sprintf("192.%d.1.10", i))
			Dump(fmap)

			t.Logf("goroutine2: %#v", Load("goroutine2"))
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 9; i++ {
			fmap := EmptyFakeMap()

			fmap.SetFakeIP(fmt.Sprintf("192.%d.1.10", i))
			fmap.SetReadIP(fmt.Sprintf("192.%d.1.10", i))
			Dump(fmap)

			t.Logf("goroutine3: %#v", LoadFakeMapByFakeIP(fmt.Sprintf("192.%d.1.10", i)))
		}
	}()

	wg.Wait()
}
