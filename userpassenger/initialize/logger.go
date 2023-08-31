package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
