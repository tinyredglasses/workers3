package cloudflare

import (
	"fmt"
	"github.com/tinyredglasses/workers3/internal/jsutil"
	"io"
	"syscall/js"

	"github.com/tinyredglasses/workers3/cloudflare/internal/cfruntimecontext"
)

// DurableObjectStorage represents interface of Durable Object's Storage instance.
type DurableObjectStorage struct {
	instance js.Value
}

type DurableObjectStorageGetOptions struct {
	allowConcurrency bool
	noCache          bool
}

type DurableObjectListOptions struct {
	start            string
	startAfter       string
	end              string
	prefix           string
	reverse          bool
	limit            int
	allowConcurrency bool
	noCache          bool
}

type DurableObjectPutOptions struct {
	allowUnconfirmed bool
	allowConcurrency bool
	noCache          bool
}

func (opts *DurableObjectStorageGetOptions) toJS(type_ string) js.Value {
	obj := jsutil.NewObject()
	obj.Set("type", type_)
	if opts == nil {
		return obj
	}
	if opts.allowConcurrency {
		obj.Set("allowConcurrency", opts.allowConcurrency)
	}
	if opts.noCache {
		obj.Set("noCache", opts.noCache)
	}
	return obj
}

func (opts *DurableObjectListOptions) toJS(type_ string) js.Value {
	obj := jsutil.NewObject()
	obj.Set("type", type_)
	if opts == nil {
		return obj
	}
	if opts.start != "" {
		obj.Set("start", opts.start)
	}
	if opts.startAfter != "" {
		obj.Set("startAfter", opts.startAfter)
	}
	if opts.end != "" {
		obj.Set("end", opts.end)
	}
	if opts.prefix != "" {
		obj.Set("prefix", opts.prefix)
	}
	if opts.reverse {
		obj.Set("reverse", opts.reverse)
	}
	if opts.limit != 0 {
		obj.Set("limit", opts.limit)
	}
	if opts.allowConcurrency {
		obj.Set("allowConcurrency", opts.noCache)
	}
	if opts.noCache {
		obj.Set("noCache", opts.noCache)
	}
	return obj
}

func (opts *DurableObjectPutOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}

	obj := jsutil.NewObject()

	if opts.allowUnconfirmed {
		obj.Set("allowUnconfirmed", opts.allowUnconfirmed)
	}

	if opts.allowConcurrency {
		obj.Set("allowConcurrency", opts.allowConcurrency)
	}
	if opts.noCache {
		obj.Set("noCache", opts.noCache)
	}
	return obj
}

func (d *DurableObjectStorage) Get(key string, opts *DurableObjectStorageGetOptions) (string, error) {
	p := d.instance.Call("get", key, opts.toJS("text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

func (d *DurableObjectStorage) GetReader(key string, opts *DurableObjectStorageGetOptions) (io.Reader, error) {
	p := d.instance.Call("get", key, opts.toJS("stream"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return jsutil.ConvertReadableStreamToReadCloser(v), nil
}

// PutString puts string value into KV with key.
//   - if a network error happens, returns error.
func (d *DurableObjectStorage) PutString(key string, value string, opts *DurableObjectPutOptions) error {
	p := d.instance.Call("put", key, value, opts.toJS())
	_, err := jsutil.AwaitPromise(p)
	if err != nil {
		return err
	}
	return nil
}

func NewDurableObjectStorage(varName string) (*DurableObjectStorage, error) {
	inst := cfruntimecontext.MustGetRuntimeContextValue("storage")

	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &DurableObjectStorage{instance: inst}, nil
}
