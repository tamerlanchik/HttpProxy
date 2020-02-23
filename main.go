package main

import (
	"fmt"
	"github.com/mkideal/cli"
	"os"
	"proxy/proxy"
	"time"
)

var (
	tunnelTimeout time.Duration
)
var root = &cli.Command{
	Desc: "https://github.com/bozaro/tech-db-forum",
	Argv: func() interface{} { return nil },
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		//os.Exit(EXIT_INVALID_COMMAND)
		fmt.Errorf("Wrong command")
		return nil
	},
}
//	Test:
//	curl -Lv --proxy https://localhost:12345 --proxy-cacert server.pem https://google.com

//	Commands
//	proxy - run proxy server
//		-port
//		-pem
//		-key
//		-timeout
//	reverse	- run reverse tool
//		list [n=20] - print latest n requests
//		repeat [id] - repeat [id] request
//		-d [id] - delete [id] request

type ProxyArgs struct {
	cli.Helper
	Port string `cli:"p,port" usage:"Working port for proxy" dft:"8080"`
	Pem string `cli:"pem" usage:"Path to PEM file" dft:"keys/server.pem"`
	Key string `cli:"key" usage:"Path to KEY file" dft:"keys/server.key"`
	Timeout int `cli:"t,timeout" usage:"Timeout for tls-tunnel, ms" fdt:"1000"`
}
var cmdProxy = &cli.Command{
	Name:               "proxy",
	Desc:               "Run proxy server",
	Argv:               func() interface{} {return new(ProxyArgs)},
	Fn:                 func(ctx *cli.Context) error {
		argv := ctx.Argv().(*ProxyArgs)

		err := proxy.Run(argv.Port, argv.Pem, argv.Key, argv.Timeout)
		fmt.Println(err)
		return nil
	},
}

var cmdReverse = &cli.Command{
	Name:               "reverse",
	Desc:               "",
	Fn:                 func(ctx *cli.Context) error {
		fmt.Println("Reverse")
		return nil
	},
	Argv:               nil,
}

type ReverseListArgs struct {
	cli.Helper
	N int `cli:"n" usage:"Number of last requests to display, 0 - all" dft:"20"`
}
var cmdReverseList = &cli.Command{
	Name:               "list",
	Desc:               "Display the list of the latest requests",
	Fn:                 func(ctx *cli.Context) error {
		//_ := ctx.Argv().(*ProxyArgs)

		fmt.Println("List")
		return nil
	},
	Argv:               func() interface{} {return new(ReverseListArgs)},
}

type ReverseRepeatArgs struct {
	cli.Helper
	Id string `cli:"i,id" usage:"Repeat the request with specified id"`
}
var cmdReverseRepeat = &cli.Command{
	Name:               "repeat",
	Desc:               "Repea the specified request",
	Fn:                 func(ctx *cli.Context) error {
		//_ := ctx.Argv().(*ProxyArgs)

		fmt.Println("List")
		return nil
	},
	Argv:               func() interface{} {return new(ReverseRepeatArgs)},
}

type ReverseDeleteArgs struct {
	cli.Helper
	Id string `cli:"i,id" usage:"Delete the request with specified id"`
}
var cmdReverseDelete = &cli.Command{
	Name:               "delete",
	Desc:               "Delete the specified request",
	Fn:                 func(ctx *cli.Context) error {
		//_ := ctx.Argv().(*ProxyArgs)

		fmt.Println("Delete")
		return nil
	},
	Argv:               func() interface{} {return new(ReverseDeleteArgs)},
}

func main() {
	cliRoot := cli.Root(root, cli.Tree(cmdProxy), cli.Tree(cmdReverse, cli.Tree(cmdReverseList)))

	if nil != cliRoot.Run(os.Args[1:]) {
		fmt.Errorf("Error: %s", os.Stderr)
		return
	}



}


