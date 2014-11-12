package corpus_generation

import (
	"fmt"
	"github.com/wangkuiyi/hmm/exp/corpus_generation/loader"
	"github.com/wangkuiyi/buildnrun"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const (
	kCSVDir  = "github.com/wangkuiyi/hmm/exp/linkedin_employee_data"
	kCSVFile = "pos_and_edu_for_LI_employees.csv"
)

func TestGenerateAndLoad(t *testing.T) {
	dir, e := ioutil.TempDir("", "converter_test")
	if e != nil {
		t.Fatalf("Cannot create temp dir: %v", e)
	}
	defer os.RemoveAll(dir)

	*flagCSV = buildnrun.Pkg(path.Join(kCSVDir, kCSVFile))
	*flagCorpus = path.Join(dir, "corpus.json")

	Run(new(PlainFeatureGenerator))

	corpus := loader.LoadJSON(α(os.Open(*flagCorpus)).(io.Reader))
	fmt.Printf("Loaded %d instances.\nThe first one:%v\nThe last one:%v\n",
		len(corpus), corpus[0], corpus[len(corpus)-1])

	if len(corpus) != 5376 {
		t.Errorf("Expecting %d instances, got %d", 5376, len(corpus))
	}
}

func TestGenerateSelectedCorpus(t *testing.T) {
	*flagCSV = buildnrun.Pkg(path.Join(kCSVDir, kCSVFile))
	*flagCorpus = "/tmp/selected_linkedin_employee_exps_corpus.json"
	Run(new(selectedMembers))
}
