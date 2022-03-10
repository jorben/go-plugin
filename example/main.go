package main

import (
	"context"
	_ "github.com/kcat-io/go-plugin/example/plugins"
	"github.com/kcat-io/go-plugin/plugin"
	"gopkg.in/yaml.v3"
	"log"
	"time"
)

var configContent = `
# 插件列表
plugins:
  - hello
  - hello2

`

type Config struct {
	PluginList []string `yaml:"plugins"`
}

// HelloMain 主执行内容
func HelloMain(ctx context.Context, in interface{}, out interface{}) error {
	pIn := in.(*string)
	pOut := out.(*string)
	*pOut = "Hello, " + *(pIn) + "!"
	time.Sleep(1 * time.Second)
	return nil
}

func main() {

	// 加载配置内容
	c := &Config{}
	err := yaml.Unmarshal([]byte(configContent), c)
	if err != nil {
		log.Fatal(err)
	}
	// 加载配置中的规则
	var ps = plugin.Plugins{}
	for _, p := range c.PluginList {
		f := plugin.GetPlugin(p)
		if f == nil {
			log.Fatalf("plugin %s no registered", p)
		}
		ps = append(ps, f)
	}

	// 创建会话
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 执行调用 主函数
	strIn := "joy"
	strOut := ""
	err = ps.Handle(ctx, &strIn, &strOut, HelloMain)
	if err != nil {
		log.Fatal(err)
	}

	// Hello, kCat!
	log.Println(strOut)

}
