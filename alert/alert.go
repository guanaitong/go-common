package alert

import (
	"bytes"
	"encoding/json"
	"github.com/guanaitong/go-common/runtime"
	"github.com/guanaitong/go-common/system"
	"github.com/guanaitong/go-common/tuple"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	URL = "http://message.frigate.devops.wuxingdev.cn"

	// way = 0
	ByAppNameUrl = URL + "/v2/message/sendMsgByAppNames"

	// way = 1
	ByGroupUrl = URL + "/v2/message/sendMsgByGroups"

	// way = 2
	ByQiWeiXinUrl = URL + "/v2/message/sendMsgByWeChatIds"

	bufferLen = 4096
)

var (
	ch = make(chan *FrigateMessage, bufferLen)
)

func init() {
	go func() {
		for {
			func() {
				defer runtime.HandleCrashWithConfig(false)
				for data := range ch {
					e := send(data)
					if e != nil {
						log.Printf("send fail,data:%v, error:%s", data, e.Error())
					}
				}

			}()
		}
	}()
}

var client = &http.Client{
	Timeout: time.Second * 5,
}

type FrigateMessage struct {
	//发送渠道，默认通过1为企业微信通知
	Channel int `json:"channel"`
	//消息标题
	Title string `json:"title"`
	//消息内容
	Content string `json:"content"`
	//当有异常堆栈时，堆栈内容
	Stack string `json:"stack"`
	//模块
	Module string `json:"module"`
	//标签
	Tags map[string]string `json:"tags"`

	// ------------------以下属于系统变量------------------------
	traceId     string `json:"traceId"`
	hostIp      string `json:"hostIp"`
	appName     string `json:"appName"`
	appInstance string `json:"appInstance"`
	workEnv     string `json:"workEnv"`
	workIdc     string `json:"workIdc"`
	//发送时间
	time   int64 `json:"time"`
	format bool  `json:"format"`

	receiveInfo tuple.Pair
	way         int8
}

func NewMessage() *FrigateMessage {
	return &FrigateMessage{
		Title:       "frigate 消息通知",
		appName:     system.GetAppName(),
		appInstance: system.GetAppInstance(),
		hostIp:      system.GetHostIp(),
		workEnv:     system.GetWorkEnv(),
		workIdc:     system.GetWorkIdc(),
		time:        time.Now().Unix() * 1000,
	}
}

// 通过工号发送消息
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByCorpCodes(channel int, msg string, corpCodes ...string) {
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	m := NewMessage()
	m.receiveInfo = tuple.Pair{
		Key:   "receiveWeChatIds",
		Value: strings.Join(corpCodes, ","),
	}
	m.Channel = channel
	m.Content = msg
	m.way = 2
	ch <- m
}

// 通过组发送消息
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByGroupId(channel int, msg string, groupId int) {
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	m := NewMessage()
	m.receiveInfo = tuple.Pair{
		Key:   "receiveGroups",
		Value: strconv.Itoa(groupId),
	}
	m.Channel = channel
	m.Content = msg
	m.way = 1
	ch <- m
}

// 通过应用名发送消息，自动获取应用名，不需要传递
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByAppName(channel int, msg string) {
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	m := NewMessage()
	m.receiveInfo = tuple.Pair{
		Key:   "receiveAppNames",
		Value: system.GetAppName(),
	}
	m.Channel = channel
	m.Content = msg
	m.way = 0
	ch <- m
}

func send(message *FrigateMessage) error {
	way := message.way

	byUrl := ""
	if way == 0 {
		byUrl = ByAppNameUrl
	} else if way == 1 {
		byUrl = ByGroupUrl
	} else if way == 2 {
		byUrl = ByQiWeiXinUrl
	}

	data, _ := json.Marshal(message)

	req, err := http.NewRequest("POST", byUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "GOLANG_UTIL")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("send request fail, error:%s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	str, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("send request fail, resp_status:%s,body:%s", resp.Status, string(str))
	return nil
}
