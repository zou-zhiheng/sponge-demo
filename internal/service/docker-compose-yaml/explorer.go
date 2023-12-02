package docker_compose_yaml

import (
	"errors"
	"log"
	"user/internal/service/common"
)

type ExplorerYaml struct {
}

type ExplorerCompose struct {
	ContainerName   string //容器名称
	ContainerDBName string //数据库容器名称
	ImageExplorer   string //浏览器镜像
	ImageExplorerDB string //数据库镜像
	PGData          string //postgresql存储数据的路径
	Certs           string //证书根目录
	FabricConfig    string //fabric配置
	WalletStore     string //钱包存储路径
	Ports           Ports  //要使用的端口
}

func (e *ExplorerYaml) Create(explorers []ExplorerCompose) ([]File, error) {
	defer Logger("NodeCompose.ExplorerYaml.Create")()

	if len(explorers) == 0 {
		return nil, errors.New("参数不能空")
	}

	var res []File
	var fileName string
	for _, explore := range explorers {

		log.Println("开始创建浏览器:", explore.ContainerName)
		//浏览器节点名称校验
		log.Println("节点名称校验")
		if explore.ContainerName == "" {
			log.Println(explore.ContainerName, "浏览器节点名称不能为空")
			break
		}

		//浏览器数据库名称校验
		log.Println(explore.ContainerName, "浏览器数据库节点名称")
		if explore.ContainerDBName == "" {
			log.Println(explore.ContainerName, "浏览器数据库节点名称不能为空")
			break
		}

		//证书校验
		log.Println(explore.ContainerName, "证书路径校验")
		if explore.Certs == "" {
			log.Println(explore.ContainerName, "证书根目录不能为空")
			break
		}

		//fabric配置校验
		log.Println(explore.ContainerName, "fabric配置路径校验")
		if explore.FabricConfig == "" {
			log.Println(explore.ContainerName, "fabric配置路径不能为空")
			break
		}

		//钱包存储校验
		log.Println(explore.ContainerName, "钱包存储路径校验")
		if explore.WalletStore == "" {
			log.Println(explore.ContainerName, "钱包存储路径不能为空")
			break
		}

		//浏览器镜像
		log.Println(explore.ContainerName, "镜像校验")
		if explore.ImageExplorer == "" {
			log.Println(explore.ContainerName, "镜像不能为空")
			break
		}

		//浏览器数据库镜像
		log.Println(explore.ContainerName, "浏览器数据库镜像校验")
		if explore.ImageExplorerDB == "" {
			log.Println(explore.ContainerName, "浏览器数据库镜像不能为空")
			break
		}

		//浏览器数据库存储
		log.Println(explore.ContainerName, "浏览器数据库存储路径校验")
		if explore.PGData == "" {
			log.Println(explore.ContainerName, "浏览器数据库存储路径不能为空")
			break
		}

		//端口赋值
		log.Println(explore.ContainerName, "镜像赋值")
		if explore.Ports.Explorer8080 == "" { //不存在
			log.Println(explore.ContainerName, "端口传入为空，使用默认端口")
			explore.Ports.Explorer8080 = "8080"
		}

		tmp := common.ReplaceVar([]byte(explorerTemplate), []common.YamlReplace{
			{
				Old: "${containerExplorerName}",
				New: explore.ContainerName,
			},
			{
				Old: "${containerExplorerDBName}",
				New: explore.ContainerDBName,
			},
			{
				Old: "${exploreDBImage}",
				New: explore.ImageExplorerDB,
			},
			{
				Old: "${exploreImage}",
				New: explore.ImageExplorer,
			},
			{
				Old: "${pgdata}",
				New: explore.PGData,
			},
			{
				Old: "${certs}",
				New: explore.Certs,
			},
			{
				Old: "${fabricConfig}",
				New: explore.FabricConfig,
			},
			{
				Old: "${walletstore}",
				New: explore.WalletStore,
			},
			{
				Old: "${8080}",
				New: explore.Ports.Explorer8080,
			},
		})

		fileName = explore.ContainerName + ".docker-compose-yaml"
		log.Println("生成文件名称", fileName)

		res = append(res, File{
			Name: fileName,
			Data: tmp,
		})
	}

	return res, nil
}
