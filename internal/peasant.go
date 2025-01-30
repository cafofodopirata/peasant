package peasant

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	peasant "github.com/candango/gopeasant"
)

type CafofoTransport struct {
	*peasant.HttpTransport
}

func NewCafofoTransport(tr *peasant.HttpTransport) *CafofoTransport {
	return &CafofoTransport{tr}
}

func (tt *CafofoTransport) Directory() (map[string]interface{}, error) {
	path := fmt.Sprintf("%s/directory/", tt.Url)
	fmt.Println(path)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	res, err := tt.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, errors.New(res.Status)
	}
	var dir map[string]any
	err = peasant.BodyAsJson(res, &dir)
	if err != nil {
		return nil, err
	}

	return dir, nil
}

func (tt *CafofoTransport) NewNonceUrl() (string, error) {
	d, err := tt.Directory()
	if err != nil {
		return "", err
	}
	return d["new-nonce"].(string), nil
}

func (tt *CafofoTransport) NewNonce() (string, error) {
	url, err := tt.NewNonceUrl()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}
	res, err := tt.Client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 {
		return "", errors.New(res.Status)
	}
	return tt.ResolveNonce(res), nil
}

func (tt *CafofoTransport) DoSomething(t *testing.T) (string, error) {
	nonce, err := tt.NewNonce()
	if err != nil {
		return "", err
	}

	d, err := tt.Directory()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, d["doSomething"].(string), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("nonce", nonce)

	res, err := tt.Client.Do(req)
	if err != nil {
		return "", err
	}

	b, err := peasant.BodyAsString(res)
	if err != nil {
		return "", err
	}

	if res.StatusCode > 299 {
		return "", errors.New(res.Status)
	}
	return b, nil
}

type CafofoPeasant struct {
	peasant.Peasant
}

func NewCandangoPesant(p peasant.Peasant) *CafofoPeasant {
	return &CafofoPeasant{p}
}

func (p *CafofoPeasant) NewNonce(t *testing.T) (string, error) {
	return p.Transport.(*CafofoTransport).DoSomething(t)
}
