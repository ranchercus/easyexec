package main

import (
	"easyexec/pkg/list"
	"easyexec/pkg/login"
	tail2 "easyexec/pkg/tailf"
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
	case "tailf":
		flag.StringVar(&f, "f", "/home/tomcat/logs/*/app.log", "文件路径")
		flag.IntVar(&l, "l", 20, "初始时显示行数")
		parse()
		tail := tail2.NewTail(f, l)
		tail.Namespace = n
		tail.DeploymentName = d
		tail.PodName = p
		tail.ContainerIndex = i
		tail.Tail()
	case "podlist":
		parse()
		podlist := list.NewPodList()
		podlist.Namespace = n
		podlist.DeploymentName = d
		podlist.PodName = p
		show := podlist.List2String()
		fmt.Println(show)
	case "exec":
		parse()
		fmt.Println("Coming Soon")
		os.Exit(0)
	case "logs":
		parse()
		fmt.Println("Coming Soon")
		os.Exit(0)
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

func parse() {
	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		flag.Usage()
		os.Exit(0)
	}
	if h {
		flag.Usage()
		os.Exit(0)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `easyexec version: easyexec/1.0.0
Usage: easyexec [ login | podlist | tailf | exec | logs] [Options]...

login: 登录操作，所有操作之前需要先登录。
podlist: 根据条件显示当前用户可用的POD列表。
tailf: 持续打印指定运行中的POD内指定的文件。
exec: 使用命令行进入POD。
logs: 持续打印POD的STDOUT信息。

Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}
