package converter

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const (
	kCSVDir  = "github.com/wangkuiyi/hmm/exp/linkedin_employee_data/data"
	kCSVFile = "pos_and_edu_for_LI_employees.csv"
)

func TestGenerateAndLoad(t *testing.T) {
	dir, e := ioutil.TempDir("", "converter_test")
	if e != nil {
		t.Fatalf("Cannot create temp dir: %v", e)
	}
	defer os.RemoveAll(dir)

	*flagCSV = path.Join(os.Getenv("GOPATH"), "src", kCSVDir, kCSVFile)
	*flagCorpus = path.Join(dir, "corpus.json")

	Run(new(PlainFeatureGenerator))

	corpus := LoadJSON(α(os.Open(*flagCorpus)).(io.Reader))
	fmt.Printf("Loaded %d instances.\nThe first one:%v\nThe last one:%v\n",
		len(corpus), corpus[0], corpus[len(corpus)-1])
}
