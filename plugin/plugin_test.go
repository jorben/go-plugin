package plugin_test

import (
	"context"
	"fmt"
	"github.com/kcat-io/go-plugin/plugin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testPlugin(id int) plugin.Plugin {
	return func(ctx context.Context, in, out interface{}, next plugin.NextHandle) (err error) {
		//fmt.Printf("enter %d plugin\n", id)
		*(in.(*string)) = fmt.Sprintf("enter %d plugin;", id) + *in.(*string)
		// do sth. before the next handle
		if next != nil {
			err = next(ctx, in, out)
		}
		*(out.(*string)) = fmt.Sprintf("exit from %d plugin;", id) + *out.(*string)
		// do sth. after the last handle
		//fmt.Printf("exit from %d plugin\n", id)
		return err
	}
}

func helloHandle(ctx context.Context, in, out interface{}) error {
	pIn := in.(*string)
	pOut := out.(*string)
	*pOut = "Hello, " + *pIn
	return nil
}

func TestGetPlugin(t *testing.T) {
	f := plugin.GetPlugin("NotExist")
	assert.Nil(t, f)
}

func TestRegister(t *testing.T) {
	name := "test"
	plugin.Register(name, testPlugin(0))
	f := plugin.GetPlugin(name)
	assert.NotNil(t, f)
}

func TestHelloHandel(t *testing.T) {
	in := "kcat"
	out := ""
	ps := plugin.Plugins{}
	err := ps.Handle(context.Background(), &in, &out, helloHandle)
	assert.Nil(t, err)
	assert.Equal(t, in, "kcat")
	assert.Equal(t, out, "Hello, kcat")
}

func TestNilHandel(t *testing.T) {
	in := "kcat"
	out := ""
	ps := plugin.Plugins{}
	err := ps.Handle(context.Background(), &in, &out, nil)
	assert.Nil(t, err)
}

func TestHandel(t *testing.T) {
	in := "kcat"
	out := ""
	ps := plugin.Plugins{
		testPlugin(0),
		testPlugin(1),
		testPlugin(2),
		testPlugin(3),
	}
	err := ps.Handle(context.Background(), &in, &out, nil)
	assert.Nil(t, err)
	assert.Equal(t, in, "enter 3 plugin;enter 2 plugin;enter 1 plugin;enter 0 plugin;kcat")
	assert.Equal(t, out, "exit from 0 plugin;exit from 1 plugin;exit from 2 plugin;exit from 3 plugin;")

}
