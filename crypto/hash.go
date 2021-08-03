package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
)

func Hash(text string) string{
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func HMac(text string) string{
	key := []byte(config.Get().API.SecretKey)
	mac := hmac.New(md5.New, key)
	mac.Write([]byte(text))

	return hex.EncodeToString(mac.Sum(nil))
}