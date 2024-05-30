package adhoc

import (
	"context"
	"fmt"
	"github.com/tinyredglasses/workers3/internal/jsutil"
	"github.com/tinyredglasses/workers3/internal/runtimecontext"
	"syscall/js"
)

var (
	messageHandler MessageHandler
)

type MessageHandler interface {
	Handle(ctx context.Context, reqObj js.Value)
}

type MessageHandlerCreator func(ctx context.Context) MessageHandler

func init() {

	handleDataCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {

		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to handleData: %d", len(args)))
		}
		eventObj := args[0]

		//fsdf1 := js.Global().Get("JSON").Call("stringify", eventObj)
		//slog.Info(fsdf1.String())

		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := handleData(eventObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("handleData", handleDataCallback)
}

func handleData(event js.Value) error {
	ctx := runtimecontext.New(context.Background(), event)

	messageHandler.Handle(ctx, event)
	return nil
}

//go:wasmimport workers ready
func ready()

func Handle(mhc MessageHandlerCreator) {
	ctx := runtimecontext.New(context.Background(), js.Value{})
	messageHandler = mhc(ctx)
	ready()
	select {}
}
