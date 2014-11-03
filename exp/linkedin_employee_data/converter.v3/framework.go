package converter

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/wangkuiyi/hmm/core"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	kNumFields = 12
)

type Record struct {
	Entry, Member                       string
	Begin, End                          time.Time
	Company, Title, Seniority, Function string
	School, Degree, DegreeRank, Field   string
}

var (
	flagCSV    = flag.String("csv", "", "Input from convert_v3.awk")
	flagCorpus = flag.String("corpus", "", "Output corpus JSON file")
)

func Run(gen Generator) {
	flag.Parse()

	if csv, e := os.Open(*flagCSV); e != nil {
		log.Fatalf("Cannot open CSV file: %v", e)
	} else {
		defer csv.Close()

		if exps, e := LoadCSV(csv); e != nil {
			log.Fatalf("LoadCSV: %v", e)
		} else {
			if corpus, e := os.Create(*flagCorpus); e != nil {
				log.Printf("Cannot create corpus file: %v", e)
			} else {
				defer corpus.Close()
				GenerateJSON(exps, gen, corpus)
			}
		}
	}
}

func LoadCSV(csv io.Reader) (map[string][]*Record, error) {
	ret := make(map[string][]*Record)
	s := bufio.NewScanner(csv)
	for s.Scan() {
		// It is important to call strings.Split instead of
		// strings.Fields, because convert_v3.awk use "\t" to seprate
		// fields strictly.
		if fs := strings.Split(s.Text(), "\t"); len(fs) != kNumFields {
			log.Printf("len(%v) != kNumFields (%d)", fs, kNumFields)
		} else {
			member := fs[1]
			if _, ok := ret[member]; !ok {
				ret[member] = make([]*Record, 0)
			}
			ret[member] = append(ret[member], &Record{
				Entry:  fs[0],
				Member: member,
				Begin:  α(ParseDate(fs[2])).(time.Time),
				End:    α(ParseDate(fs[3])).(time.Time),

				Company:   fs[4],
				Title:     fs[5],
				Seniority: fs[6],
				Function:  fs[7],

				School:     fs[8],
				Degree:     fs[9],
				DegreeRank: fs[10],
				Field:      fs[11]})
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
		log.Fatalf("cannot parse date %s", date)
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

func GenerateJSON(exps map[string][]*Record, gen Generator, corpus io.Writer) {
	correct := 0 // count corrected encoded instances.
	en := json.NewEncoder(corpus)

	// Sort memberSk so to access exps in deterministic order.
	members := make([]string, 0, len(exps))
	for member, _ := range exps {
		members = append(members, member)
	}
	sort.Strings(members)

	//  JSON-encode instances, each corresponds to a member.
	for _, memberSk := range members {
		exp := exps[memberSk]
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
					for _, f := range gen.Feature(r, c) {
						if len(f) > 0 {
							inst.Obs[y][c][f]++
						}
					}
				}
			}
		}

		if y := emptyYear(inst.Obs); y < 0 {
			if e := en.Encode(inst); e != nil {
				log.Fatalf("Cannot encode instance: %v, %e", inst, e)
			}
			correct++
		} else {
			log.Printf("Error: Member %s has gap year %d", memberSk, y+minYear)
		}
	}

	log.Printf("Output %d instances", correct)
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
