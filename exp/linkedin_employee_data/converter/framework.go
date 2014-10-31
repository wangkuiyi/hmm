package converter

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/wangkuiyi/hmm/core"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	kNumFields = 11
)

type Record struct {
	CompanyOrSchool       string
	PosOrEduRank          string
	TitleOrDegree         string
	SeniorityOrDegreeRank string
	FunctionOrField       string
	Begin                 time.Time
	End                   time.Time
	IsJob                 bool
	IsEdu                 bool
}

var (
	flagCSV    = flag.String("csv", "", "Input CSV file")
	flagCorpus = flag.String("corpus", "", "Output corpus JSON file")
)

func Run(gen Generator) {
	flag.Parse()

	if csv, e := os.Open(*flagCSV); e != nil {
		log.Fatalf("Cannot open CSV file %s: %v", *flagCSV, e)
	} else {
		defer csv.Close()

		if exp, e := LoadCSV(csv); e != nil {
			log.Printf("LoadCSV: %v", e)
		} else {
			if corpus, e := os.Create(*flagCorpus); e != nil {
				log.Printf("Cannot create corpus file: %v", e)
			} else {
				defer corpus.Close()
				exps := GenerateJSON(exp, gen)
				e := json.NewEncoder(corpus).Encode(exps)
				if e != nil {
					log.Printf("Failed generating corpus")
				} else {
					log.Printf("Wrote %d members", len(exps))
				}
			}
		}
	}
}

func LoadCSV(csv io.Reader) (map[string][]*Record, error) {
	ret := make(map[string][]*Record)
	firstLine := true
	s := bufio.NewScanner(csv)
	for s.Scan() {
		if firstLine {
			firstLine = false // Skip the table header.
			continue
		}

		if fs := strings.Fields(s.Text()); len(fs) != kNumFields {
			log.Printf("len(fs) != kNumFs. fs=%v", fs)
		} else {
			if _, ok := ret[fs[1]]; !ok {
				ret[fs[1]] = make([]*Record, 0)
			}
			ret[fs[1]] = append(ret[fs[1]], &Record{
				CompanyOrSchool:       fs[2],
				PosOrEduRank:          fs[3],
				TitleOrDegree:         fs[4],
				SeniorityOrDegreeRank: fs[5],
				FunctionOrField:       fs[6],
				Begin:                 α(ParseDate(fs[7])).(time.Time),
				End:                   α(ParseDate(fs[8])).(time.Time),
				IsJob:                 α(strconv.ParseBool(fs[9])).(bool),
				IsEdu:                 α(strconv.ParseBool(fs[10])).(bool)})
		}
	}

	if e := s.Err(); e != nil {
		log.Fatalf("reading CSV:", e)
	}

	return ret, nil
}

func ParseDate(date string) (time.Time, error) {
	fs := strings.Split(date, "/")
	if len(fs) != 3 {
		return time.Now(), nil
	}
	return time.Date(
		α(strconv.Atoi(fs[2])).(int),
		time.Month(α(strconv.Atoi(fs[0])).(int)),
		α(strconv.Atoi(fs[1])).(int),
		0, 0, 0, 0, time.Local), nil
}

// α is used to get the first of multiple return values, when the
// second parameter, must be of type error, is nil.
func α(args ...interface{}) interface{} {
	if args[1] != nil {
		log.Fatalf("α error: %v", args[1].(error))
	}
	return args[0]
}

func GenerateJSON(exps map[string][]*Record, gen Generator) []*core.Instance {
	ret := make([]*core.Instance, 0, len(exps))
	for memberSk, exp := range exps {
		minYear := 3000
		maxYear := 1000
		for _, r := range exp {
			if r.Begin.Year() < minYear {
				minYear = r.Begin.Year()
			}
			if r.End.Year() > maxYear {
				maxYear = r.End.Year()
			}
		}

		years := maxYear - minYear + 1
		if years <= 0 {
			log.Printf("Member %s has 0 years of experience")
			continue
		}

		inst := &core.Instance{
			Obs:     makeObservedMatrix(years, gen.NumChannels()),
			Periods: ones(years)}

		for _, r := range exp {
			for year := r.Begin.Year(); year <= r.End.Year(); year++ {
				y := year - minYear
				for c := 0; c < gen.NumChannels(); c++ {
					if f := gen.Feature(r, c); len(f) > 0 {
						inst.Obs[y][c][f]++
					}
				}
			}
		}

		if y := emptyYear(inst.Obs); y < 0 {
			ret = append(ret, inst)
		} else {
			log.Printf("Member %s has no feature in %d", memberSk, y+minYear)
		}
	}
	return ret
}

func makeObservedMatrix(years, channels int) [][]core.Observed {
	ret := make([][]core.Observed, years)
	for i, _ := range ret {
		ret[i] = make([]core.Observed, channels)
		for c, _ := range ret[i] {
			ret[i][c] = make(core.Observed)
		}
	}
	return ret
}

func ones(length int) []int {
	ret := make([]int, length)
	for i := 0; i < length; i++ {
		ret[i] = 1
	}
	return ret
}

func emptyYear(obs [][]core.Observed) int {
	for y, ob := range obs {
		empty := true
		for ch, _ := range ob {
			if len(ob[ch]) > 0 {
				empty = false
				break
			}
		}

		if empty {
			return y
		}
	}
	return -1
}
