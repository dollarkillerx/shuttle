package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/dollarkillerx/urllib"
	"google.dev/google/shuttle/core/app/client"
	"google.dev/google/shuttle/core/app/client/conf"
)

func TestMain(m *testing.M) {
	client.RouterRegister()
	addr := "127.0.0.1:8985"
	server := &http.Server{Addr: addr, Handler: nil}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	os.Exit(m.Run())
}

func TestClientRun(t *testing.T) {
	i, bytes, err := urllib.Post("http://127.0.0.1:8985/link").SetJsonObject(conf.ClientConfig{
		Server: conf.Server{Protocol: "wss", Addr: "192.227.234.228:8283"},
		HTTP:   conf.HTTP{},
		WS:     conf.WS{Path: "/56801e42-10a2-4246-9d5e-0643a2668f47"},
		Pac:    false,
	}).Byte()

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(i)
	t.Log(string(bytes))

	i, bytes, err = urllib.Get("http://127.0.0.1:8985/status").Byte()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(i)
	t.Log(string(bytes))
}

func TestClientStatus(t *testing.T) {
	i, bytes, err := urllib.Get("http://127.0.0.1:8985/status").Byte()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(i)
	t.Log(string(bytes))
}

func TestClientLink(t *testing.T) {
	i, bytes, err := urllib.Post("http://127.0.0.1:8985/link").SetJsonObject(conf.ClientConfig{
		Server: struct {
			Protocol string `toml:"protocol" json:"protocol"`
			Addr     string `toml:"address" json:"addr"`
		}(struct {
			Protocol string `toml:"protocol"`
			Addr     string `toml:"address"`
		}{
			Protocol: "wss",
			Addr:     "192.227.234.228:8283",
			//ManagerRPCAddr: "127.0.0.1:8283",
		}),
		HTTP: struct {
			Path string `toml:"path" default:"/" json:"path"`
		}(struct {
			Path string `toml:"path" default:"/"`
		}{}),
		WS: struct {
			Path string `toml:"path" default:"/" json:"path"`
		}(struct {
			Path string `toml:"path" default:"/"`
		}{
			Path: "/56801e42-10a2-4246-9d5e-0643a2668f47",
		}),
		Pac: false,
	}).Byte()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(i)
	t.Log(string(bytes))
}

func TestClientStop(t *testing.T) {
	i, bytes, err := urllib.Post("http://127.0.0.1:8985/stop").Byte()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(i)
	t.Log(string(bytes))
}
