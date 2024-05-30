package cron

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/tinyredglasses/workers3/internal/jsutil"
	"github.com/tinyredglasses/workers3/internal/runtimecontext"
)

type Task func(ctx context.Context) error

var scheduledTask Task

func runScheduler(eventObj js.Value) error {
	ctx := runtimecontext.New(context.Background(), eventObj)
	if err := scheduledTask(ctx); err != nil {
		return err
	}
	return nil
}

func init() {
	runSchedulerCallback := js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 1 {
			panic(fmt.Errorf("invalid number of arguments given to runScheduler: %d", len(args)))
		}
		eventObj := args[0]
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				err := runScheduler(eventObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(js.Undefined())
			}()
			return js.Undefined()
		})

		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("runScheduler", runSchedulerCallback)
}

//go:wasmimport workers ready
func ready()

// ScheduleTask sets the Task to be executed
func ScheduleTask(task Task) {
	scheduledTask = task
	ready()
	select {}
}
