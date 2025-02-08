package peasant

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	peasant "github.com/candango/gopeasant"
)

type CafofoTransport struct {
	*peasant.HttpTransport
}

func NewCafofoTransport(tr *peasant.HttpTransport) *CafofoTransport {
	return &CafofoTransport{tr}
}

func (ct *CafofoTransport) Auth() error {
	nonce, err := ct.NewNonce()
	if err != nil {
		return err
	}

	d, err := ct.Directory()
	if err != nil {
		return err
	}

	ds, ok := d["security"].(map[string]any)
	if !(ok) {
		return errors.New("error converting security directory")
	}

	data := url.Values{}
	secret := "secret correto"
	data.Set("secret", secret)
	req, err := http.NewRequest(http.MethodPost, ds["auth"].(string), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("nonce", nonce)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := ct.Client.Do(req)
	if err != nil {
		return err
	}

	b, err := peasant.BodyAsString(res)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return errors.New(res.Status)
	}
	fmt.Println(b)
	return nil
}

type CafofoDirectoryProvider struct {
	directory map[string]any
	url       string
	http.Client
}

func NewCafofoDirectoryProvider(url string) *CafofoDirectoryProvider {
	return &CafofoDirectoryProvider{url: url}
}

func (p *CafofoDirectoryProvider) Directory() (map[string]any, error) {
	if p.directory == nil {
		path := fmt.Sprintf("%s/directory/", p.GetUrl())
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		res, err := p.Client.Do(req)
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
		p.directory = dir
	}

	//TODO: Expire the directory
	return p.directory, nil
}

func (p *CafofoDirectoryProvider) GetUrl() string {
	return p.url
}

type CafofoPeasant struct {
	peasant.Peasant
}

func NewCandangoPesant(p peasant.Peasant) *CafofoPeasant {
	return &CafofoPeasant{p}
}
