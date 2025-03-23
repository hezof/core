package core

import (
	"fmt"
	"github.com/hezof/log"
)

/*************************************************
 * 容器功能
 *************************************************/

// Register 注册组件工厂, 重复注册会panic!
func Register(base string, factory ManagedFactory) {
	err := _managedContext.RegisterFactory(base, factory)
	if err != nil {
		panic(fmt.Errorf("register factory %v error: %v", base, err))
	}
}

// Component 断言组件实例, 若无base工厂会panic!
func Component[T any](base string, name string) T {
	component, err := _managedContext.RetrieveComponent(base, name)
	if err != nil {
		panic(fmt.Errorf("retrieve component %v.%v error: %v", base, name, err))
	}
	return component.(T)
}

/*************************************************
 * 核心功能
 *************************************************/

func Init() {
	datas, err := ReadFile(ConfigTomlFile())
	if err != nil {
		panic(fmt.Errorf("init config context error: %v", err))
	}
	InitTomlData(datas...)
}

func InitTomlData(datas ...[]byte) {
	// 0. 初始配置
	if err := _configContext.SetTomlData(datas...); err != nil {
		panic(fmt.Errorf("init config context error: %v", err))
	}

	// 初始化日志. 需早于其他组件!
	if err := InitLogger(); err != nil {
		panic(fmt.Errorf("init logger error: %v", err))
	}

	// 1. 回调钩子
	ExecHook(BeforeInit, _configContext, _managedContext)

	// 2. 初始托管(错误中止)
	if err := _managedContext.Init(_configContext); err != nil {
		panic(fmt.Errorf("init managed context error: %v", err))
	}

	// 最后回调钩子
	ExecHook(AfterInit, _configContext, _managedContext)
}

func Reload(reloadPolicy func(base string, config *ManagedConfig, newValues map[string]any) bool) error {
	datas, err := ReadFile(ConfigTomlFile())
	if err != nil {
		return fmt.Errorf("reload config context error: %v", err)
	}
	return ReloadTomlData(reloadPolicy, datas...)
}

func ReloadTomlData(reloadPolicy func(base string, config *ManagedConfig, newValues map[string]any) bool, datas ...[]byte) error {
	// 0. 重载配置
	if err := _configContext.SetTomlData(datas...); err != nil {
		return fmt.Errorf("reload config context error: %v", err)
	}

	// 1. 回调钩子
	ExecHook(BeforeReload, _configContext, _managedContext)

	// 2. 重载托管(错误忽略)
	_managedContext.Reload(_configContext, reloadPolicy)

	// 最后回调钩子
	ExecHook(AfterReload, _configContext, _managedContext)

	return nil
}

func Exit(hints ...func(base string, config *ManagedConfig, err error)) {

	/*
		注意: 因为Exit通常在defer中执行, 因此不建议再在Exit添加defer流程!
	*/

	// 退出开始执行flush logger(避免丢失关键消息)
	log.Flush()

	// 1. 回调钩子
	ExecHook(BeforeReload, _configContext, _managedContext)

	// 2. 退出托管(容忍错误)
	_managedContext.Exit(hints...)

	// 最后回调钩子
	ExecHook(AfterReload, _configContext, _managedContext)

	// 退出最后执行flush logger.
	log.Flush()
}
