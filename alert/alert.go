package alert

import (
	"fmt"
	"github.com/guanaitong/go-common/runtime"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	URL = "http://message.frigate.devops.wuxingdev.cn"

	BY_CORPCODE_URL = URL + "/message/sendMsgsToUsersByQiyeWeChatId"

	BY_GROUP_URL = URL + "/message/sendMsgsToGroup"

	BY_APPNAME_URL = URL + "/message/sendMsgsToUserByAppName"

	bufferLen = 1000
)

var (
	appName   string
	msgHeader string
	workEnv   string
	ch        = make(chan map[string]interface{}, bufferLen)
)

func init() {
	appName = os.Getenv("APP_NAME")
	if appName == "" {
		appName = "unknown"
	}
	workEnv = os.Getenv("WORK_ENV")
	if os.Getenv("WORK_ENV") != "" {
		msgHeader = msgHeader + "env[" + os.Getenv("WORK_ENV") + "],"
	}

	msgHeader = msgHeader + "server_ip[" + getLocalIP() + "]"

	if os.Getenv("HOSTNAME") != "" {
		msgHeader = msgHeader + ",host[" + os.Getenv("HOSTNAME") + "]"
	}

	msgHeader = msgHeader + ".msg: "

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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("get local ip error:%s", err.Error())
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

var client = &http.Client{
	Timeout: time.Second * 5,
}

// 通过工号发送消息
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByCorpCodes(channel int, msg string, corpCodes ...string) {
	m := map[string]interface{}{
		"wechatIdList": strings.Join(corpCodes, ","),
		"way":          1,
	}
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	ch <- buildMsg(channel, msg, m)
}

// 通过组发送消息
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByGroupId(channel int, msg string, groupId int) {
	m := map[string]interface{}{
		"groupId": groupId,
		"way":     2,
	}
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	ch <- buildMsg(channel, msg, m)
}

// 通过应用名发送消息，自动获取应用名，不需要传递
// 0-没有渠道，1-企业微信，2-邮件，3-短信，4-企业微信+邮件，5-企业微信+短信，6-邮件+短信，7-企业微信+邮件+短信
func SendByAppName(channel int, msg string) {
	m := map[string]interface{}{
		"appName": appName,
		"way":     3,
	}
	if len(ch) >= bufferLen { //缓存区满，丢弃
		log.Printf("abort msg %s", msg)
		return
	}
	ch <- buildMsg(channel, msg, m)
}

func buildMsg(channel int, msg string, m map[string]interface{}) map[string]interface{} {
	m["msgContent"] = msgHeader + msg
	m["channel"] = channel
	m["time"] = time.Now().Unix() * 1000
	if workEnv != "" {
		m["workEnv"] = workEnv
	}
	return m
}

func send(data map[string]interface{}) error {

	way := data["way"].(int)

	byUrl := ""
	if way == 1 {
		byUrl = BY_CORPCODE_URL
	} else if way == 2 {
		byUrl = BY_GROUP_URL
	} else if way == 3 {
		byUrl = BY_APPNAME_URL
	}
	delete(data, "way")

	value := url.Values{}
	for k, v := range data {
		value.Add(k, fmt.Sprint(v))
	}
	req, err := http.NewRequest("POST", byUrl, strings.NewReader(value.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "GOLANG_UTIL")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
