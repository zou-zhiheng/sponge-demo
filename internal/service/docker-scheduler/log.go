package docker_scheduler

import "log"

func Logger(name string) func(detail ...string) {
	log.Println("执行" + name + "===>")
	return func(detail ...string) {
		if detail != nil {
			log.Println("执行中的详细信息", detail)
		}
		log.Println("执行" + name + "===>结束")
	}
}
