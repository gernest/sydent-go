package clients

import (
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff"
)

// FederatedTripper implements http.RoundTripper interface with support for exponential
// backoff retries. This is used to talk to federated servers
type FederatedTripper struct {
	backoff   backoff.BackOff
	round     http.RoundTripper
	resolve   SRVResolveFunc
	wellKnown bool
}

func NewFederatedTripper(maxRetries uint64) *FederatedTripper {
	b := backoff.NewExponentialBackOff()
	return &FederatedTripper{
		backoff: backoff.WithMaxRetries(b, maxRetries),
		round:   http.DefaultTransport,
		resolve: SrvResolver(),
	}
}

func (tr *FederatedTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var ri *RoutingInfo
	var terr error
	err := backoff.Retry(func() error {
		ri, terr = RouteMatrixURI(tr.resolve, req.URL, tr.wellKnown)
		if terr != nil {
			return terr
		}
		if h := req.Header.Get("Host"); h == "" {
			req.Header.Set("Host", ri.HostHeader)
		}
		req.URL.Host = fmt.Sprintf("%s:%d", ri.TargetHost, ri.TargetPort)
		res, terr = tr.round.RoundTrip(req)
		if terr != nil {
			return terr
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return fmt.Errorf("matrixid/clients:received not 200 status code :%d", res.StatusCode)
		}
		return nil
	}, tr.backoff)
	if err != nil {
		return nil, err
	}
	return res, nil
}
