package main

import (
	"fmt"
	"math/rand"
	"ontts/xf"
	"os"
	"testing"
	"time"

	"github.com/imroc/log"
)

func TestSpeed(t *testing.T) {
	file, _ := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	log.Set(log.DEBUG, file, log.LstdFlags)

	xf.SetTTSParams("voice_name = xiaoqi, text_encoding = UTF8, sample_rate = 8000, speed = 50, volume = 50, pitch = 50, rdn = 2")
	xf.SetSleep(0)

	err := xf.Login("appid = 5718a335, work_dir = .")

	if err != nil {
		log.Error("err:%v", err)
		return
	}

	now := time.Now()

	var txt string
	for i := 1; i < 1000; i++ {
		txt = getRandomString(r.Intn(21) + 40)
		err = xf.TextToSpeech(txt, fmt.Sprintf("/home/roc/wav/%d.wav", i))
		if err != nil {
			log.Error("err:%v", err)
			return
		}
		log.Debug("已生成第%d个：%s，共用%f秒", i, txt, time.Since(now).Seconds())
	}

	log.Info("执行完毕,总消耗：%f秒", time.Since(now).Seconds())
	xf.Logout()

}

var start rune = 0x4e00
var stop rune = 0x9fa5
var n int32 = int32(stop - start)
var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomString(l int) string {
	s := make([]rune, l)
	for i := 0; i < l; i++ {
		s[i] = rune(r.Int31n(n) + start)
	}
	return string(s)
}
