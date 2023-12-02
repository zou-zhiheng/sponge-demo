package docker_scheduler

import (
	"errors"
	"fmt"
	"log"
	"user/internal/model"
	"user/internal/service/ssh"
)

type Scheduler struct {
}

type VHost struct {
	IP        string
	PortRange string //开放的端口范围，需有序（从小到大） 格式 "8001-8005;9001-9005"
	Hardware  VHostHardware
	Filter    VHostFilter
	Score     float64
}

type VHostHardware struct {
	UsedMemory   uint64
	UsedCpus     float32
	UsedDataDisk uint64

	TotalMemory   uint64
	TotalCpus     float32
	TotalDataDisk uint64
}

type VHostFilter struct {
	PortOccupy    string //端口号 格式 "8001-8005;9001-9005"
	PortFree      string //可使用的端口号 格式 "8001-8005;9001-9005"
	ContainerName string //容器名称
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// 获取硬件分数
func (s *Scheduler) GetHardwareScore(pod VHost, vHost []VHost, vHostWeight map[VHostResource]float64) ([]VHost, error) {
	defer Logger("GetHardwareScore")()

	//Pod和vHost默认不为空
	var res []VHost

	//计算每一台主机的得分
	var found bool  //当前主机是否满足Pod所需资源
	var tmp float64 //当前主机得分
	for i, vh := range vHost {

		tmp = 0
		found = true
		log.Println("计算主机硬件" + vh.IP + "得分")
		//CPU
		if (vh.Hardware.TotalCpus - vh.Hardware.UsedCpus) < pod.Hardware.UsedCpus {
			found = false
			break
		}
		tmp += (float64(vh.Hardware.TotalCpus-vh.Hardware.UsedCpus-pod.Hardware.UsedCpus) / float64(vh.Hardware.TotalCpus)) * vHostWeight[CPU]

		//内存
		if (vh.Hardware.TotalMemory - vh.Hardware.UsedMemory) < pod.Hardware.UsedMemory {
			found = false
			break
		}
		tmp += (float64(vh.Hardware.TotalMemory-vh.Hardware.UsedMemory-pod.Hardware.UsedMemory) / float64(vh.Hardware.TotalMemory)) * vHostWeight[CPU]

		//磁盘
		if float64(vh.Hardware.TotalDataDisk-vh.Hardware.UsedDataDisk-pod.Hardware.UsedDataDisk)/float64(vh.Hardware.TotalDataDisk) < 0.1 {
			found = false
			break
		}
		tmp += (float64(vh.Hardware.TotalDataDisk-vh.Hardware.UsedDataDisk-pod.Hardware.UsedDataDisk) / float64(vh.Hardware.TotalDataDisk)) * vHostWeight[DataDisk]

		log.Println("主机"+vh.IP+"硬件资源是否满足", found)
		if found { //找到合适的主机
			vHost[i].Score = tmp
			res = sortVHostByScore(vHost[i], res)
		}
	}

	if len(res) == 0 {
		return nil, errors.New("主机资源不满足")
	}

	return res, nil
}

// 获取主机信息，作为公共包时应该删除这个方法
func (s *Scheduler) GetVHostInfo(vHostDB []model.TbVirtualHostInfo) []VHost {
	defer Logger("GetHostInfo")()
	var vhost []VHost
	for _, vh := range vHostDB {
		host := VHost{}
		//登录到目标主机获取（端口、数据卷、容器名称信息）
		log.Println("登录主机" + vh.FwalIP)
		cli := ssh.NewCli(vh.FwalIP, vh.FsshUser, vh.FsshPwd)
		log.Println("获取主机" + vh.FwalIP + "端口信息")
		//获取端口信息
		output, err := cli.Run("{ ss -tuln | awk '/^tcp/ {print $5}' ;" +
			" ss -uln | awk '/^udp/ {print $5}' ; } | awk -F':' '{print $NF}' | sort -n | uniq")
		if err != nil {
			fmt.Println("ssh.Connect===>", vh.FwalIP, err)
			continue
		}
		host.Filter.PortOccupy = output

		log.Println("获取主机" + vh.FwalIP + "容器名称信息")
		//获取容器名称信息
		output, err = cli.Run("sudo docker ps -a --format \"{{.Names}}\"")
		if err != nil {
			continue
		}
		host.Filter.ContainerName = output

		////获取已挂载的数据卷
		//output, err = cli.Run("sudo docker ps -a -q | xargs -n 1 sudo docker inspect -f '{{range .Mounts}}{{.Source}}{{\"\\n\"}}{{end}}'")
		//if err != nil {
		//	continue
		//}
		//field[ReelDataOccupy] = FieldValue{Value: output}

		//内存
		host.Hardware.UsedMemory = vh.Fmemory
		host.Hardware.TotalMemory = 10000
		//CPU
		host.Hardware.UsedCpus = vh.Fcpu
		host.Hardware.TotalCpus = 12
		//磁盘
		host.Hardware.UsedDataDisk = vh.FdataDisk
		host.Hardware.TotalDataDisk = 10000
		//IP
		host.IP = vh.FwalIP

		//存储查询到的主机信息
		vhost = append(vhost, host)
	}

	return vhost

}
