package cloudflare

import (
	"fmt"
	"syscall/js"

	"github.com/tinyredglasses/workers3/internal/jsutil"

	"github.com/tinyredglasses/workers3/cloudflare/internal/cfruntimecontext"
)

// DurableObjectStorage represents interface of Durable Object's Storage instance.
type DurableObjectStorage struct {
	instance js.Value
}

type DurableObjectStorageGetOptions struct {
	AllowConcurrency bool
	NoCache          bool
}

type DurableObjectListOptions struct {
	Start            string
	StartAfter       string
	End              string
	Prefix           string
	Reverse          bool
	Limit            int
	AllowConcurrency bool
	NoCache          bool
}

type DurableObjectListResponse struct {
}

type DurableObjectPutDeleteOptions struct {
	AllowUnconfirmed bool
	AllowConcurrency bool
	NoCache          bool
}

func (opts *DurableObjectStorageGetOptions) toJS(type_ string) js.Value {
	obj := jsutil.NewObject()
	// obj.Set("type", type_)
	if opts == nil {
		return obj
	}
	if opts.AllowConcurrency {
		obj.Set("allowConcurrency", opts.AllowConcurrency)
	}
	if opts.NoCache {
		obj.Set("noCache", opts.NoCache)
	}
	return obj
}

func (opts *DurableObjectListOptions) toJS(type_ string) js.Value {
	obj := jsutil.NewObject()
	// obj.Set("type", type_)
	if opts == nil {
		return obj
	}
	if opts.Start != "" {
		obj.Set("start", opts.Start)
	}
	if opts.StartAfter != "" {
		obj.Set("startAfter", opts.StartAfter)
	}
	if opts.End != "" {
		obj.Set("end", opts.End)
	}
	if opts.Prefix != "" {
		obj.Set("prefix", opts.Prefix)
	}
	if opts.Reverse {
		obj.Set("reverse", opts.Reverse)
	}
	if opts.Limit != 0 {
		obj.Set("limit", opts.Limit)
	}
	if opts.AllowConcurrency {
		obj.Set("allowConcurrency", opts.AllowConcurrency)
	}
	if opts.NoCache {
		obj.Set("noCache", opts.NoCache)
	}
	return obj
}

func (opts *DurableObjectPutDeleteOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}

	obj := jsutil.NewObject()

	if opts.AllowUnconfirmed {
		obj.Set("allowUnconfirmed", opts.AllowUnconfirmed)
	}

	if opts.AllowConcurrency {
		obj.Set("allowConcurrency", opts.AllowConcurrency)
	}
	if opts.NoCache {
		obj.Set("noCache", opts.NoCache)
	}
	return obj
}

func (d *DurableObjectStorage) Get(key string, opts DurableObjectStorageGetOptions) (string, error) {
	p := d.instance.Call("get", key, opts.toJS("text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

func (d *DurableObjectStorage) GetMany(keys []string, opts DurableObjectStorageGetOptions) (map[string]string, error) {
	p := d.instance.Call("get", keys, opts.toJS("text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return jsutil.MapToMap(v), nil
}

func (d *DurableObjectStorage) List(opts DurableObjectListOptions) (map[string]string, error) {
	p := d.instance.Call("list", opts.toJS("text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return jsutil.MapToMap(v), nil
}

func (d *DurableObjectStorage) DeleteAll(opts DurableObjectPutDeleteOptions) (string, error) {
	p := d.instance.Call("deleteAll", opts.toJS())
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

// PutString puts string value into KV with key.
//   - if a network error happens, returns error.
func (d *DurableObjectStorage) PutString(key string, value string, opts DurableObjectPutDeleteOptions) error {
	p := d.instance.Call("put", key, value, opts.toJS())
	_, err := jsutil.AwaitPromise(p)
	if err != nil {
		return err
	}
	return nil
}

func (d *DurableObjectStorage) PutStrings(entries map[string]string, opts DurableObjectPutDeleteOptions) error {
	val := jsutil.MapToStrRecord(entries)
	toString := val.Call("toString").String()
	fmt.Println("PutString val: " + toString)
	p := d.instance.Call("put", val, opts.toJS())
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
