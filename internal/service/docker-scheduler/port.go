package docker_scheduler

import (
	"strings"
)

// 处理后的格式 "8001-8005;9001-9005" ,src(待处理的字符串，字符串需有序),split(待处理的字符串中的分割符)
func (s *Scheduler) StrFormat(src string, split string) string {
	if len(src) == 0 {
		return ""
	} else if len(src) == 1 || len(split) > len(src) {
		return src
	}

	var res string
	str := strings.Split(src, split)
	//判断头部是否有分割符号
	if str[0] == "" {
		str = str[1:]
	}
	//判断尾部是否有分割符号
	if str[len(str)-1] == "" {
		str = str[:len(str)-1]
	}

	var start, end string
	var flag = len(str) - 1
	start = str[0]
	end = start

	for i := 1; i < len(str); i++ {
		if str[i] == "" {
			continue
		}
		//比较当前值与end是否连续
		if strContinue(end, str[i]) {
			//更新终点
			end = str[i]
			if i == flag { //是否是最后一个节点
				//创建区间
				res += start + "-" + end
				break
			}
		} else { //不连续则创建新区间
			//起点和终点相同则合并
			if start == end {
				res += start + ";"
			} else {
				res += start + "-" + end + ";"
			}

			//重新初始化区间
			start = str[i]
			end = start

			if i == flag { //是否是最后一个节点
				res += start
			}
		}

	}

	return res
}

// 返回可用端口的区间段;vHostPorts 格式 "8001-8005;9001-9005";返回格式 "1-8000;8006-9000;9006-65535"
func (s *Scheduler) GetValidPort(vHostPorts string) string {

	if len(vHostPorts) == 0 {
		return ""
	} else if len(vHostPorts) <= 5 {
		return vHostPorts
	}

	var res, start, end, index string
	var group []string
	//按区间分组
	vHostPortGroup := strings.Split(vHostPorts, ";")

	group = strings.Split(vHostPortGroup[0], "-") //单独处理第一组
	end, index = getValidPortHelper(group, "-1", "1")
	if strCompare(end, "1") > 0 {
		if strCompare(end, "3") > 0 {
			res += "1-" + end + ";"
		} else { // 只有一个端口(1)
			res += end + ";"
		}

		//start = index
		//end = start
	}

	start = index
	end = start

	for i := 1; i < len(vHostPortGroup); i++ {
		group = strings.Split(vHostPortGroup[i], "-")
		end, index = getValidPortHelper(group, "-1", "1")
		if start == end {
			res += start + ";"
		} else {
			res += start + "-" + end + ";"
		}

		start = index
		end = start

	}

	//处理最后一个区间
	if strCompare("65535", start) > 0 {
		if strCompare("65534", start) > 0 {
			res += start + "-65535"
		} else {
			res += start
		}
	}

	return res
}

// 指定端口的使用范围来获取当前端口的可用区间， vHostPorts(已经使用的端口) 格式 "8001-8005;9001-9005"; portRange(主机开放的端口) 端口开放的范围格式 1-65535; 返回格式 "1-8000;8006-9000;9006-65535"
func (s *Scheduler) GetValidPortByPortRange(vHostPorts string, portRanges string) string {

	var res string
	//主机开放端口为空，则使用默认值:1-65535
	if portRanges == "" {
		portRanges = "1-65535"
	}

	if len(vHostPorts) == 0 {
		return portRanges
	}

	//获取主机全部的端口可用区间
	allRange := strings.Split(s.GetValidPort(vHostPorts), ";")
	//获取主机已使用的端口区间
	portRange := strings.Split(portRanges, ";")

	//allRange,vHostPort已有序，必须

	var s1, e1, s2, e2, tmp string
	for i := 0; i < len(allRange); i++ {
		for j := 0; j < len(portRange); j++ {
			//获取区间1
			s1, e1 = getValidPortHelper(strings.Split(allRange[i], "-"), "0", "0")
			//获取区间2
			s2, e2 = getValidPortHelper(strings.Split(portRange[j], "-"), "0", "0")
			//获取区间1和区间2的交集
			tmp = getIntersection(s1, e1, s2, e2)
			if tmp != "" {
				res += tmp + ";"
				break
			}
		}
	}

	if len(res) != 0 {
		res = res[:len(res)-1]
	}

	return res
}
