package converter

import (
	"log"
)

type selectedMembers struct {
	p *PlainFeatureGenerator
}

func (*selectedMembers) NumChannels() int {
	return 1
}

func (s *selectedMembers) Feature(r []*Record, year, ch int) []string {
	if r[0].Member == "105509708" ||
		r[0].Member == "103275950" ||
		r[0].Member == "63143482" ||
		r[0].Member == "99944053" ||
		r[0].Member == "63631912" {
		log.Println(r[0].Member, year, s.p.Feature(r, year, ch))
		return s.p.Feature(r, year, ch)
	}
	return nil
}
