package docker_compose_yaml

import (
	"errors"
	"log"
	"user/internal/service/common"
)

type CouchDBYaml struct {
}

type CouchDBCompose struct {
	ContainerName string //容器名称
	Image         string //镜像
	DB            string //数据存储路径
	Ports         Ports  //节点部署时使用的端口
}

func (c *CouchDBYaml) Create(couchDBs []CouchDBCompose) ([]File, error) {
	defer Logger("NodeYaml.CouchDB.Create")()

	if len(couchDBs) == 0 {
		return nil, errors.New("参数不能空")
	}

	var res []File
	var fileName string
	for _, couchDB := range couchDBs {
		log.Println("开始创建couchDB:")
		//节点名称校验
		log.Println("节点名称校验")
		if couchDB.ContainerName == "" {
			log.Println("节点名称不能为空")
			break
		}

		//数据存储校验
		log.Println("数据存储校验")
		if couchDB.DB == "" {
			log.Println("数据存储路径不能为空")
			break
		}

		//镜像
		log.Println(couchDB.ContainerName, "镜像校验")
		if couchDB.Image == "" {
			log.Println(couchDB.ContainerName, "镜像不能为空")
			break
		}

		//端口赋值
		log.Println(couchDB.ContainerName, "端口赋值")
		if couchDB.Ports.Couchdb5984 == "" { //不存在
			log.Println(couchDB.ContainerName, "端口传入为空，使用默认端口")
			couchDB.Ports.Couchdb5984 = "5984"
		}

		tmp := common.ReplaceVar([]byte(couchDBTemplate), []common.YamlReplace{
			{
				Old: "${containerCouchdbName}",
				New: couchDB.ContainerName,
			},
			{
				Old: "${db}",
				New: couchDB.DB,
			},
			{
				Old: "${image}",
				New: couchDB.Image,
			},
			{
				Old: "${5984}",
				New: couchDB.Ports.Couchdb5984,
			},
		})

		fileName = couchDB.ContainerName + ".docker-compose-yaml"
		log.Println("生成文件名称", fileName)

		res = append(res, File{
			Name: fileName,
			Data: tmp,
		})
	}

	return res, nil
}
