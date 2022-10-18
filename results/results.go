package results

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-msvc/utils/stringutils"
)

type Results interface {
	Add(code int, name, doc string) Result //must be uniq name and code
	AllResults() []Result
	ResultByName(name string) Result
	ResultByCode(code int) Result
}

type Result interface {
	error
	Results() Results //parent
	Code() int
	Name() string
	Doc() string
}

func New() (Results, Result) {
	rr := &results{
		results:      []result{},
		resultByCode: map[int]Result{},
		resultByName: map[string]Result{},
	}
	resultSuccess := rr.Add(0, "SUCCESS", "Success")
	return rr, resultSuccess
}

type results struct {
	results      []result
	resultByCode map[int]Result
	resultByName map[string]Result
}

//code must be 1..99 (0=success)
//name must be SNAKE_CODE
//doc is human readable text to describe the result
func (rr *results) Add(code int, name, doc string) Result {
	if code < 0 || code > 99 {
		panic(fmt.Sprintf("result code %d not 1..99 for %s", code, name)) //0=reserved for success
	}
	if name == "" || strings.ToUpper(name) == "" || !stringutils.IsSnakeCase(strings.ToLower(name)) {
		panic(fmt.Sprintf("result name \"%s\" is not uppercase snake (valid example: FAILED_TO_ADD)", name))
	}
	doc = strings.TrimSpace(doc)
	if doc == "" {
		panic(fmt.Sprintf("missing doc for result %d,%s", code, name))
	}
	if existingResult, ok := rr.resultByCode[code]; ok {
		panic(fmt.Sprintf("duplicate code %d for %s and %s", code, name, existingResult.Name()))
	}
	if existingResult, ok := rr.resultByName[name]; ok {
		panic(fmt.Sprintf("duplicate name %s for code %d and %d", name, code, existingResult.Code()))
	}
	r := result{
		results: rr,
		code:    code,
		name:    name,
		doc:     doc,
	}
	rr.results = append(rr.results, r)
	sort.Slice(rr.results, func(i, j int) bool { return rr.results[i].code < rr.results[j].code })
	rr.resultByCode[r.code] = r
	rr.resultByName[r.name] = r
	return r
}

func (rr results) AllResults() []Result {
	l := make([]Result, len(rr.results))
	for i, r := range rr.results {
		l[i] = r
	}
	return l
}

func (rr results) ResultByName(name string) Result {
	return rr.resultByName[name]
}

func (rr results) ResultByCode(code int) Result {
	return rr.resultByCode[code]
}

type result struct {
	results Results
	code    int    //0..99
	name    string //UPPER_SNAKE
	doc     string
}

func (r result) Results() Results { return r.results }
func (r result) Code() int        { return r.code }
func (r result) Name() string     { return r.name }
func (r result) Doc() string      { return r.doc }
func (r result) Error() string    { return fmt.Sprintf("RESULT(%d:%s)", r.code, r.name) }
