package plugins

import (
	"context"
	"github.com/kcat-io/go-plugin/plugin"
	"log"
	"time"
)

func init() {
	plugin.Register("hello", HelloPlugin())
}

func HelloPlugin() plugin.Plugin {
	return func(ctx context.Context, in, out interface{}, next plugin.NextHandle) (err error) {
		log.Printf("enter hello plugin")
		begin := time.Now()
		// do sth. before the next handle
		if next != nil {
			err = next(ctx, in, out)
		}
		// do sth. after the last handle
		end := time.Now()
		log.Printf("exit from hello plugin, cost:%s", end.Sub(begin))
		return err
	}
}
