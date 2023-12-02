package docker_scheduler

import (
	"fmt"
	"testing"
)

func TestStrFormat(t *testing.T) {
	s := NewScheduler()
	//解析从服务器拿到的端口信息
	strFormat := s.StrFormat("\n22\n53\n179\n6010\n9099\n9253\n9254\n9353\n10248\n10249\n10250\n10256\n39449\n", "\n")
	fmt.Println("StrFormat主机端口已被占用列表===>", strFormat)

	//获取可用端口区间
	getValidPort := s.GetValidPortByPortRange(strFormat, "1-6553;6556-65535")
	fmt.Println("GetValidPort主机端口未被使用列表===>", getValidPort)

	//判断端口是否可用
	isValidPort, err := s.IsValidPort("10249", getValidPort)
	if err != nil {
		fmt.Println("IsValidPort", err)
		return
	}

	fmt.Println("IsValidPort主机端口可被使用的列表===>", isValidPort, "端口可用")

}
