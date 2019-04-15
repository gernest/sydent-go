package clients

import (
	"fmt"
	"net"
	"net/url"
	"testing"
)

func TestSrvResolver(t *testing.T) {
	resolve := SrvResolver()
	s, err := resolve("matrix.org")
	if err != nil {
		t.Fatal(err)
	}
	if s.Records == nil {
		t.Error("expected srv records")
	}
	cname := "_matrix._tcp.matrix.org."
	if s.CName != cname {
		t.Errorf("expected cname to be %s got %s", cname, s.CName)
	}
}

func TestRouteMatrixURI(t *testing.T) {
	target := "matrixid-example.org"
	port := uint16(8000)
	sample := []struct {
		desc   string
		url    string
		expect string
		well   bool
	}{
		{"ip no port", "https://192.0.2.1", "192.0.2.1,192.0.2.1,8448", false},
		{"ip with port", "https://192.0.2.1:70", "192.0.2.1,192.0.2.1,70", false},
		{"host with port", "https://localhost:70", "localhost,localhost,70", false},
		{"host no port", "https://localhost", fmt.Sprintf("localhost,%s,%d", target, port), false},
	}

	fmtResult := func(info *RoutingInfo, err error) string {
		if err != nil {
			return fmt.Sprintf("error=%q", err.Error())
		}
		return fmt.Sprintf("%s,%s,%d", info.TLSServerName, info.TargetHost, info.TargetPort)
	}

	resolve := func(addr string) (*Server, error) {
		return &Server{
			Records: []*net.SRV{
				{Target: target, Port: port},
			},
		}, nil
	}
	for _, v := range sample {
		u, err := url.Parse(v.url)
		if err != nil {
			t.Fatal(err)
		}
		got := fmtResult(RouteMatrixURI(resolve, u, v.well))
		if got != v.expect {
			t.Errorf("%s: expected %s got %s", v.desc, v.expect, got)
		}
	}
}
