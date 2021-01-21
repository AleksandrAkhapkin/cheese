package main

import (
	"flag"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/mysql"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/qiwi"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/yoomoney"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/server"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/server/handlers"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/service"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types/config"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	defer func() {
		r := recover()

		if r == nil {
			return
		}

		err := fmt.Errorf("PANIC:'%v'\nRecovered in: %s", r, infrastruct.IdentifyPanic())
		logger.LogError(err)
	}()
	configPath := new(string)
	debug := new(bool)
	flag.StringVar(configPath, "configs-path", "configs/configs.yaml", "path to yaml configs file")
	flag.BoolVar(debug, "debug", false, "debug-alarm")
	flag.Parse()
	f, err := os.Open(*configPath)
	if err != nil {
		logger.LogFatal(fmt.Errorf("err with open configs file %v, %s", err, *configPath))
	}
	logger.Debug = *debug
	cnf := config.Config{}
	if err = yaml.NewDecoder(f).Decode(&cnf); err != nil {
		logger.LogFatal(fmt.Errorf("err with parse configs %v, %s", err, *configPath))
	}

	err = logger.NewLogger(cnf.Telegram)
	if err != nil {
		logger.LogFatal(err)
	}

	pg, err := mysql.NewMySQL(cnf.MySQLDsn)
	if err != nil {
		logger.LogFatal(err)
	}

	srv, err := service.NewService(pg, &cnf)
	if err != nil {
		logger.LogFatal(err)
	}

	err = yoomoney.NewYoomoney(cnf.YoomoneySecret)
	if err != nil {
		logger.LogFatal(err)
	}

	err = qiwi.NewQiwi(cnf.YoomoneySecret)
	if err != nil {
		logger.LogFatal(err)
	}

	handls := handlers.NewHandlers(srv, &cnf)
	logger.CheckDebug()
	logger.LogInfo(fmt.Sprintf("Start server in port: %s", cnf.ServerPort))
	server.StartServer(handls, cnf.ServerPort, cnf.UseSSL, cnf.CertPathSSL, cnf.KeyPathSSL)
}
