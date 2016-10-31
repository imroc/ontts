ontts
==============
ontts 是go写的语音合成服务(文本转语音)，原理是cgo调用的科大讯飞的在线语音合成linux SDK(用go封装SDK)。支持单次合成与后台合成服务两种模式，后台合成服务是订阅redis中的数据(业务系统可将待合成文本发布到redis)，一有数据立马合成并存储到磁盘

## 安装

``` sh
go get github.com/imroc/ontts
```

## 运行
需要将libmsc.so加入环境变量
``` sh
cp xf/lib/x64/libmsc.so /usr/local/lib/
vi ~/.bashrc
```
export LD_LIBRARY_PATH = /usr/local/lib


## 使用示例
##### 单次合成:
``` sh
./ontts -t "云喇叭快递，快递小管家，您的快递到了，请于下午6点前到学校后门申通快递取件" -o test.wav -lp "appid = 5808ae7e, work_dir = ."
```

##### 启动合成后台服务:
``` sh
./ontts -r ":6379" -d /tmp/out -lp "appid = 5808ae7e, work_dir = ."

```

## redis数据
后台合成的服务是订阅的redis中"tts"的通道中的数据,redis发布数据示例：
```sh
redis-cli>publish tts "{\"id\":\"245671051\",\"txt\":\"这是一段测试语音\"}"
```
注：生成的语音文件名是id加".wav"后缀

## 命令参数
<pre>
讯飞语音参数选项:
    -tp <param>                 TTS合成参数[有默认值]
    -lp <param>                 登录参数

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
</pre>

## 目录
<pre>
── ontts
   ├── glide.yaml (glide依赖配置)
   ├── main.go (程序入口)
   ├── README.md
   ├── server (TTS合成主体逻辑的package)
   │   └── server.go
   ├── speed_test.go (速度测试)
   └── xf (讯飞SDK的Go封装)
       ├── doc (讯飞语音linux SDK相关参考)
       ├── include (cgo需要用到的头文件)
       ├── lib (动态链接库 SDK)
       ├── README.md
       └── xf.go
</pre>
