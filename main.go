package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/imroc/ontts/server"

	"github.com/imroc/log"
)

var usageStr = `
Usage: ontts [options]

讯飞语音参数选项:
    -tp <param>                 TTS合成参数[有默认值]
    -lp <param>                 登录参数[有默认值]

单次合成模式选项:
    -t <text>                	待合成的文本
    -o <file>               	音频输出路径 

合成服务模式选项:
    -d <dir>                    音频保存的目录 
    -s <digit>                  合成速度级别(1-10),数值越小速度越快，越耗CPU[默认为1]
    -r <addr>                   redis连接地址
    -rp <pass>                  redis密码

日志选项:
    -l <file>                   日志输出路径[默认./ontts.log]
    -ll <level>                 日志输出级别(debug,info,warn,error)

其他:
    -h                          查看帮助 
`

func main() {
	opts := &server.Options{}

	var txt string
	var out string
	var help bool
	var logFile string
	var logLevel string

	flag.StringVar(&txt, "t", "", "单次合成的文本")
	flag.StringVar(&out, "o", "", "单次合成的输出路径")
	flag.StringVar(&logFile, "l", "ontts.log", "日志输出路径")
	flag.StringVar(&logLevel, "ll", "debug", "日志输出级别")
	flag.BoolVar(&help, "h", false, "Help")

	flag.StringVar(&opts.TTSParams, "tp", "voice_name = xiaoqi, text_encoding = UTF8, sample_rate = 8000, speed = 50, volume = 50, pitch = 50, rdn = 2", "TTS合成参数")
	flag.StringVar(&opts.LoginParams, "lp", "", "登录参数")
	flag.StringVar(&opts.RedisAddr, "r", ":6379", "redis连接地址")
	flag.StringVar(&opts.RedisPass, "rp", "", "redis连接密码")
	flag.StringVar(&opts.OutDir, "d", "", "音频输出目录")
	flag.IntVar(&opts.Speed, "s", 1, "合成速度")

	flag.Parse()

	if help {
		fmt.Printf("%s\n", usageStr)
		return
	}

	err := configureLog(logFile, logLevel)
	if err != nil {
		log.Error("日志配置失败:%v", err)
		return
	}

	s := server.New(opts)

	if txt != "" { // 单次合成
		if out == "" {
			out = txt + ".wav"
		}
		log.Debug("合成文本:%q,输出:%s", txt, out)
		if err = s.Once(txt, out); err != nil {
			log.Error("%v", err)
			return
		}
	}

	s.Start()

}

func configureLog(logFile, logLevel string) error {
	level := log.DEBUG

	switch strings.ToLower(logLevel) {
	case "debug":
		level = log.DEBUG
	case "info":
		level = log.INFO
	case "warn":
		level = log.WARN
	case "error":
		level = log.ERROR
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.Set(level, file, log.Lshortfile|log.LstdFlags)

	return nil
}
