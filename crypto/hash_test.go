package crypto

import (
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"testing"
	"time"
)

func TestHMac(t *testing.T) {
	config.LoadYAML("../settings.yaml")
	t.Log(HMac("HelloWorld"))
	t.Log(Hash("HelloWorld"))
}

func TestGenerateValidURL(t *testing.T) {
	config.LoadYAML("../settings.yaml")
	url := "/v1/processor/workers"
	ts := time.Now().Unix()
	url = fmt.Sprintf("%v?t=%d",url, ts)

	token := HMac(url)
	url = url + "&s=" + token
	t.Log(url)
}