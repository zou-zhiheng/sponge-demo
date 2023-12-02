package docker_scheduler

import (
	"errors"
	"log"
	"strings"
)

// 根据规则过滤不适合的主机
func (s *Scheduler) GetVHostByFilter(pod VHost, vHost []VHost) ([]VHost, error) {

	defer Logger("GetVHostByFilter")()

	var res, lost []VHost
	var found bool //当前主机是否满足Pod所需资源,主机是否拥有pod资源项
	for i := range vHost {
		found = true
		log.Println("根据规则校验主机" + vHost[i].IP)

		//端口校验
		vHost[i].Filter.PortOccupy = s.StrFormat(vHost[i].Filter.PortOccupy, "\n")
		vHost[i].Filter.PortFree = s.GetValidPortByPortRange(vHost[i].Filter.PortOccupy, vHost[i].PortRange)
		ok, err := s.IsValidPort(pod.Filter.PortOccupy, vHost[i].Filter.PortFree)
		if ok != "" || err != nil {
			found = false
		}

		//容器名称校验
		fOK := s.IsValidContainerName(pod.Filter.ContainerName, vHost[i].Filter.ContainerName)
		if !fOK {
			found = false
		}

		log.Println("主机"+vHost[i].IP+"是否通过校验", found)
		host := VHost{}
		host.IP = vHost[i].IP
		if found {
			//找到合适的主机，返回主机IP
			res = append(res, host)
		} else {
			//端口可用列表、已使用的容器名称
			host.Filter.PortFree = vHost[i].Filter.PortFree
			host.Filter.ContainerName = vHost[i].Filter.ContainerName
			lost = append(lost, host)
		}
	}

	if res == nil {
		return lost, errors.New("无符合条件主机===>返回各主机可用列表")
	}

	return res, nil

}

// 端口是否可用
func (s *Scheduler) IsValidPort(portSrc string, vHostPorts string) (string, error) {

	ports := strings.Split(portSrc, ";")

	if len(ports) == 0 {
		return "", errors.New("端口不能为空")
	}

	group := strings.Split(vHostPorts, ";")
	var vStart, vEnd, pStart, pEnd string
	var i, j int

	var found bool
	//判断是否在可用区间内
	for i = 0; i < len(ports); i++ {
		found = false
		pStart, pEnd = getValidPortHelper(strings.Split(ports[i], "-"), "0", "0")
		for j = 0; j < len(group); j++ {
			vStart, vEnd = getValidPortHelper(strings.Split(group[j], "-"), "0", "0")
			//判断是否在可用区间内
			if strCompare(pStart, vStart) >= 0 && strCompare(vEnd, pEnd) >= 0 {
				found = true
				break
			}
		}

		//有端口不在可用区间内
		if !found {
			return vHostPorts, errors.New("端口已占用")
		}
	}

	return "", nil
}

// 容器名称是否可用 ,可用:ture ; 不可用:false
func (s *Scheduler) IsValidContainerName(pod string, vhost string) bool {

	names := strings.Split(vhost, "\n")
	for i := range names {
		if pod == names[i] {
			return false
		}
	}

	return true
}
