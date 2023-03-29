package cloudflare

import (
	"context"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

func WaitUntil(ctx context.Context, task func()) {
	executionContext := cfruntimecontext.GetExecutionContext(ctx)

	executionContext.Call("waitUntil", jsutil.NewPromise(js.FuncOf(func(this js.Value, args []js.Value) any {
		task()
		return nil
	})))
}