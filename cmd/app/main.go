package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/MaxFando/application-design/configs"
	"github.com/MaxFando/application-design/internal/adapters/driving/http"
	"github.com/MaxFando/application-design/internal/providers"
	"github.com/MaxFando/application-design/pkg/utils"
)

func main() {
	utils.InitializeLogger("INFO")

	ctx, cancel := context.WithCancel(context.Background())
	configs.InitializeConfig()

	repositoryProvider := providers.NewRepositoryProvider()
	repositoryProvider.RegisterDependencies()

	serviceProvider := providers.NewServiceProvider()
	serviceProvider.RegisterDependencies(repositoryProvider)

	useCaseProvider := providers.NewUseCaseProvider()
	useCaseProvider.RegisterDependencies(serviceProvider)

	ctx = context.WithValue(ctx, providers.UseCaseProviderKey, useCaseProvider)

	httpServer := http.NewHttpServer(http.NewHandler(ctx), ":"+configs.Config.HTTPServer.Port)
	httpServer.Serve()
	utils.Logger.Info("Приложение стартовало в режиме", zap.String("log_level", configs.Config.App.LogLevel), zap.String("env", configs.Config.App.Env))
	utils.Logger.Info("На порту " + configs.Config.HTTPServer.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		utils.Logger.Info("app - Run - signal: " + s.String())
	case errHTTP := <-httpServer.Notify():
		utils.Logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", errHTTP))
	}

	cancel()

	_ = httpServer.Shutdown()
	utils.Logger.Info("HTTP сервер завершил работу")
	utils.Logger.Info("Приложение завершило работу")
}
