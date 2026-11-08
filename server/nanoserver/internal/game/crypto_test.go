package game

import (
	"testing"

	"github.com/lonng/nano/pipeline"
)

func BenchmarkCrypto_Inbound(b *testing.B) {
	c := &Crypto{[]byte("hKKJdfskj997sdSk")}
	payload := []byte(`[{"name":"test","length":1.06666672229767,"segments":[{"t":0.233333334326744,"v":4.44000005722046},{"t":0.200000002980232,"v":2.62499976158142},{"t":0.266666650772095,"v":0.686249911785126},{"t":0.166666686534882,"v":1.34915959835052},{"t":0.200000047683716,"v":2.28395414352417}]}]`)
	msg := &pipeline.Message{Data: payload}
	_ = c.outbound(nil, msg)
	for i := 0; i < b.N; i++ {
		c.inbound(nil, msg)
	}
}

func BenchmarkCrypto_Outbound(b *testing.B) {
	c := &Crypto{[]byte("hKKJdfskj997sdSk")}
	payload := []byte(`[{"name":"test","length":1.06666672229767,"segments":[{"t":0.233333334326744,"v":4.44000005722046},{"t":0.200000002980232,"v":2.62499976158142},{"t":0.266666650772095,"v":0.686249911785126},{"t":0.166666686534882,"v":1.34915959835052},{"t":0.200000047683716,"v":2.28395414352417}]}]`)
	for i := 0; i < b.N; i++ {
		msg := &pipeline.Message{Data: payload}
		c.outbound(nil, msg)
	}
}
