package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/imroc/ontts/xf"

	"github.com/garyburd/redigo/redis"
	"github.com/imroc/log"
)

type Server struct {
	opts *Options
}

type Options struct {
	OutDir      string //音频输出目录
	Level       int    //音频生成速度级别，越快越耗CPU，级别1~10,数字越小速度越快
	TTSParams   string
	LoginParams string
	RedisAddr   string
	RedisPass   string
	Speed       int
}

type Speech struct {
	Id  string `json:"id"`
	Txt string `json:"txt"`
}

func New(opts *Options) *Server {
	return &Server{
		opts: opts,
	}
}

func (s *Server) Start() {
	if s.opts.LoginParams == "" {
		log.Error("no login params")
		return
	}
	var c redis.Conn
	var err error
	if s.opts.RedisPass == "" {
		c, err = redis.Dial("tcp", s.opts.RedisAddr)
	} else {
		c, err = redis.Dial("tcp", s.opts.RedisAddr, redis.DialPassword(s.opts.RedisPass))

	}
	if err != nil {
		log.Error("failed to connect redis:%v")
		return
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	err = psc.Subscribe("tts")
	if err != nil {
		log.Error("failed to subscribe:%v", err)
		return
	}

	sub, ok := psc.Receive().(redis.Subscription)
	if !ok {
		log.Error("first message is not subscription")
		return
	}
	if sub.Count == 0 {
		log.Error("redis subscription count is 0")
		return
	}

	err = setXF(s.opts.Speed, s.opts.TTSParams, s.opts.LoginParams)
	if err != nil {
		log.Error("failed to set xunfei params:%v", err)
		return
	}

	if s.opts.OutDir != "" && s.opts.OutDir[len(s.opts.OutDir)-1] != os.PathSeparator {
		s.opts.OutDir += string(os.PathSeparator)
	}

	var speech Speech
	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			err := json.Unmarshal(n.Data, &speech)
			if err != nil {
				log.Error("error unmarshal:%v", err)
				continue
			}
		TTS:
			err = xf.TextToSpeech(speech.Txt, s.opts.OutDir+speech.Id+".wav")
			if err != nil {
				log.Error("error convert:%v,tts ID:%s,TXT:%s", err, speech.Id, speech.Txt)
				time.Sleep(5 * time.Second)
				goto TTS
			}
			log.Debug("合成ID:%s,TXT:%s", speech.Id, speech.Txt)
		case error:
			log.Error("error redis message:%v", n)
			time.Sleep(10 * time.Second)
		default:
			log.Warn("unknown message:%v", n)
		}
	}

}

func (s *Server) Once(txt string, desPath string) error {
	log.Debug("tts:%s,login:%s", s.opts.TTSParams, s.opts.LoginParams)
	xf.SetTTSParams(s.opts.TTSParams)
	err := xf.Login(s.opts.LoginParams)
	if err != nil {
		return err
	}
	//不SetSleep,默认为0,单次合成以高性能模式
	log.Debug("txt:%s,output path:%s", txt, desPath)
	err = xf.TextToSpeech(txt, desPath)
	if err != nil {
		return err
	}
	return nil
}

func setXF(speedLevel int, ttsParams, loginParams string) error {
	if speedLevel < 1 || speedLevel > 10 {
		return fmt.Errorf("wrong speed level:%d,it should between 1 and 10", speedLevel)
	}

	sleepTime := 15000 * (speedLevel - 1)

	xf.SetSleep(sleepTime)

	xf.SetTTSParams(ttsParams)

	err := xf.Login(loginParams)
	if err != nil {
		return err
	}
	return nil
}
