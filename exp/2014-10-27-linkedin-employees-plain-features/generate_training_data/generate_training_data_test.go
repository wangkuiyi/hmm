package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"testing"
)

const (
	kTestIn = "src/futureyou/linkedin_employees/generate_training_data" +
		"/testdata/test.csv"
)

func TestLoadCSV(t *testing.T) {
	truth := `{
  "43099587": [
  {
  "CompanyOrSchool": "145411",
  "PosOrEduRank": "-1",
  "TitleOrDegree": "1174",
  "SeniorityOrDegreeRank": "4",
  "FunctionOrField": "24",
  "Begin": "2012-05-01T00:00:00-07:00",
  "End": "2014-09-29T00:00:00-07:00",
  "IsJob": true,
  "IsEdu": false
  }
  ],
  "74749430": [
  {
  "CompanyOrSchool": "1294365",
  "PosOrEduRank": "2",
  "TitleOrDegree": "2498",
  "SeniorityOrDegreeRank": "4",
  "FunctionOrField": "4",
  "Begin": "2012-03-01T00:00:00-08:00",
  "End": "2013-05-01T00:00:00-07:00",
  "IsJob": true,
  "IsEdu": false
  },
  {
  "CompanyOrSchool": "82005",
  "PosOrEduRank": "4",
  "TitleOrDegree": "13337",
  "SeniorityOrDegreeRank": "4",
  "FunctionOrField": "25",
  "Begin": "2007-05-01T00:00:00-07:00",
  "End": "2008-12-01T00:00:00-08:00",
  "IsJob": true,
  "IsEdu": false
  }
  ]
  }`

	if f, e := os.Open(path.Join(os.Getenv("GOPATH"), kTestIn)); e != nil {
		t.Fatalf("Cannot open file: %v", e)
	} else {
		if m, e := LoadCSV(f); e != nil {
			t.Fatalf("LoadCSV: %v", e)
		} else {
			if b, e := json.MarshalIndent(m, "  ", ""); e != nil {
				t.Fatalf("json.MarshalIndent: %v", e)
			} else {
				if string(b) != truth {
					t.Errorf("Expecting %s, got %s", truth, string(b))
				}
			}
		}
	}
}

func TestGenerateJSON(t *testing.T) {
	const truth = `[{"Obs":[[{"company145411":1},{},{"title1174":1},{"seniority4":1},{"function24":1}],[{"company145411":1},{},{"title1174":1},{"seniority4":1},{"function24":1}],[{"company145411":1},{},{"title1174":1},{"seniority4":1},{"function24":1}]],"Periods":[1,1,1]}]
`

	if f, e := os.Open(path.Join(os.Getenv("GOPATH"), kTestIn)); e != nil {
		t.Fatalf("Cannot open file: %v", e)
	} else {
		if m, e := LoadCSV(f); e != nil {
			t.Fatalf("LoadCSV: %v", e)
		} else {
			var buf bytes.Buffer
			e := json.NewEncoder(&buf).Encode(GenerateJSON(m))
			if e != nil {
				t.Fatalf("JSON encoding: %v", e)
			} else {
				if buf.String() != truth {
					t.Errorf("Expecting %s, got %s", truth, buf.String())
				}
			}
		}
	}
}
