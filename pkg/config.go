package pkg

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ArangoDBInfo ArangoDB 相关信息
type ArangoDBInfo struct {
	Url             string            `yaml:"url"`
	Username        string            `yaml:"username"`
	Passwd          string            `yaml:"passwd"`
	Database        string            `yaml:"database"`
	ModelColNameMap map[string]string `yaml:"collections"` // key:各模块名（模块名定义在 pkg/GLOBAL.go 中），value:数据库中的集合名
}

// Config 配置信息
type Config struct {
	Addr     string       `yaml:"addr"`
	ArangoDB ArangoDBInfo `yaml:"arangoDB"`
	LogPath  string       `yaml:"logPath"`
}

// NewConfig 构造 Config ，读取配置文件，获取配置数据
func NewConfig(filepath string) *Config {
	// 解析 conf.yaml 文件
	inFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	conf := new(Config)
	err = yaml.Unmarshal(inFile, conf)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Print(fmt.Sprintf("%+v", conf))
	return conf
}
