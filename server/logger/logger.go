package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	Log      *zap.Logger
	onceInit sync.Once
)

func Init(lvl int) error {
	var err error

	onceInit.Do(func() {
		globalLvl := zapcore.Level(lvl)

		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})

		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= globalLvl && lvl < zapcore.ErrorLevel
		})

		consoleInfo := zapcore.Lock(os.Stdout)
		consoleErrs := zapcore.Lock(os.Stderr)

		zcnf := zap.NewProductionEncoderConfig()
		consoleEncoder := zapcore.NewJSONEncoder(zcnf)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrs, highPriority),
			zapcore.NewCore(consoleEncoder, consoleInfo, lowPriority),
		)

		Log = zap.New(core)
		zap.RedirectStdLog(Log)

	})

	return err
}
