package docker_scheduler

type VHostResource int

// 需要统一单位,可以全部是k,全部是MB，全部是G等，统一单位即可；硬件资源不能用string类型
const (
	// 单项资源默认最大数量
	ContainerResourceMax float64 = 100000

	CPU VHostResource = iota
	GPU
	Memory
	BandWidth //带宽
	SystemDisk
	DataDisk
)

var (
	// 资源权重，可自定义权重
	VHostWeight = map[VHostResource]float64{
		CPU:       2,
		GPU:       1,
		Memory:    2,
		BandWidth: 1,
		DataDisk:  1,
	}
)
