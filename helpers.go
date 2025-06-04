package compression

import (
	"fmt"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
)

// convert Sobek array to Go byte array
func ToNativeBytes(rt *sobek.Runtime, bytes sobek.Value) []byte {
	exported := bytes.Export()
	d, ok := exported.([]byte)

	if !ok {
		common.Throw(rt, fmt.Errorf("type error: expecting JS array"))
	}

	return d
}
