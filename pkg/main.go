package main

import (
	"easyexec/pkg/login"
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
	case "tail":
		flag.StringVar(&f, "f", "/home/tomcat/logs/*/app.log", "文件路径")
		parse()
		fmt.Println("tail")
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
Usage: easyexec [ login | tail | exec | logs] [Options]...

Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}
