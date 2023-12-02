package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/robfig/cron/v3"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
	"testing"
)

var namingClient naming_client.INamingClient
var configClient config_client.IConfigClient

func init() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "6d617f7a-fcd2-4924-a541-aae176275dca", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "192.168.182.136",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	var err error
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

}

func TestNacosRegister(t *testing.T) {
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "10.0.0.12",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(success)
	select {}

}

func TestNacosSubscribe(t *testing.T) {

	err := namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: "demo.go",
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Println(services[0], err)
		},
	})

	if err != nil {
		panic(err)
	}

	select {}

}

func TestNacosDeleteInstance(t *testing.T) {

	success, err := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          "10.0.0.11",
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true,
		Cluster:     "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(success)
}

func TestNacosPubConfig(t *testing.T) {

	success, err := configClient.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "DEFAULT_GROUP",
		Content: "hello world!22222",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(success)
}

func TestNacosDeleteConfig(t *testing.T) {

	success, err := configClient.DeleteConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(success)

}

func TestNacosListenConfig(t *testing.T) {

	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("listening")
	select {}
}

func TestNacosRegisterTry(t *testing.T) {

	nacosCli := nacos.New(namingClient)
	watcher, err := nacosCli.Watch(context.Background(), "demo.go")
	if err != nil {
		panic(err)
	}
	instances, err := watcher.Next()
	if err != nil {
		panic(err)
	}
	fmt.Println(instances[0].ID)

}

func NacosPingAndUpdate() {
	instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: "demo.go",
		HealthyOnly: true,
	})
	if err != nil {
		fmt.Println("NacosPingAndUpdate:", err.Error())
		return
	}

	for i := range instances {
		fmt.Println(instances[i])
	}

}

func TestTicker(t *testing.T) {
	c := cron.New()
	// 每5分钟执行一次任务
	_, err := c.AddFunc("*/5 * * * *", NacosPingAndUpdate)
	if err != nil {
		fmt.Println("添加定时任务失败：", err)
		return
	}

	// 启动定时任务
	c.Start()

	// 程序运行60秒后退出，可以根据需要调整
	select {}

	// 停止定时任务
	//c.Stop()
}

func TestDemo(t *testing.T) {
	//NacosPingAndUpdate()
	TestNacosRegister(&testing.T{})
}
