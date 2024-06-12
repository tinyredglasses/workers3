package adhoc

import (
	"context"
	"fmt"
	"github.com/tinyredglasses/workers3/internal/jsutil"
	"github.com/tinyredglasses/workers3/internal/runtimecontext"
	"syscall/js"
)

var (
	handler Handler
)

type Handler interface {
	Handle(ctx context.Context, reqObj js.Value)
}

type HandlerCreator func(ctx context.Context) Handler

func init() {

	handleDataCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {

		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to handle: %d", len(args)))
		}
		eventObj := args[0]

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := handle(eventObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("handle", handleDataCallback)
}

func handle(event js.Value) error {
	ctx := runtimecontext.New(context.Background(), event)

	handler.Handle(ctx, event)
	return nil
}

//go:wasmimport workers ready
func ready()

func Handle(hc HandlerCreator) {
	ctx := runtimecontext.New(context.Background(), js.Value{})
	handler = hc(ctx)
	ready()
	select {}
}
