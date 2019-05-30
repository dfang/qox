package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	Logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(Logger)
}

// var singleton zap.Logger
// var once sync.Once

// // Init initializes a thread-safe singleton logger
// // This would be called from a main method when the application starts up
// // This function would ideally, take zap configuration, but is left out
// // in favor of simplicity using the example logger.
// func Init() {
// 	// once ensures the singleton is initialized only once
// 	once.Do(func() {
// 		singleton = zap.New(
// 			zap.NewJSONEncoder(),
// 		)
// 	})
// }

// // Log a message at the given level with given fields
// func Log(level zap.Level, message string, fields ...zap.Field) {
// 	singleton.Log(level, message, fields...)
// }

// // Debug logs a debug message with the given fields
// func Debug(message string, fields ...zap.Field) {
// 	singleton.Debug(message, fields...)
// }

// // Info logs a debug message with the given fields
// func Info(message string, fields ...zap.Field) {
// 	singleton.Info(message, fields...)
// }

// // Warn logs a debug message with the given fields
// func Warn(message string, fields ...zap.Field) {
// 	singleton.Warn(message, fields...)
// }

// // Error logs a debug message with the given fields
// func Error(message string, fields ...zap.Field) {
// 	singleton.Error(message, fields...)
// }

// // Fatal logs a message than calls os.Exit(1)
// func Fatal(message string, fields ...zap.Field) {
// 	singleton.Fatal(message, fields...)
// }
