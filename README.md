# go-plugin
[![Go](https://github.com/kcat-io/go-plugin/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/kcat-io/go-plugin/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/kcat-io/go-plugin?status.svg)](http://godoc.org/github.com/kcat-io/go-plugin)

一个插件库（插件引擎），可以使你的Go程序轻松支持插件式扩展。

- 这是一个链表式插件库，它把各个插件按顺序连接起来
- 每个插件均可在next handle的前后分别实现`前序逻辑`和`后续逻辑`
- 各插件的执行顺序如下图所示
[![go-plugin](https://kcat.io/usr/uploads/2022/03/3388664059.png)](https://kcat.io)

## Install
```shell
go get -u github.com/kcat-io/go-plugin/plugin
```

## Example
- Please refer to the case in the example directory

编写一个插件非常的简单
```go
package plugins

import (
	"context"
	"github.com/kcat-io/go-plugin/plugin"
	"log"
)

func init() {
	// 在这里向插件库注册声明本插件
	plugin.Register("demoPlugin", DemoPlugin())
}

// DemoPlugin 插件实现逻辑
func DemoPlugin() plugin.Plugin {
	return func(ctx context.Context, in, out interface{}, next plugin.NextHandle) (err error) {
		log.Printf("enter demoPlugin")
		// do sth. before the next handle
		// 可以在这里实现`前序逻辑`
		
		if next != nil {
			// 这里是链式调用下一个插件或func的入口
			err = next(ctx, in, out)
		}
		
		// do sth. after the last handle
		// 可以在这里实现`后续逻辑`
		log.Printf("exit from demoPlugin")
		return err
	}
}
```

插件编写好之后我们只需要在main函数中调用即可
```go
package main

import (
	"context"
	// 这里需要import刚编写的插件，此处以example中的示例为参考
	_ "github.com/kcat-io/go-plugin/example/plugins"
	"github.com/kcat-io/go-plugin/plugin"
	"log"
	"time"
)

// DemoMain 在这里实现你的常规业务逻辑
func DemoMain(ctx context.Context, in interface{}, out interface{}) error {
	pIn := in.(*string)
	pOut := out.(*string)
	*pOut = "Hello, " + *(pIn) + "!"
	return nil
}

func main(){
	// 在这里注明需要启用的插件（通常把它们放在配置文件中）
	var PluginList = []string{
		"demoPlugin",
    }
	
	// 加载配置中的插件
	var ps = plugin.Plugins{}
	for _, p := range PluginList {
		f := plugin.GetPlugin(p)
		if f == nil {
			log.Fatalf("Plugin %s is not registered yet.", p)
		}
		ps = append(ps, f)
	}

	// 创建会话
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 执行调用 主函数
	strIn := "kCat"
	strOut := ""
	err = ps.Handle(ctx, &strIn, &strOut, DemoMain)
	if err != nil {
		log.Fatal(err)
	}

	// Hello, kCat!
	log.Println(strOut)
}
```

## License
MIT - See [LICENSE][license] file

[license]: https://github.com/kcat-io/go-plugin/blob/master/LICENSE