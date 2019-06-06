package qtrade

import (
	"net/http"
	"path"
)

type Balances struct {
	Data struct {
		Balances []struct {
			Balance  float64 `json:"balance,string"`
			Currency string  `json:"currency"`
		} `json:"balances"`
	} `json:"data"`
}

func (p *Client) Balances() (*Balances, error) {
	p.URL.Path = path.Join(VERSION, "user", "balances")
	res, err := p.do(http.MethodGet, p.URL.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	b := new(Balances)
	decoder(res, &b)

	return b, err
}
