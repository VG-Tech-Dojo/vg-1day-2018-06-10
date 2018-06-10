package bot

import (
	"regexp"
	"strings"

	"fmt"

	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/ranbo/env"
	"github.com/VG-Tech-Dojo/vg-1day-2018-06-10/ranbo/model"
	"net/url"
)

const (
	keywordAPIURLFormat = "https://jlp.yahooapis.jp/KeyphraseService/V1/extract?appid=%s&sentence=%s&output=json"
	talkAPIURLFormat    = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
	youtubeAPIURLFormat = "https://www.googleapis.com/youtube/v3/search?part=id&type=video&key=%s&q=%s"
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

	TalkProcessor struct{}

	GachaProcessor struct{}

	// youtube bot
	YoutubeProcessor struct{}
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

// youtube bot
func (p *YoutubeProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	r := regexp.MustCompile("\\Ayoutube (.+)")
	matchedStrings := r.FindStringSubmatch(msgIn.Body)
	if len(matchedStrings) != 2 {
		return nil, fmt.Errorf("bad message: %s", msgIn.Body)
	}

	text := matchedStrings[1]

	requestURL := fmt.Sprintf(youtubeAPIURLFormat, env.YoutubeAPIAppID, url.QueryEscape(text))

	type youtubeAPIResponse map[string]interface{}

	res := &struct {
		Items []struct {
			Id struct {
				VideoId string `json:videoid`
			} `json:id`
		} `json:items`
	}{}
	get(requestURL, res)

	fmt.Println(res.Items)

	// set videos this
	var maxlen int = len(res.Items)
	if maxlen > 3 {
		maxlen = 3
	}
	
	var videos string
	for i := 0; i < maxlen; i++ {
		if i == 0 {
			videos += "<iframe width=\"560\" height=\"315\" src=\"https://www.youtube.com/embed/" + res.Items[i].Id.VideoId + "\" frameborder=\"0\" allow=\"autoplay; encrypted-media\" allowfullscreen></iframe><br>"
			continue
		}
		videos += "<a href=\"https://www.youtube.com/watch?v=" + res.Items[i].Id.VideoId + "\" target=\"_blank\">" + res.Items[i].Id.VideoId + "</a><br>"
	}

	return &model.Message{
		Body: videos,
	}, nil
}

func (p *GachaProcessor) Process(msgIn *model.Message) (*model.Message, error) {
	fortunes := []string{
		"SSR",
		"SR",
		"R",
		"N",
	}
	result := fortunes[randIntn(len(fortunes))]
	return &model.Message{
		Body: result,
	}, nil
}

// func (p *KeywordProcessor) Talk(msgIn *model.Message) (*model.Message, error) {
// 	r := regexp.MustCompile("\\Atalk (.+)")
// 	matchedStrings := r.FindStringSubmatch(msgIn.Body)
// 	if len(matchedStrings) != 2 {
// 		return nil, fmt.Errorf("bad message: %s", msgIn.Body)
// 	}
//
// 	text := matchedStrings[1]
//
// 	// requestURL := fmt.Sprintf(keywordAPIURLFormat, env.TalkAPIAppID, url.QueryEscape(text), "おはよう")
//
// 	type keywordAPIResponse map[string]interface{}
// 	var response keywordAPIResponse
// 	val := url.Values{}
// 	val.Set("apikey", env.KeywordAPIAppID)
// 	val.Add("query", "おはよう")
//
// 	post(talkAPIURLFormat, val, &response)
// 	talks := make([]string, 0, len(response))
// 	for k, v := range response {
// 		if k == "Error" {
// 			return nil, fmt.Errorf("%#v", v)
// 		}
// 		talks = append(talks, k)
// 	}
//
// 	return &model.Message{
// 		Body: "キーワード：" + strings.Join(talks, ", "),
// 	}, nil
// }
