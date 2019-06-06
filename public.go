package qtrade

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"
)

type Ticker struct {
	Data struct {
		Ask             float64 `json:"ask,string"`
		Bid             float64 `json:"bid,string"`
		DayAvgPrice     float64 `json:"day_avg_price,string"`
		DayChange       float64 `json:"day_change,string"`
		DayHigh         float64 `json:"day_high,string"`
		DayLow          float64 `json:"day_low,string"`
		DayOpen         float64 `json:"day_open,string"`
		DayVolumeBase   float64 `json:"day_volume_base,string"`
		DayVolumeMarket float64 `json:"day_volume_market,string"`
		ID              int     `json:"id"`
		IDHr            string  `json:"id_hr"`
		Last            float64 `json:"last,string"`
	} `json:"data"`
}

func (p *Client) Ticker(code string) (*Ticker, error) {
	p.URL.Path = path.Join(VERSION, "ticker", code)
	res, err := p.do(http.MethodGet, p.URL.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	t := new(Ticker)
	decoder(res, &t)

	return t, err
}

type Orderbook struct {
	Data struct {
		BestAsk, BestBid float64
		Buy              Books `json:"buy"`
		Sell             Books `json:"sell"`
	} `json:"data"`
}

type Books struct {
	Books []Book
}

type Book struct {
	Price float64
	Size  float64
}

func (p *Client) Orderbook(code string) (*Orderbook, error) {
	p.URL.Path = path.Join(VERSION, "orderbook", code)
	res, err := p.do(http.MethodGet, p.URL.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	o := new(Orderbook)
	decoder(res, &o)

	o.getSpread()

	return o, err
}

func (p *Books) UnmarshalJSON(b []byte) error {
	var j map[string]string
	json.Unmarshal(b, &j)

	for key, val := range j {
		price, err := strconv.ParseFloat(key, 64)
		if err != nil {
			continue
		}
		size, err := strconv.ParseFloat(val, 64)
		if err != nil {
			continue
		}
		p.Books = append(p.Books, Book{
			Price: price,
			Size:  size,
		})
	}

	return nil
}

func (p *Orderbook) getSpread() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		var bid float64
		for _, v := range p.Data.Buy.Books {
			if bid < v.Price {
				bid = v.Price
			}
		}
		p.Data.BestBid = bid
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ask := math.Inf(1)
		for _, v := range p.Data.Sell.Books {
			if ask > v.Price {
				ask = v.Price
			}
		}
		p.Data.BestAsk = ask
		wg.Done()
	}()

	wg.Wait()
}

type Executions struct {
	Data struct {
		Trades []struct {
			Amount      float64   `json:"amount,string"`
			CreatedAt   time.Time `json:"created_at"`
			Price       float64   `json:"price,string"`
			SellerTaker bool      `json:"seller_taker"`
		} `json:"trades"`
	} `json:"data"`
}

func (p *Client) Executions(code string) (*Executions, error) {
	p.URL.Path = fmt.Sprintf("%s/market/%s/trades", VERSION, code)
	res, err := p.do(http.MethodGet, p.URL.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	e := new(Executions)
	decoder(res, &e)

	return e, err
}
