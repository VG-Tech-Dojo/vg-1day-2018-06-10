package bot

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/syasin-5d/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/syasin-5d/model"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	chatbotAPIURLFormat = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
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

	// ChatBotProcessor はメッセージ本文に返答するprocessorの構造体です
	ChatBotProcessor struct{}
)

// Process は"hello, world!"というbodyがセットされたメッセージのポインタを返します
func (p *HelloWorldProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	return &model.Message{
		Body: msgIn.Body + ", world!",
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

// Talk はメッセージ本文からreplyを生成します
func (p *ChatBotProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Atalk (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)

	text := matchedStrings[1]

	params := url.Values{}
	params.Set("apikey", env.ChatBotAPIAppID)
	params.Add("query", text)
	
	res := &struct {
		Status int64 `json:status`
		Message string `json:message`
		Results []struct {
			Perplexity float64 `json:perplexity`
			Reply string `json:reply`
		} `json:results`
	}{}

	post(chatbotAPIURLFormat, params, res)

	if res.Status != 0 {
		return nil, fmt.Errorf("%#v", res)
	}

	return &model.Message{
		Body: res.Results[0].Reply,
	}, nil
}
