package compression

import (
	"fmt"

	"github.com/grafana/sobek"
	"github.com/klauspost/compress/zstd"
	"go.k6.io/k6/js/modules"
)

type CompressionModule struct{}

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/compression", new(CompressionModule))
}

func (m *CompressionModule) zstdCompress(data []byte, rt *sobek.Runtime) sobek.Value {
	zw, err := zstd.NewWriter(nil)

	if err != nil {
		panic(fmt.Errorf("failed to initialize zstd Writer: %v", err))
	}

	defer zw.Close()

	dst := make([]byte, 0, len(data))
	dst = zw.EncodeAll(data, dst)

	return rt.ToValue(dst)
}

func (m *CompressionModule) zstdDecompress(compressed []byte, rt *sobek.Runtime) sobek.Value {
	zw, err := zstd.NewReader(nil)

	if err != nil {
		panic(fmt.Errorf("failed to initialize zstd Reader: %v", err))
	}

	defer zw.Close()

	out, decodeErr := zw.DecodeAll(compressed, nil)
	if decodeErr != nil {
		panic(fmt.Errorf("failed to decode data: %v", decodeErr))
	}

	return rt.ToValue(out)
}
