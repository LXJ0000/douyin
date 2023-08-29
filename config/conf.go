package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"strings"
)

type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	ParseTime bool `toml:"parse_time"`
	Loc       string
}

type Server struct {
	IP   string
	Port int
	URL  string
}

type Config struct {
	DB     Mysql `toml:"mysql"`
	Server `toml:"server"`
	Path   `toml:"path"`
}

// Path ffmpeg_path和static_path
type Path struct {
	FfmpegPath       string `toml:"ffmpeg_path"`
	StaticSourcePath string `toml:"static_source_path"`
}

var Info Config

// 包初始化加载时候会调用的函数
func init() {
	// toml加载配置文件xxx.toml
	if _, err := toml.DecodeFile("./config/config.toml", &Info); err != nil {
		panic(err)
	}
	//去除左右的空格
	strings.Trim(Info.Server.IP, " ")
	strings.Trim(Info.DB.Host, " ")
}

// DBConnectString 填充得到数据库连接字符串
func DBConnectString() string {
	arg := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Info.DB.Username, Info.DB.Password, Info.DB.Host, Info.DB.Port, Info.DB.Database,
		Info.DB.Charset, Info.DB.ParseTime, Info.DB.Loc)
	log.Println(arg)
	return arg
}
