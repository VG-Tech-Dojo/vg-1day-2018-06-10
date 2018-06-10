package bot

import (
	"regexp"
	"strings"

	"fmt"

	"net/url"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/nakata/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/nakata/model"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	chatAPIURLFormat    = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
)

type BitflyerTicker struct {
	Symbol string     `json:"product_code"`
	Price float64   `json:"best_ask"`
}

type CoincheckTicker struct {
	Price float64   `json:"ask"`
}

type ZaifTicker struct {
	Price float64   `json:"last_price"`
}

type (
	// Processor はmessageを受け取り、投稿用messageを作るインターフェースです
	Processor interface {
		Process(message *model.Message) (*model.Message, error)
	}

	// HelloWorldProcessor は"hello, world!"メッセージを作るprocessorの構造体です
	HelloWorldProcessor struct{}

	// OmikujiProcessor は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかをランダムで作るprocessorの構造体です
	OmikujiProcessor struct{}

	// KeywordProcessor はメッセージ本文からキーワードを抽出するprocessorの構造体です
	KeywordProcessor struct{}

	GachaProcessor struct{}

	ChatProcessor struct{}

	BtcProcessor struct{}

	SpreadProcessor struct {}
)

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	return &model.Message{
		Body: msgIn.Body + ", world!",
	}, nil
}

// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	ranks := []string{
		"SSレア",
		"Sレア",
		"レア",
		"ノーマル",
	}
	result := ranks[randIntn(len(ranks))]
	return &model.Message{
		Body: result,
	}, nil
}

// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *OmikujiProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	fortunes := []string{
		"大吉",
		"吉",
		"中吉",
		"小吉",
		"末吉",
		"凶",
	}
	result := fortunes[randIntn(len(fortunes))]
	return &model.Message{
		Body: result,
	}, nil
}

// Process はメッセージ本文からキーワードを抽出します
func (p *KeywordProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Akeyword (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	if len(matchedStrings) != 2 {
		return nil, fmt.Errorf("bad message: %s", msgIn.Body)
	}

	text := matchedStrings[1]

	requestURL := fmt.Sprintf(keywordAPIURLFormat, env.KeywordAPIAppID, url.QueryEscape(text))

	type keywordAPIResponse map[string]interface{}
	var response keywordAPIResponse
	get(requestURL, &response)

	keywords := make([]string, 0, len(response))
	for k, v := range response {
		if k == "Error" {
			return nil, fmt.Errorf("%#v", v)
		}
		keywords = append(keywords, k)
	}

	return &model.Message{
		Body: "キーワード：" + strings.Join(keywords, ", "),
	}, nil
}

// Process はメッセージ本文からキーワードを抽出します
func (p *ChatProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	if len(matchedStrings) != 2 {
		return nil, fmt.Errorf("bad message: %s", msgIn.Body)
	}

	text := matchedStrings[1]

	response := &struct {
		Status  int64  `json:status`
		Message string `json:message`
		Results []struct {
			Perplexity float64 `json:perplexity`
			Reply      string  `json:reply`
		} `json:results`
	}{}
	values := url.Values{}
	values.Add("apikey", env.ChatAPIKey)
	values.Add("query", text)
	err := post(chatAPIURLFormat, values, response)

	if err != nil {
		return nil, fmt.Errorf("%#v", err)
	}

	return &model.Message{
		Body: response.Results[0].Reply,
	}, nil
}

// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *BtcProcessor) Process(msgIn *model.Message) (*model.Message, error) {

	r := regexp.MustCompile("\\Abtc (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	if len(matchedStrings) != 2 {
		return nil, fmt.Errorf("bad message: %s", msgIn.Body)
	}

	exchange := matchedStrings[1]
	var str string

	// bitflyer の振る舞い
	if exchange == "bitflyer"{
		url := "https://api.bitflyer.jp/v1/ticker"
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)

		var ticker string = string(b)
		jsonBytes := ([]byte)(ticker)
		data := new(BitflyerTicker)

		if err := json.Unmarshal(jsonBytes, data); err != nil {
			return nil, fmt.Errorf("%#v", err)
		}

		str = fmt.Sprintf("bitflyerの%s価格は%f円", data.Symbol, data.Price)
	}

	// coincheck の振る舞い
	if exchange == "coincheck"{
		url := "https://coincheck.com/api/ticker"
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)

		var ticker string = string(b)
		jsonBytes := ([]byte)(ticker)
		data := new(CoincheckTicker)

		if err := json.Unmarshal(jsonBytes, data); err != nil {
			return nil, fmt.Errorf("%#v", err)
		}

		str = fmt.Sprintf("coincheckのBTC価格は%f円", data.Price)
	}

	// zaif の振る舞い
	if exchange == "zaif"{
		url := "https://api.zaif.jp/api/1/last_price/btc_jpy"
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)

		var ticker string = string(b)
		jsonBytes := ([]byte)(ticker)
		data := new(ZaifTicker)

		if err := json.Unmarshal(jsonBytes, data); err != nil {
			return nil, fmt.Errorf("%#v", err)
		}

		str = fmt.Sprintf("zaifのBTC価格は%f円", data.Price)
	}

	return &model.Message{
		Body: str,
	}, nil
}



func getPrice(exchange string) float64{
	var price float64

	// bitflyer の振る舞い
	if exchange == "bitflyer"{
		url := "https://api.bitflyer.jp/v1/ticker"
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)

		var ticker string = string(b)
		jsonBytes := ([]byte)(ticker)
		data := new(BitflyerTicker)

		if err := json.Unmarshal(jsonBytes, data); err != nil {
			return 0
		}

		price = data.Price
	}

	// coincheck の振る舞い
	if exchange == "coincheck"{
		url := "https://coincheck.com/api/ticker"
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)

		var ticker string = string(b)
		jsonBytes := ([]byte)(ticker)
		data := new(CoincheckTicker)

		if err := json.Unmarshal(jsonBytes, data); err != nil {
			return 0
		}

		price = data.Price
	}

	// zaif の振る舞い
	if exchange == "zaif"{
		url := "https://api.zaif.jp/api/1/last_price/btc_jpy"
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)

		var ticker string = string(b)
		jsonBytes := ([]byte)(ticker)
		data := new(ZaifTicker)

		if err := json.Unmarshal(jsonBytes, data); err != nil {
			return 0
		}

		price = data.Price
	}
	return price
}


// Process は"大吉", "吉", "中吉", "小吉", "末吉", "凶"のいずれかがbodyにセットされたメッセージへのポインタを返します
func (p *SpreadProcessor) Process(msgIn *model.Message) (*model.Message, error) {

	r := regexp.MustCompile("\\Aspread (.+) (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)


	exchange1 := matchedStrings[1]
	exchange2 := matchedStrings[2]

	price1 := getPrice(exchange1)
	price2 := getPrice(exchange2)

	spread := price1 - price2

	str := fmt.Sprintf("%s, %sのBTC価格差は%f円", exchange1, exchange2, spread)

	return &model.Message{
		Body: str,
	}, nil
}

