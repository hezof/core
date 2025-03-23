package core

import "github.com/hezof/log"

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

func ExitLogger() {
	log.Flush()
}
