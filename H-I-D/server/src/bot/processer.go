package bot

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"model"
	"net/url"
	"math/rand"
	"time"
	"strings"
	"strconv"
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
	
	// UranaiProcesser
	//
	UranaiProcesser struct {
	}
	
	// WarikanProcesser
	//
	WarikanProcesser struct {
	}

	// TimelineProcesser
	// homeのtimelineのtweetを1つ取得するProcesser
	TimelineProcesser struct {
		Api *anaconda.TwitterApi
	}
)

func (p *EchoProcesser) Process(msgIn *model.Message) *model.Message {
	return &model.Message{Body: "[echo] " + msgIn.Body}
}

func (p *GreetProcesser) Process(msgIn *model.Message) *model.Message {
	txt := "[greet] nice to meet you! my name is " + p.Name
	return &model.Message{Body: txt}
}

func (p *UranaiProcesser) Process(msgIn *model.Message) *model.Message {
	rand.Seed(time.Now().UnixNano())
	result := "大吉"
	switch rand.Intn(4) {
	case 0:	result = "中吉"
	case 1:	result = "吉"
	case 2:	result = "末吉"
	case 3:	result = "凶"
	}
	txt := "result: " + result
	return &model.Message{Body: txt}
}

func (p *WarikanProcesser) Process(msgIn *model.Message) *model.Message {
	inputs := strings.Split(msgIn.Body, " ")
	sum, _ := strconv.Atoi(inputs[1])
	each, _ := strconv.Atoi(inputs[2])
	warikan := sum / each
	txt := "一人 " + strconv.Itoa(warikan) + "円です"
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
