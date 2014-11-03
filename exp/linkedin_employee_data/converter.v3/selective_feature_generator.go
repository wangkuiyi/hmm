package converter

import (
	"log"
)

type selectiveFeatureGenerator struct {
	p *PlainFeatureGenerator
}

func (*selectiveFeatureGenerator) NumChannels() int {
	return 1
}

func (s *selectiveFeatureGenerator) Feature(r *Record, year, ch int) []string {
	if r.Member == "105509708" ||
		r.Member == "103275950" ||
		r.Member == "63143482" ||
		r.Member == "99944053" ||
		r.Member == "63631912" {
		log.Println(r.Member, year, s.p.Feature(r, year, ch))
		return s.p.Feature(r, year, ch)
	}
	return nil
}
