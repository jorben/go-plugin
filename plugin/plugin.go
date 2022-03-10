package plugin

import (
	"context"
	"sync"
)

var (
	plugins = make(map[string]Plugin)
	lock    = sync.RWMutex{}
)

type (
	// NextHandle 链接接口定义
	NextHandle func(ctx context.Context, in, out interface{}) error
	// Plugin 插件回包定义
	Plugin func(ctx context.Context, in, out interface{}, next NextHandle) error
	// Plugins 插件组
	Plugins []Plugin
)

// Register 注册插件
func Register(name string, f Plugin) {
	lock.Lock()
	plugins[name] = f
	lock.Unlock()
}

// GetPlugin 获取插件方法
func GetPlugin(name string) Plugin {
	lock.RLock()
	f := plugins[name]
	lock.RUnlock()
	return f
}

// Handle 链式插件执行入口方法
func (s Plugins) Handle(ctx context.Context, in, out interface{}, next NextHandle) error {
	// 把已注册的插件装载成链表，头节点为第一个插件
	for i := len(s) - 1; i >= 0; i-- {
		n, p := next, s[i]
		next = func(ctx context.Context, in, out interface{}) error {
			return p(ctx, in, out, n)
		}
	}
	if next != nil {
		return next(ctx, in, out)
	}
	return nil
}
