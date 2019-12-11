//go:generate go run pkg/common/gensaconfig/main.go -kubeconfig pkg/common/gensaconfig/config
package main

import (
	"easyexec/pkg/common"
	exec2 "easyexec/pkg/exec"
	"easyexec/pkg/list"
	"easyexec/pkg/login"
	logs2 "easyexec/pkg/logs"
	"flag"
	"fmt"
	"os"
)

var (
	h bool
	p string
	d string
	n string
	f string
	u string
	i int
	l int
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助")
	flag.Usage = usage
}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
	}
	opt := os.Args[1]
	if opt == "login" {
		flag.StringVar(&u, "u", "", "Rancher用户名")
		flag.StringVar(&p, "p", "", "Rancher密码")
		parse()
		l := login.NewLogin(u, p)
		err := l.Login()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("登录成功")
		}
		os.Exit(0)
	}

	commonFlag()
	switch opt {
	case "tail", "tailf":
		flag.StringVar(&f, "f", "/home/tomcat/logs/*/app.log", "文件路径")
		flag.IntVar(&l, "l", 20, "初始时显示行数")
		commonType := parse()
		plusf := ""
		if opt == "tailf" {
			plusf = "-f"
		}
		cmd := fmt.Sprintf("tail %s -n %d %s", plusf, l, f)
		exec := exec2.NewExec(commonType, cmd)
		exec.Exec()
	case "ll", "ls":
		flag.StringVar(&f, "f", "/home/tomcat/logs/*/", "文件(文件夹)路径")
		commonType := parse()
		cmd := fmt.Sprintf("ls -l %s", f)
		exec := exec2.NewExec(commonType, cmd)
		exec.Exec()
	case "cat":
		flag.StringVar(&f, "f", "/home/tomcat/logs/*/app.log", "文件路径")
		commonType := parse()
		cmd := fmt.Sprintf("cat %s", f)
		exec := exec2.NewExec(commonType, cmd)
		exec.Exec()
	case "podlist":
		commonType := parse()
		podlist := list.NewPodList(commonType)
		show := podlist.List2String()
		fmt.Println(show)
	case "exec":
		parse()
		fmt.Println("Coming Soon")
		os.Exit(0)
	case "logs":
		commonType := parse()
		logs := logs2.NewLogs(commonType)
		logs.Logs()
	default:
		flag.Usage()
	}
}

func commonFlag() {
	flag.StringVar(&p, "p", "", "pod名，deployment和pod二选其一")
	flag.StringVar(&d, "d", "", "deployment名，deployment和pod二选其一")
	flag.StringVar(&n, "n", "", "namespace，如果不填将会智能查找")
	flag.IntVar(&i, "i", 1, "容器序号")
}

func parse() common.CommonType {
	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		flag.Usage()
		os.Exit(0)
	}
	if h {
		flag.Usage()
		os.Exit(0)
	}
	commonType := &common.CommonType{
		PodName:        p,
		Namespace:      n,
		DeploymentName: d,
		ContainerIndex: i,
	}
	return *commonType
}

func usage() {
	fmt.Fprintf(os.Stderr, `easyexec version: easyexec/1.0.1
Usage: easyexec [ login | podlist | tail |tailf | ll | ls | cat | logs] [Options]...

login: 登录操作，所有操作之前需要先登录。
podlist: 根据条件显示当前用户可用的POD列表。
tail: 打印指定运行中的POD内指定的文件末尾数行。
tailf: 持续打印指定运行中的POD内指定的文件，等同于tail -f。
ll: 查看指定目录下文件信息。
ls: 查看指定目录下文件信息。
cat: 打印文件全部内容。(Linux下运行可使用命令: 'easyexec cat -d xxx > abc.txt'将内容保存至本地)
logs: 持续打印POD的STDOUT信息，等同于docker logs。

Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}
