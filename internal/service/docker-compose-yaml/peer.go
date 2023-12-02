package docker_compose_yaml

import (
	"errors"
	"log"
	"user/internal/service/common"
)

type PeerYaml struct {
}

type PeerCompose struct {
	OrgName string //组织名称
	Compose []*PeerComposeDetail
}

type PeerComposeDetail struct {
	IP            string //节点IP
	DNS           string //节点域名
	Image         string //镜像
	Certs         string //证书根目录
	ContainerName string //容器名称
	PeerAddress   string //peer节点地址
	DB            string //数据存储路径
	Ports         Ports  //节点部署时使用的端口
}

func (p *PeerYaml) Create(peerComps []PeerCompose) ([]File, error) {
	defer Logger("NodeYaml.Peer.Create")()

	if peerComps == nil {
		return nil, errors.New("参数不能为空")
	}

	var res []File
	var fileName string
	orderIndex := 1
	for _, orgComp := range peerComps {
		for _, compose := range orgComp.Compose {

			log.Println(orgComp.OrgName, "开始创建peer:")

			//容器名称
			log.Println(orgComp.OrgName, "容器名称校验")
			if compose.ContainerName == "" {
				log.Println("容器名称不能为空")
				break
			}

			//域名校验
			log.Println(compose.ContainerName, "域名校验")
			if compose.DNS == "" {
				log.Println(compose.ContainerName, "域名不能为空")
				break
			}

			//组织名称校验
			log.Println(compose.ContainerName, "组织名称校验")
			if orgComp.OrgName == "" {
				log.Println(compose.ContainerName, "组织名称不能为空")
				break
			}

			//镜像
			log.Println(compose.ContainerName, "镜像校验")
			if compose.Image == "" {
				log.Println(compose.ContainerName, "镜像不能为空")
				break
			}

			//端口赋值
			log.Println(compose.ContainerName, "端口赋值")

			//端口7051
			if compose.Ports.Peer7051 == "" { //不存在
				log.Println(compose.ContainerName, "端口传入为空，使用默认端口:7051")
				compose.Ports.Peer7051 = "7051"
			}

			//端口7052
			if compose.Ports.Peer7052 == "" { //不存在
				log.Println(compose.ContainerName, "端口传入为空，使用默认端口:7052")
				compose.Ports.Peer7052 = "7052"
			}

			tmp := common.ReplaceVar([]byte(peerTemplate), []common.YamlReplace{
				{
					Old: "${dns}",
					New: compose.DNS,
				},
				{
					Old: "${orgName}",
					New: orgComp.OrgName,
				},
				{
					Old: "${containerPeerName}",
					New: compose.ContainerName,
				},
				{
					Old: "${certs}",
					New: compose.Certs,
				},
				{
					Old: "${image}",
					New: compose.Image,
				},
				{
					Old: "${7051}",
					New: compose.Ports.Peer7051,
				},
				{
					Old: "${7052}",
					New: compose.Ports.Peer7052,
				},
			})

			fileName = compose.ContainerName + ".docker-compose-yaml"
			log.Println("生成文件名称", compose.ContainerName)

			res = append(res, File{
				Name: fileName,
				Data: tmp,
			})
		}
		orderIndex++
	}

	return res, nil
}
