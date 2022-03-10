package plugins

import (
	"context"
	"github.com/kcat-io/go-plugin/plugin"
	"log"
)

func init() {
	plugin.Register("hello2", HelloPlugin2())
}

func HelloPlugin2() plugin.Plugin {
	return func(ctx context.Context, in, out interface{}, next plugin.NextHandle) (err error) {
		log.Printf("enter hello2 plugin")

		pIn := in.(*string)
		// modify input value in plugin
		*pIn = "kCat"
		// do sth. before the next handle
		if next != nil {
			err = next(ctx, in, out)
		}
		// do sth. after the last handle
		log.Printf("exit from hello2 plugin")
		return err
	}
}
