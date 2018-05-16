package main

import (
	"fmt"

	"github.com/deadcheat/awsset/examples/viper/configbin"
	"github.com/spf13/viper"
)

// Server example struct for config
type Server struct {
	Host string
	Port int
}

func main() {
	viper.SetConfigType("toml")
	f, _ := configbin.Assets.File("/config/configfile.toml")
	viper.ReadConfig(f)
	var s Server
	_ = viper.UnmarshalKey("server", &s)
	fmt.Printf("%#+v", s)
}
