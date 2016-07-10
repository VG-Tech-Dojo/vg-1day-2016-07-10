package bot

import (
	"fmt"
	"math/rand"
	"model"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type (
	// Processer
	// もととなるmessageをもとに，新規投稿するためのmessageを作成する
	Processer interface {
		Process(*model.Message) *model.Message
	}

	// EchoProcesser
	// 入力値のテキストそのままで返すProcesser
	EchoProcesser struct {
	}

	// Greetprocesser
	// 自身の名前で挨拶するProcesser
	GreetProcesser struct {
		Name string
	}

	UranaiProcesser struct {
	}

	// TimelineProcesser
	// homeのtimelineのtweetを1つ取得するProcesser
	TimelineProcesser struct {
		Api *anaconda.TwitterApi
	}

	ShiritoriProcesser struct {
		Dict map[string][]string
	}
)

func (p *EchoProcesser) Process(msgIn *model.Message) *model.Message {
	return &model.Message{Body: "[echo] " + msgIn.Body}
}

func (p *GreetProcesser) Process(msgIn *model.Message) *model.Message {
	txt := "[greet] nice to meet you! my name is " + p.Name
	return &model.Message{Body: txt}
}

func (p *TimelineProcesser) Init() {
	p.Api = getTwitterApi()
}

func (*UranaiProcesser) Process(msgIn *model.Message) *model.Message {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(3)
	text := ""
	switch random {
	case 0:
		text = "大吉"
	case 1:
		text = "吉"
	case 2:
		text = "凶"
	}
	return &model.Message{Body: text}
}

func (p *TimelineProcesser) Process(msgIn *model.Message) *model.Message {
	if p.Api == nil {
		return &model.Message{Body: "[timeline] api can not available"}
	}
	query := url.Values{}
	query.Add("count", "1")

	timeline, _ := p.Api.GetHomeTimeline(query)
	if len(timeline) < 1 {
		return &model.Message{Body: "[timeline] no result found"}
	}

	tweet := timeline[0]
	return &model.Message{Body: fmt.Sprintf("[timeline:%s] %s", tweet.User.Name, tweet.Text)}
}

func (p *ShiritoriProcesser) Init() {
	p.Dict = make(map[string][]string)
	p.Dict["あ"] = []string{
		"あんこ",
	}
}

func (p *ShiritoriProcesser) Process(msgIn *model.Message) *model.Message {
	s := strings.Fields(msgIn.Body)
	head := s[1][len(s[1])-3 : len(s[1])]
	if candidates, ok := p.Dict[head]; ok {
		ans := candidates[0]
		p.Dict[head] = append(candidates[1:])
		if len(p.Dict[head]) == 0 {
			delete(p.Dict, head)
		}
		return &model.Message{Body: ans, Username: "shiritoribot"}
	} else {
		return &model.Message{Body: "参りました", Username: "shiritoribot"}
	}
}
