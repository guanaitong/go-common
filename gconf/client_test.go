package gconf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type es struct {
	ClusterName string
	hosts       string
	Port        int32
}

func TestGetConfigCollection(t *testing.T) {
	c := GetConfigCollection("userdoor")
	v1 := c.GetConfig("es.properties")
	fmt.Println(v1)
	fmt.Println(c.AsMap())
	c1 := GetConfigCollection("userdoor")
	fmt.Println(c1.AsMap())

	v2 := c1.GetConfigAsStructuredMap("es.properties")
	fmt.Println(v2)
	p := new(es)
	c1.GetConfigAsBean("es.properties", p)
	fmt.Println(fmt.Sprint(*p))

	c2 := GetConfigCollection("jaeger")
	fmt.Println(c2.AsMap())
	time.Sleep(time.Second * 120)
}

type testBean struct {
	a string `config:"a"`
	b int
	c uint64
}

func TestReflectBase(t *testing.T) {

	p := new(testBean)
	v := reflect.ValueOf(p).Elem()

	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		fieldInfo1 := v.Field(i)
		fmt.Println(i)
		fmt.Println(fieldInfo.Name)
		fmt.Println(fieldInfo.Type)
		fmt.Println(fieldInfo.Tag)
		fmt.Println(fieldInfo1.Type().Name())

	}
}

func TestSplit(tt *testing.T) {
	t := "mysql_url=amy:1qazxsw@@tcp(mdb.servers.dev.ofc:3306)/notifyagent?charset=utf8"
	ss := strings.SplitN(t, "=", 2)
	for _, s := range ss {
		fmt.Println(s)
	}
}

func TestSplit1(tt *testing.T) {
	t := "mysql_url="
	ss := strings.SplitN(t, "=", 2)
	for _, s := range ss {
		fmt.Println(s)
	}
}
