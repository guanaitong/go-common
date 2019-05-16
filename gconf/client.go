package gconf

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	mathRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// gconf client就不适用glog了，因为glog需要flag parse，过多依赖不好。
var cache = map[string]*ConfigCollection{}
var RegionId int
var url string
var mux = new(sync.Mutex)
var clientId = uuid()

func init() {
	url = getUrl()
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")

	if len(host) > 0 || len(port) > 0 {
		log.Printf("app in k8s,use hc://gconf.kube-system/")
		url = "hc://gconf.kube-system/"
	}

	initRegionId()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("recover=>" + fmt.Sprint(r))
			}
		}()
		for {
			if len(cache) == 0 {
				time.Sleep(time.Second * 10)
				continue
			}
			var keys []string
			for k := range cache {
				keys = append(keys, k)
			}
			configAppIdList := strings.Join(keys, ",")
			needChangeAppIdList, err := httpGetListResp(url + "api/watch?configAppIdList=" + configAppIdList + "&clientId=" + clientId)
			if err != nil {
				log.Printf("wath error" + err.Error())
				time.Sleep(time.Second * 10)
				continue
			}

			for _, appId := range needChangeAppIdList {
				cache[appId].refreshData()
			}
		}
	}()
}

func uuid() (uuid string) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		uuid = fmt.Sprint(mathRand.Int63n(time.Now().UnixNano()))
		return
	}
	uuid = base64.URLEncoding.EncodeToString(b)
	return
}

func initRegionId() {
	resp, err := http.Get(url + "regionId")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	RegionId, err = strconv.Atoi(string(bs))

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Printf("reigonId initialized, value:%d", RegionId)
}

func getUrl() string {
	workEnv, workIdc := os.Getenv("WORK_ENV"), os.Getenv("WORK_IDC")
	log.Printf("workEnv %s workIdc %s", workEnv, workIdc)
	if workEnv == "dev" && workIdc == "ofc" {
		RegionId = 3
		return "hc://gconf.services.dev.ofc/"
	}
	if workEnv == "test" && workIdc == "jx" {
		RegionId = 8
		return "hc://gconf.services.test.jx/"
	}
	if workEnv == "product" && workIdc == "sh" {
		RegionId = 1
		return "hc://gconf.services.product.sh/"
	}
	if workEnv == "prepare" && workIdc == "sh" {
		RegionId = 11
		return "hc://gconf.services.product.sh/"
	}
	if workEnv == "product" && workIdc == "ali" {
		RegionId = 2
		return "hc://gconf.services.product.ali/"
	}
	RegionId = 3
	return "hc://gconf.services.dev.ofc/"
}

func GetConfigCollection(appId string) *ConfigCollection {
	res, ok := cache[appId]
	if ok {
		return res
	}

	mux.Lock()
	defer mux.Unlock()

	//double check
	res, ok = cache[appId]
	if ok {
		return res
	}

	configApp, err := httpGetMapResp(url + "api/getConfigApp?configAppId=" + appId)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	res = &ConfigCollection{
		appId:     appId,
		name:      configApp["name"],
		data:      map[string]*configData{},
		listeners: map[string][]ConfigChangeListener{},
	}
	res.refreshData()
	cache[appId] = res
	return res
}

func GetWorkEnvByRegionId(regionId int) string {
	if regionId == 3 {
		return "dev"
	}
	if regionId == 8 {
		return "test"
	}
	if regionId == 1 {
		return "product"
	}
	if regionId == 11 {
		return "prepare"
	}
	if regionId == 2 {
		return "product"
	}
	return "dev"
}

func GetWorkIdcByRegionId(regionId int) string {
	if regionId == 3 {
		return "ofc"
	}
	if regionId == 8 {
		return "jx"
	}
	if regionId == 1 {
		return "sh"
	}
	if regionId == 11 {
		return "sh"
	}
	if regionId == 2 {
		return "ali"
	}
	return "ofc"
}
