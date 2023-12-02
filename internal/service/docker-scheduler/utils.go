package docker_scheduler

import (
	"strconv"
)

// string转int64相加后再转回string
func addStr(str1, str2 string) string {

	val1, _ := strconv.ParseInt(str1, 10, 64)
	val2, _ := strconv.ParseInt(str2, 10, 64)

	tmp := val2 + val1

	return strconv.FormatInt(int64(tmp), 10)
}

// string转int64进行比较
func strCompare(str1, str2 string) int {
	val1, _ := strconv.ParseInt(str1, 10, 64)
	val2, _ := strconv.ParseInt(str2, 10, 64)

	tmp := val1 - val2

	if tmp == 0 {
		return 0
	} else if tmp < 0 {
		return -1
	}
	return 1
}

// string转int64进行比较，返回最大值
func strMax(str1, str2 string) string {
	if strCompare(str1, str2) > 0 {
		return str1
	}

	return str2
}

// string转int64进行比较，返回最小值
func strMin(str1, str2 string) string {
	if strCompare(str1, str2) < 0 {
		return str1
	}

	return str2
}

// 判断数字(字符串类型)是否连续,字符串需先统一排序：从小到大或从大到小
func strContinue(str1, str2 string) bool {

	val1, _ := strconv.ParseInt(str1, 10, 64)
	val2, _ := strconv.ParseInt(str2, 10, 64)

	tmp := val2 - val1
	if tmp == 1 || tmp == -1 {
		return true
	}

	return false
}

// 根据区间获得返回的新区间的起点和终点
func getValidPortHelper(group []string, start, end string) (string, string) {

	if len(group) == 0 || len(group) > 2 {
		return "", ""
	}

	if len(group) == 1 {
		return addStr(group[0], start), addStr(group[0], end)
	} else {
		return addStr(group[0], start), addStr(group[1], end)
	}

}

// 根据主机打分排序
func sortVHostByScore(host VHost, vHost []VHost) []VHost {

	if len(vHost) == 0 {
		return []VHost{host}
	}

	var res []VHost

	var i int
	for i = range vHost {

		if host.Score > vHost[i].Score { //找到插入点
			break
		}
	}

	res = append(res, vHost[:i]...)
	res = append(res, host)
	res = append(res, vHost[i:]...)

	return res
}

// 获取两个区间的交集
func getIntersection(s1, e1, s2, e2 string) string {

	var res string

	if strCompare(s1, e2) > 0 { //区间1的起点大于区间2的终点，无交集
		return res
	}

	if strCompare(s2, e1) > 0 { //区间2的起点大于区间1的终点，无交集
		return res
	}

	var start, end string

	//确定起点
	start = strMax(s1, s2) //区间1的起点和区间2的起点最大值为起点
	//确定终点
	end = strMin(e1, e2) //区间1的终点和区间2的终点的最小值为终点

	if start == end {
		return start
	} else {
		return start + "-" + end
	}
}
