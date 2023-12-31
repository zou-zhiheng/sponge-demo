// code generated by https://user

package config

import (
	"github.com/zhufuyi/sponge/pkg/conf"
)

func NewCenter(configFile string) (*Center, error) {
	nacosConf := &Center{}
	err := conf.Parse(configFile, nacosConf)
	return nacosConf, err
}

type Center struct {
	Nacos Nacos `docker-compose-yaml:"nacos" json:"nacos"`
}

type Nacos struct {
	ContextPath string `docker-compose-yaml:"contextPath" json:"contextPath"`
	DataID      string `docker-compose-yaml:"dataID" json:"dataID"`
	Format      string `docker-compose-yaml:"format" json:"format"`
	Group       string `docker-compose-yaml:"group" json:"group"`
	IPAddr      string `docker-compose-yaml:"ipAddr" json:"ipAddr"`
	NamespaceID string `docker-compose-yaml:"namespaceID" json:"namespaceID"`
	Port        int    `docker-compose-yaml:"port" json:"port"`
	Scheme      string `docker-compose-yaml:"scheme" json:"scheme"`
}
