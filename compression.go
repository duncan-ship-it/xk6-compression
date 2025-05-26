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

func (m *Compression) zstdCompress(data []byte) sobek.Value {
	zw, err := zstd.NewWriter(nil)

	if err != nil {
		panic(fmt.Errorf("failed to initialize zstd Writer: %v", err))
	}

	defer zw.Close()

	dst := make([]byte, 0, len(data))
	dst = zw.EncodeAll(data, dst)

	rt := m.vu.Runtime()
	return rt.ToValue(dst)
}

func (m *Compression) zstdDecompress(compressed []byte) sobek.Value {
	zw, err := zstd.NewReader(nil)

	if err != nil {
		panic(fmt.Errorf("failed to initialize zstd Reader: %v", err))
	}

	defer zw.Close()

	out, decodeErr := zw.DecodeAll(compressed, nil)
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
