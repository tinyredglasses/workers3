package cloudflare

import (
	"context"
	"syscall/js"

	"github.com/tinyredglasses/workers3/cloudflare/internal/cfruntimecontext"
	"github.com/tinyredglasses/workers3/internal/runtimecontext"
)

// Getenv gets a value of an environment variable.
//   - https://developers.cloudflare.com/workers/platform/environment-variables/
//   - This function panics when a runtime context is not found.
func Getenv(name string) string {
	return cfruntimecontext.MustGetRuntimeContextEnv().Get(name).String()
}

// GetBinding gets a value of an environment binding.
//   - https://developers.cloudflare.com/workers/platform/bindings/about-service-bindings/
//   - This function panics when a runtime context is not found.
func GetBinding(name string) js.Value {
	return cfruntimecontext.MustGetRuntimeContextEnv().Get(name)
}

func GetRuntimeContextValue(name string) js.Value {
	return cfruntimecontext.MustGetRuntimeContextValue(name)
}

func GetCtx(name string) js.Value {
	return cfruntimecontext.MustGetExecutionContext().Get(name)
}

func GetTriggerObject(ctx context.Context) js.Value {
	return runtimecontext.MustExtractTriggerObj(ctx)
}
