package docker_compose_yaml

import (
	"errors"
	"log"
	"user/internal/service/common"
)

type OrderYaml struct {
}

type OrderCompose struct {
	ContainerName string //节点名称
	OrgName       string //组织名称
	DNS           string //域名
	Image         string //镜像
	Certs         string //证书根目录
	Ports         Ports  //使用的端口
}

func (o *OrderYaml) Create(orders []OrderCompose) ([]File, error) {
	defer Logger("NodeYaml.Order.Create")()

	if len(orders) == 0 {
		return nil, errors.New("参数不能为空")
	}

	var res []File
	var fileName string

	for _, order := range orders {

		log.Println("开始创建order:", order.DNS)

		//容器名称
		log.Println("节点名称校验")
		if order.ContainerName == "" {
			log.Println("order节点名称不能为空")
			break
		}

		//域名
		log.Println(order.ContainerName, "域名校验")
		if order.DNS == "" {
			log.Println("域名不能为空")
			break
		}

		//组织名称
		log.Println(order.ContainerName, "组织名称校验")
		if order.OrgName == "" {
			log.Println(order.DNS, "组织名称不能为空")
			break
		}

		//证书根目录
		log.Println(order.ContainerName, "证书根目录校验")
		if order.Certs == "" {
			log.Println(order.ContainerName, "证书根目录不能为空")
			break
		}

		//镜像
		log.Println(order.ContainerName, "镜像校验")
		if order.Image == "" {
			log.Println(order.ContainerName, "镜像不能为空")
			break
		}

		//端口赋值
		log.Println(order.ContainerName, "端口赋值")
		if order.Ports.Order7053 == "" { //使用默认端口
			log.Println(order.ContainerName, "端口传入为空，使用默认端口:7053")
			order.Ports.Order7053 = "7053"
		}

		log.Println(order.ContainerName, "端口赋值")
		if order.Ports.Order7050 == "" { //使用默认端口
			log.Println(order.ContainerName, "端口传入为空，使用默认端口:7050")
			order.Ports.Order7050 = "7050"
		}

		tmp := common.ReplaceVar([]byte(orderTemplate), []common.YamlReplace{
			{
				Old: "${dns}",
				New: order.DNS,
			},
			{
				Old: "${orgName}",
				New: order.OrgName,
			},
			{
				Old: "${image}",
				New: order.Image,
			},
			{
				Old: "${certs}",
				New: order.Certs,
			},
			{
				Old: "${containerOrderName}",
				New: order.ContainerName,
			},
			{
				Old: "${7050}",
				New: order.Ports.Order7050,
			},
			{
				Old: "${7053}",
				New: order.Ports.Order7053,
			},
		})

		fileName = order.ContainerName + ".docker-compose-yaml"
		log.Println("生成文件名称", fileName)

		res = append(res, File{
			Name: fileName,
			Data: tmp,
		})
	}

	return res, nil
}
