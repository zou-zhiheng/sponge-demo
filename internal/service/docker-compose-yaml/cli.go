package docker_compose_yaml

import (
	"errors"
	"log"
	"user/internal/service/common"
)

type CliYaml struct {
}

type CliCompose struct {
	OrgName string //组织名称
	Compose []*CliComposeDetail
}

type CliComposeDetail struct {
	IP            string //节点IP
	DNS           string //节点域名
	Image         string //镜像
	Certs         string //证书根目录
	ContainerName string //容器名称
	PeerAddress   string //peer节点地址
	DB            string //数据存储路径
	Ports         Ports  //节点部署时使用的端口
}

func (c *CliYaml) Create(orgKeyMap []CliCompose) ([]File, error) {
	defer Logger("NodeYaml.Cli.Create")()

	if orgKeyMap == nil {
		return nil, errors.New("参数不能为空")
	}

	var res []File
	var fileName string
	orderIndex := 1
	for key, orgComp := range orgKeyMap {
		for _, compose := range orgComp.Compose {

			log.Println(key, "开始创建cli:")

			//容器名称
			log.Println(orgComp.OrgName, "容器名称校验")
			if compose.ContainerName == "" {
				log.Println(key, compose.DNS, "容器名称不允许为空")
				break
			}

			//域名校验
			log.Println(compose.ContainerName, "域名校验")
			if compose.DNS == "" {
				log.Println(key, "域名不能为空")
				break
			}

			//组织名称校验
			log.Println(compose.ContainerName, "组织名称校验")
			if orgComp.OrgName == "" {
				log.Println(orgComp.OrgName, compose.DNS, "组织名称不能为空")
				break
			}

			//节点地址
			log.Println(compose.ContainerName, "Peer节点地址赋值")
			if compose.PeerAddress == "" {
				log.Println(key, compose.DNS, "Peer节点地址不允许为空")
				break
			}

			//镜像
			log.Println(compose.ContainerName, "镜像校验")
			if compose.Image == "" {
				log.Println(compose.ContainerName, "镜像不能为空")
				break
			}

			//yaml文件生成
			tmp := common.ReplaceVar([]byte(cliTemplate), []common.YamlReplace{
				{
					Old: "${dns}",
					New: compose.DNS,
				},
				{
					Old: "${orgName}",
					New: orgComp.OrgName,
				},
				{
					Old: "${containerName}",
					New: compose.ContainerName,
				},
				{
					Old: "${orgName}",
					New: key,
				},
				{
					Old: "${dns}",
					New: compose.DNS,
				},
				{
					Old: "${image}",
					New: compose.Image,
				},
			})

			//拼接文件名称
			fileName = compose.ContainerName + ".docker-compose-yaml"
			log.Println("生成文件名称", fileName)

			res = append(res, File{
				Name: fileName,
				Data: tmp,
			})
		}
		orderIndex++
	}

	return res, nil
}
