package core

import "github.com/hezof/log"

/*
	框架内置组件. 例如日志之类, 往往贯穿所有组件的初始流程.
*/

func InitLogger() error {
	cfg := new(log.FileConfig)
	ok, err := ConfigStruct("log", cfg, "")
	if err != nil {
		return err
	}
	if ok {
		lgr, err := log.NewFileLogger(cfg)
		if err != nil {
			return err
		}
		log.InitLogger(lgr)
	}
	return nil
}
