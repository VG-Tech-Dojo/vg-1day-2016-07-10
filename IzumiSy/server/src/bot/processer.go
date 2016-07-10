package bot

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"math/rand"
	"model"
	"net/url"
	"strconv"
	"strings"
	"time"
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

	// GreetProcesser
	// 自身の名前で挨拶するProcesser
	GreetProcesser struct {
		Name string
	}

	// UranaiProcessor
	// 占い用のProcessor
	UranaiProcessor struct {
	}

	// WarikanProecssor
	// 割り勘用のProcessor
	WarikanProcessor struct {
	}

	// TimelineProcesser
	// homeのtimelineのtweetを1つ取得するProcesser
	TimelineProcesser struct {
		Api *anaconda.TwitterApi
	}
)

func (p *WarikanProcessor) Process(msgIn *model.Message) *model.Message {
	_params := strings.Split(msgIn.Body, " ")
	var _result string

	if len(_params) != 3 {
		return &model.Message{Body: "割り勘はできなそうですね"}
	}

	_baseAmount, _ := strconv.Atoi(_params[1])
	_division, _ := strconv.Atoi(_params[2])

	if _baseAmount <= 0 || _division <= 0 {
		return &model.Message{Body: "割り勘できません"}
	}

	_result = fmt.Sprintf("一人%d円です", (_baseAmount / _division))

	return &model.Message{Body: _result}
}

func (p *UranaiProcessor) Process(msgIn *model.Message) *model.Message {
	rand.Seed(time.Now().UnixNano())

	var _result string
	switch rand.Intn(2) {
	case 0:
		_result = "大吉"
	case 1:
		_result = "吉"
	case 2:
	default:
		_result = "凶"
	}

	return &model.Message{Body: _result}
}

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
