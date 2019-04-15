package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/gernest/sydent-go/config"
)

const (
	serviceName  = "matrix"
	serviceProto = "tcp"
)

// Server stores SRV records of a host that supports matrix federation.
type Server struct {
	CName   string
	Records []*net.SRV
}

// Pick returns a random srv record.
func (s *Server) Pick() (*net.SRV, error) {
	if s.Records == nil {
		return nil, errors.New("matrixid: no records to pick from")
	}
	if len(s.Records) == 1 {
		return s.Records[0], nil
	}
	n := rand.Intn(len(s.Records))
	return s.Records[n], nil
}

// SRVResolveFunc is a function used to resolve srv records offederated matrix
// servers
type SRVResolveFunc func(address string) (*Server, error)

// SrvResolver returns a function that lookup for srv records of a host
// configured for matrix federation.
//
// This caches the records, and the returned function is safe for concurrent use.
func SrvResolver() SRVResolveFunc {
	var mu sync.RWMutex
	cache := make(map[string]*Server)
	return func(host string) (*Server, error) {
		mu.RLock()
		v, ok := cache[host]
		if ok {
			return v, nil
		}
		mu.RUnlock()
		cname, addrs, err := net.LookupSRV(serviceName, serviceProto, host)
		if err != nil {
			return nil, err
		}
		s := &Server{
			CName:   cname,
			Records: addrs,
		}
		mu.Lock()
		cache[host] = s
		mu.Unlock()
		return s, nil
	}
}

// RoutingInfo Contains the parameters needed to direct a federation connection
// to a particular  server.
// Where a SRV record points to several servers, this object contains a single server
// chosen from the list.
type RoutingInfo struct {
	// The value we should assign to the Host header (host:port from the matrix
	// URI, or .well-known).
	HostHeader string
	// The server name we should set in the SNI (typically host, without port, from
	// the  matrix URI or .well-known)
	TLSServerName string
	// The hostname (or IP literal) we should route the TCP connection to (the
	// target of the  SRV record, or the hostname from the URL/.well-known)
	TargetHost string
	// The port we should route the TCP connection to (the target of the SRV
	// record, or the port from the URL/.well-known, or 8448)
	TargetPort uint16
}

// WellKnownResponse is a response returned from a federated matrix server for
// queries on /.well-known/matrix/server path.
type WellKnownResponse struct {
	Server string `json:"m.server"`
}

// GetWellKnown sends a GET request to address looking for .well=known host
// address.
//
// No caching is performed.
func GetWellKnown(address string) (*WellKnownResponse, error) {
	uri := fmt.Sprintf("https://%s/.well-known/matrix/server", address)
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var o WellKnownResponse
	err = json.Unmarshal(b, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// RouteMatrixURI returns RoutingInfo that contains information about which
// server the request should be sent to.
//
// resolve is used to resolve hostname for non IP hosts.
func RouteMatrixURI(resolve SRVResolveFunc, u *url.URL, lookupWellKnown bool) (*RoutingInfo, error) {
	host := u.Hostname()
	if ip := net.ParseIP(host); ip != nil {
		var port uint16
		p := u.Port()
		if p == "" {
			port = 8448
		} else {
			up, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			port = uint16(up)
		}
		return &RoutingInfo{
			HostHeader:    u.Host,
			TLSServerName: host,
			TargetHost:    host,
			TargetPort:    port,
		}, nil
	}
	if port := u.Port(); port != "" {
		up, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		return &RoutingInfo{
			HostHeader:    u.Host,
			TLSServerName: host,
			TargetHost:    host,
			TargetPort:    uint16(up),
		}, nil
	}
	if lookupWellKnown {
		// TODO : add well known
	}
	s, err := resolve(host)
	if err != nil {
		return nil, err
	}
	t, err := s.Pick()
	if err != nil {
		return nil, err
	}
	return &RoutingInfo{
		HostHeader:    u.Host,
		TLSServerName: host,
		TargetHost:    t.Target,
		TargetPort:    t.Port,
	}, nil
}

// NewFederatedClient returns a http.Client that uses the federated transport.
func NewFederatedClient() *http.Client {
	return &http.Client{
		Transport: NewFederatedTripper(config.MaxRetries),
	}
}

// Fed is a global federated matrix client. Safe for concurrent use.
var Fed = NewFederatedClient()
