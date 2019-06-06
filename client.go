package qtrade

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/labstack/gommon/log"

	"net/http"
	"net/url"
	"time"
)

/*
Markets
	MarketID	Pair		Maker Fee	Taker Fee
	1			BTC/LTC		0%			0.5%
	7			BTC/RVN		0%			0.5%
	9			BTC/BTM		0%			0.5%
	13			BTC/DEFT	0%			0.5%
	15			BTC/VEO		0%			1%
	19			BTC/SNOW	0%			1.5%
	20			BTC/BIS		0%			0.5%
	21			BTC/PHL		0%			0.5%
	24			BTC/NYZO	0%			1.5%
	25			BTC/TAO		0%			0.75%
	26			BTC/XTRI	0%			0.5%
	27			BTC/VLS		0%			0.5%
	28			BTC/ZANO	0%			0.5%
	30			BTC/PASC	0%			0.5%
*/

const (
	URL     = "https://api.qtrade.io/"
	VERSION = "v1"

	VEOBTC = "VEO_BTC"
	LTCBTC = "LTC_BTC"
)

type Client struct {
	Key, Secret string

	URL        *url.URL
	HTTPClient *http.Client
}

func New(key, secret string) *Client {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}

	client := http.DefaultClient
	tr := &http.Transport{
		MaxIdleConnsPerHost: 24,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	return &Client{
		Key:    key,
		Secret: secret,

		URL:        u,
		HTTPClient: client,
	}
}

// Request do
func (p *Client) do(method, url string, params *map[string]string, body io.Reader) (*http.Response, error) {
	head := http.Header{}
	head.Set("Content-Type", "application/json")

	if params != nil {
		q := p.URL.Query()
		for k, v := range *params {
			q.Set(k, v)
		}
		p.URL.RawQuery = q.Encode()
	}

	if strings.Contains(url, "user") { // Private API 処理
		timestamp := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
		// sets Header
		head.Set("Authorization", fmt.Sprintf("HMAC-SHA256 %s:%s", p.Key, p.signiture(method, timestamp, body)))
		head.Set("HMAC-Timestamp", timestamp)
		head.Set("Access-Control-Allow-Origin", "https://qtrade.io")
	}

	req, err := http.NewRequest(method, p.URL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = head

	res, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%d, response error status %s", res.StatusCode, res.Status)
	}

	return res, nil
}

func (p *Client) signiture(method, timestamp string, body io.Reader) string {
	// ## python
	// request_details = req.method + "\n"
	// request_details += url_obj.path + url_obj.params + "\n"
	// request_details += timestamp + "\n"
	// if req.body:
	// 	request_details += req.body + "\n"
	// else:
	// 	request_details += "\n"

	text := method + "\n"
	text += "/" + p.URL.Path + p.URL.Query().Encode() + "\n"
	text += timestamp + "\n"
	if body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		text += buf.String()
	}
	text += "\n"

	text += p.Key

	fmt.Printf("%+v\n", text)
	mac := makeHMAC(p.Secret, text)
	fmt.Printf("%+v\n", mac)

	return mac
}

func makeHMAC(key, str string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

// with close()
func decoder(res *http.Response, v interface{}) {
	json.NewDecoder(res.Body).Decode(&v)
	defer res.Body.Close()
}
