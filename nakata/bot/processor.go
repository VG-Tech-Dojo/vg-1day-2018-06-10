package bot

import (
	"regexp"
	"strings"

	"fmt"

	"net/url"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/nakata/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/nakata/model"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	chatAPIURLFormat    = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
	hatenaAPIURLFormat  = "http://b.hatena.ne.jp/entry/jsonlite/?url=%s"
)

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

	HatenaProcessor struct{}
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

// Process はメッセージ本文からキーワードを抽出します
func (p *HatenaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Ahatena (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	if len(matchedStrings) != 2 {
		return nil, fmt.Errorf("bad message: %s", msgIn.Body)
	}

	in_url := matchedStrings[1]
	requestURL := fmt.Sprintf(hatenaAPIURLFormat, url.QueryEscape(in_url))

	response := &struct {
		Title     string `json:status`
		Bookmarks []struct {
			Comment string `json:comment`
		} `json:bookmarks`
	}{}
	err := get(requestURL, &response)

	fmt.Println(requestURL)
	fmt.Println(response)
	if err != nil {
		return nil, fmt.Errorf("%#v", err)
	}
	var out string
	for _, v := range response.Bookmarks {
		if v.Comment != "" {
			out = v.Comment
			break
		}
	}
	return &model.Message{
		Body: out,
	}, nil
}
