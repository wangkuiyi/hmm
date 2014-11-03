package converter

import (
	"encoding/json"
	"github.com/wangkuiyi/hmm/core"
	"io"
	"log"
)

func LoadJSON(r io.Reader) []*core.Instance {
	ret := make([]*core.Instance, 0)
	dec := json.NewDecoder(r)
	for {
		inst := new(core.Instance)
		if e := dec.Decode(inst); e != nil {
			if e == io.EOF {
				break
			} else {
				log.Fatalf("Failed decode instance: %v", e)
			}
		}
		ret = append(ret, inst)
	}
	return ret
}
