package cmdPlus

import (
	"fmt"
	"os/exec"
	"time"
)


func main() {
	cmd := exec.Command("/Users/justinh/Downloads/wdtsend", "-transfer_id", "334455", "-destination", "127.0.0.1", "-start_port", "23456", "Office 安装程序.pkg")
	cc := spider.NewCmdPlus(cmd, "/Users/justinh/Downloads");

	cc.SetTriggerRegexpKeys(`^.+progress\s(?P<progress>.+),.+\s(?P<completed>.+%),\sAverage\sthroughput\s(?P<Average>.+),\sRecent\sthroughput\s(?P<Recent>.+)\.`)

	cc.TriggerKeyCallback = func(vals map[string]string) {
		fmt.Println("----TriggerRegexp:", vals)
	}

	cc.OutPutCallback = func(v string) {
		fmt.Println("VVVV:", v)
	}
	isFinish := false
	cc.FinishCallback = func(e error) {
		isFinish = true
	}

	cc.Exec()
	fmt.Println("OOOooooooooo")
	for !isFinish {
		time.Sleep(time.Millisecond * 100)

	}


	fmt.Println("所有操作已完成！")
	//	fmt.Println("code！",cc.ReturnCode)
}
