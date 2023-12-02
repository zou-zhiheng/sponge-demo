package docker_scheduler

import (
	"fmt"
	"log"
	"testing"
	"user/configs"
	"user/internal/config"
	"user/internal/model"
)

func init() {
	err := config.Init(configs.Path("user.yaml"))
	if err != nil {
		panic(err)
	}
}

var vHost = []VHost{
	{
		IP:        "192.168.182.134",
		PortRange: "1-660;662-65535",
		Hardware: VHostHardware{
			UsedCpus:      0.2,
			UsedMemory:    2,
			UsedDataDisk:  5,
			TotalCpus:     12,
			TotalMemory:   10,
			TotalDataDisk: 100,
		},
	},
	{
		IP:        "192.168.182.133",
		PortRange: "1-660;662-65535",
		Hardware: VHostHardware{
			UsedCpus:      0.3,
			UsedMemory:    2,
			UsedDataDisk:  5,
			TotalCpus:     12,
			TotalMemory:   10,
			TotalDataDisk: 100,
		},
	},
}

var pod = VHost{
	Hardware: VHostHardware{
		UsedCpus:     0.2,
		UsedMemory:   1,
		UsedDataDisk: 2,
	},
	Filter: VHostFilter{
		PortOccupy:    "21",
		ContainerName: "test",
	},
}

func TestScheduler(t *testing.T) {
	s := NewScheduler()
	//获取硬件打分
	vh, err := s.GetHardwareScore(pod, vHost, VHostWeight)
	if err != nil {
		fmt.Println("GetHardwareScore===>\n", err)
		return
	}
	for i := range vh {
		fmt.Println(vh[i].IP, vh[i].Score)
	}
	fmt.Println("GetHardwareScore===>OK")
	fmt.Printf("\n\n\n")

	//根据主机规则过滤
	res, err := s.GetVHostByFilter(pod, vh)

	if err != nil {
		fmt.Println("GetVHostByFilter===>\n", err)
	} else {
		fmt.Println("GetVHostByFilter===>OK")
	}

	for i := range res {
		fmt.Println(res[i].IP, res[i].Filter.PortOccupy, res[i].Filter.PortFree)
	}

}

func TestScheduler_GetOptimalHost(t *testing.T) {
	s := NewScheduler()
	//获取最优主机
	host, vhost, err := s.GetOptimalHost(pod, vHost, VHostWeight)
	if err != nil {
		fmt.Println("GetOptimalHost===>\n", err)
		for i := range vhost { //打印主机可用端口范围，其他占用信息在ssh登录时已知
			fmt.Println(vhost[i].IP, vHost[i].Filter.PortFree)
		}
		return
	}

	fmt.Println("GetOptimalHost===>OK", host)

}

func TestSchedulerDemo(t *testing.T) {
	s := NewScheduler()
	//获取主机信息
	db := model.GetDB() //dsn:"root:123456@(10.0.51.35:3306)/hll?parseTime=true&loc=Local&charset=utf8,utf8mb4"
	var vHostDB []model.TbVirtualHostInfo
	rep := db.Raw("SELECT * from tb_virtual_host_info where fwal_ip!=\"192.168.189.128\" GROUP BY fwal_ip").Scan(&vHostDB)
	if rep.Error != nil {
		panic(rep.Error)
	}
	log.Println("获取数据库中的虚拟主机信息")
	vhost := s.GetVHostInfo(vHostDB)
	if vhost == nil {
		fmt.Println("GetVHostInfo===>无主机信息")
		return
	}
	fmt.Println("GetVHostInfo===>OK")
	fmt.Printf("\n\n\n")

	//获取硬件打分
	vh, err := s.GetHardwareScore(pod, vhost, VHostWeight)
	if err != nil {
		fmt.Println("GetHardwareScore===>\n", err)
		return
	}
	for i := range vh {
		fmt.Println(vh[i].IP, vh[i].Score)
	}
	fmt.Println("GetHardwareScore===>OK")
	fmt.Printf("\n\n\n")

	//根据主机规则过滤
	res, err := s.GetVHostByFilter(pod, vh)

	if err != nil {
		fmt.Println("GetVHostByFilter===>\n", err)
	} else {
		fmt.Println("GetVHostByFilter===>OK")
	}

	for i := range res {
		fmt.Println(res[i].IP, res[i].Filter.PortOccupy, res[i].Filter.PortFree)
	}

}
