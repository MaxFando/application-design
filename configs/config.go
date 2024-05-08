package configs

import (
	"sync"
)

type config struct {
	App        App
	HTTPServer HTTPServer
}

var (
	Config config
	once   sync.Once
)

func InitializeConfig() {
	once.Do(func() {
		Config = config{
			App: App{
				ServiceName: "ms-github.com/MaxFando/application-design",
				Env:         "local",
				LogLevel:    "info",
			},
			HTTPServer: HTTPServer{
				Port: "80",
			},
		}
	})
}
