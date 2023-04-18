/**
 * @Author Nil
 * @Description pkg/config/config.go
 * @Date 2023/3/28 01:06
 **/

package config

import (
	"fmt"
	"github.com/ha5ky/hu5ky-bot/pkg/util"
	"github.com/ha5ky/hu5ky-bot/pkg/watcher/filewatcher"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"sync"
)

var (
	SysCache = &SysConfig{}

	DefaultConfig = map[string]string{
		"test":    "manifests/profiles/config.yaml",
		"debug":   "manifests/profiles/config.yaml",
		"release": "manifests/profiles/config.yaml",
	}
)

type SysConfig struct {
	HttpConfig struct {
		Port               uint16 `yaml:"port"`
		ReadTimeout        int    `yaml:"readTimeout"`
		WriteTimeout       int    `yaml:"writeTimeout"`
		MaxHeaderBytes     int    `yaml:"maxHeaderBytes"`
		MaxMultipartMemory int64  `yaml:"maxMultipartMemory"`
	} `yaml:"httpConfig"`

	ServerConfig struct {
		LogLevel string `yaml:"logLevel"`
		Mode     string `yaml:"mode"`
		Storage  string `yaml:"storage"`
	} `yaml:"serverConfig"`

	GPT struct {
		OpenaiAPIKey string `yaml:"openaiApiKey"`
	} `yaml:"gpt"`

	DB struct {
		QdRant struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		} `yaml:"qdRant"`

		Mysql struct {
			User      string `yaml:"user"`
			Pwd       string `yaml:"pwd"`
			Host      string `yaml:"host"`
			Port      string `yaml:"port"`
			DBName    string `yaml:"dbName"`
			LogLevel  string `yaml:"logLevel"`
			Charset   string `yaml:"charset"`
			ParseTime string `yaml:"parseTime"`
			Loc       string `yaml:"loc"`
		} `yaml:"mysql"`
	} `yaml:"db"`
}

func yaml2Cache(mode string) {
	configFile, err := os.ReadFile(DefaultConfig[mode])
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, SysCache)
	if err != nil {
		panic(err)
	}
}

func Cache2yaml(conf *SysConfig) (err error) {
	var out []byte
	out, err = yaml.Marshal(conf)
	if err != nil {
		return
	}
	err = os.WriteFile(DefaultConfig[conf.ServerConfig.Mode], out, fs.ModePerm)
	if err != nil {
		return
	}
	return
}

var once sync.Once

func Watcher(mode string) {
	// at first time, we should trigger convert actively
	once.Do(func() {
		yaml2Cache(mode)
		fmt.Println("... init system cache ...")
		util.DumpPretty(*SysCache)
	})
	filewatcher.AddFileWatcher(filewatcher.NewWatcher(), DefaultConfig[mode], func() {
		yaml2Cache(mode)
		fmt.Println("system cache has been modified.")
		util.DumpPretty(*SysCache)
	})
}
