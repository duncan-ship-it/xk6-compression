package compression

import (
	"fmt"

	"github.com/grafana/sobek"
	"github.com/klauspost/compress/zstd"
	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/compression", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a VU
		vu          modules.VU
		compression *Compression
	}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

type Compression struct {
	vu modules.VU
}

// NewModuleInstance implements the modules.Module interface returning a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:          vu,
		compression: &Compression{vu: vu},
	}
}

func (m *Compression) ZstdCompress(data sobek.Value) sobek.Value {
	zw, err := zstd.NewWriter(nil)
	if err != nil {
		panic(fmt.Errorf("failed to initialize zstd Writer: %v", err))
	}
	defer zw.Close()

	rt := m.vu.Runtime()

	exported := data.Export()
	d, ok := exported.([]byte)

	if !ok {
		panic("expecting JS array")
	}

	dst := make([]byte, 0, len(d))
	dst = zw.EncodeAll(d, dst)

	return rt.ToValue(dst)
}

func (m *Compression) ZstdDecompress(compressed sobek.Value) sobek.Value {
	zw, err := zstd.NewReader(nil)

	if err != nil {
		panic(fmt.Errorf("failed to initialize zstd Reader: %v", err))
	}

	defer zw.Close()

	exported := compressed.Export()
	data, ok := exported.([]byte)

	if !ok {
		panic("expecting JS array")
	}

	out, decodeErr := zw.DecodeAll(data, nil)
	if decodeErr != nil {
		panic(fmt.Errorf("failed to decode data: %v", decodeErr))
	}

	rt := m.vu.Runtime()
	return rt.ToValue(out)
}

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: mi.compression,
	}
}
