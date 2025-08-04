package main

import "github.com/nondramatic/practice/smooth"

func main() {

	c := smooth.NewSuperServerConfig()
	c.HttpListenAddr = ":28567"
	server := smooth.SuperServer{
		Executor:    Main,
		ListenAddrs: []string{":9090"},
		Config:      c,
	}
	_ = server.Run()
}

func Main(fds []*smooth.Fd) {

}
