ontts
==============
ontts 是语音在线合成服务

## 编译

##### GOPATH
项目源文件需要放在$GOPATH下

##### glide
通过glide管理依赖，若没有安装glide，需先安装

Ubuntu下glide安装方法:
``` sh
sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
sudo apt-get install glide
```
在源文件目录下执行以下命令下载依赖
``` sh
glide install
```

##### 编译
``` sh
go build
```

##### 运行
需要将libmsc.so加入环境变量
``` sh
mv xf/lib/libmsc.so /usr/local/lib/
vi ~/.bashrc
```
export LD_LIBRARY_PATH=/usr/local/lib


## 使用示例
##### 单次合成:
``` sh
./ontts -t "云喇叭快递，快递小管家，您的快递到了，请于下午6点前到学校后门申通快递取件" -o test.wav
```

##### 启动合成后台服务:
``` sh
./ontts -r ":6379" -d /tmp/out
```

## 命令参数
<pre>
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
       │   └── libmsc.so
       ├── README.md
       └── xf.go
</pre>
