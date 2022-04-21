package main

import (
	"flag"
	"fmt"
	"os"
	"rexagent/pkg/conf"
	"rexagent/pkg/server"
)

type arrayFlags []string

func (f *arrayFlags) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var configs arrayFlags
	var configDirs arrayFlags
	params := make(map[string]string)
	flag.Var(&configs, "conf", "config files path")
	flag.Var(&configDirs, "confdir", "config directories path")

	flag.Parse()

	config, err := conf.LoadXmlConfig(configs, configDirs, params)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
	err = server.NewServer(&config).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

}
