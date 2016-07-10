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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"github.com/jmoiron/jsonq"
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

	// ImageProcesser
	//
	ImageProcesser struct {
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

func (p *ImageProcesser) Process(msgIn *model.Message) *model.Message {
	inputs := strings.Split(msgIn.Body, " ")
	query := inputs[1]
	imageurl := Get(query);
	txt := imageurl + query
	return &model.Message{Body: txt}
}


/* imagebot
-------------------*/
func getAPIKey() (apikey string, err error) {
	apikey = os.Getenv("BING_API_KEY")
	if apikey == "" {
		return "", errors.New("Not Set APIKEY")
	}
	return
}

func getImageType(contentType string) (imageType string, err error) {
	switch contentType {
	case "image/jpeg":
		imageType = "jpeg"
	case "image/png":
		imageType = "png"
	case "image/gif":
		imageType = "gif"
	default:
		return "", errors.New("Unknown ContentType")
	}

	return
}

func getJSON(query string) (json string, err error) {
	client := &http.Client{}
	URL := "https://api.datamarket.azure.com/Bing/Search/Image"
	apikey, err := getAPIKey()
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Add("Query", query)
	values.Add("$format", "json")
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	req.URL.RawQuery = values.Encode()
	req.SetBasicAuth(apikey, apikey)
	response, _ := client.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json = string(body)

	return
}

func parseJSON(jsonStr string) (urls [][]string, err error) {

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	for i := 0; i < 1; i++ {
		url, _ := jq.String("d", "results", strconv.Itoa(i), "MediaUrl")
		contentType, _ := jq.String("d", "results", strconv.Itoa(i), "ContentType")
		imageType, err := getImageType(contentType)
		if err != nil {
			return nil, err
		}

		urls = append(urls, []string{url, imageType})
	}
	return
}

func Get(query string) string {

	jsonStr, err := getJSON(query)
	if err != nil {
		panic(err)
	}
	urls, err := parseJSON(jsonStr)
	if err != nil {
		panic(err)
	}

	timeStamp := time.Now().Format("20060102150405")

	dirName := "mimorin-" + timeStamp
	if err := os.Mkdir(dirName, 0777); err != nil {
		panic(err)
	}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	statusChan := make(chan string)
	for idx, url := range urls {
		filePath := dirName + "/" + "mimorin" + strconv.Itoa(idx) + "." + url[1]
		go func(url, filePath string) {
			statusChan <- (url)
		}(url[0], filePath)
	}
	return <-statusChan
}
