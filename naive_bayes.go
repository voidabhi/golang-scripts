package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	punc = regexp.MustCompile(`\W+`)
)

func main() {
	start := time.Now()
	counter := NewCounter()
	// data set for authorship identification
	//https://archive.ics.uci.edu/ml/datasets/Reuter_50_50#
	dir := "C50/C50train"
	out := readC50DataSet(dir, ExtractFeatures)
	rows := 0
	for row := range out {
		counter.Add(row.Class, row.Features)
		rows++
	}

	fmt.Printf("Data rows %d \nClasses %d \nFeatures %d \nTrain duration %+v\n",
		rows, len(counter.classes), len(counter.features), time.Since(start)) // output for debug

	start = time.Now()
	dir = "C50/C50test"
	out = readC50DataSet(dir, ExtractFeatures)
	countTotal := 0
	correct := 0
	for row := range out {
		countTotal++
		class := NaiveBayes(counter, row.Features)
		if row.Class == class {
			correct++
		}
	}
	fmt.Printf("Predictions %d \nTest duration %+v\n",
		countTotal, time.Since(start)) // output for debug

	fmt.Printf("Correct rate %f\n", float64(correct)*100/float64(countTotal)) // output for debug

}

type Row struct {
	Class    string
	Features []string
}

type ICounter interface {
	NumberOfTuples() int
	Pairs(class, feature string) (int, bool)
	Classes() map[string]int
}

func ExtractFeatures(src []byte) []string {
	// TODO rm most used 500 words (experiment)
	txt := string(punc.ReplaceAll(src, []byte(" ")))
	m := make(map[string]struct{})
	words := make([]string, 0, len(m))
	for _, str := range strings.Split(txt, " ") {
		if len(str) > 1 {
			if _, ok := m[str]; !ok {
				m[str] = struct{}{}
				words = append(words, str)
			}
		}
	}

	return words
}

// returns best candidate from sampled classes for given features
// Naive Bayesian Classification
func NaiveBayes(c ICounter, features []string) string {
	type candidate struct {
		class string
		prob  float64
	}
	var max candidate
	classes := len(c.Classes())
	out := make(chan candidate, classes)
	for class, classCount := range c.Classes() {
		go func(class string, classCount float64, tuples float64) {
			// max float64 to prevent overflow
			// posterior probability of features with given class
			pcf := 1.797693134862315708145274237317043567981e+308
			// prior probability
			pc := classCount / tuples
			for _, f := range features {
				if v, ok := c.Pairs(class, f); ok {
					if v > 0 {
						// TODO test as continuous-valued attr
						pcf *= float64(v) / classCount
					} else if v == 0 {
						// TODO use Laplacian correction for missing feature class tuples
						pcf *= 1.0 / classCount
					}
				}
			}
			out <- candidate{class, pc * pcf}
		}(class, float64(classCount), float64(c.NumberOfTuples()))
	}

	for can := range out {
		if max.prob < can.prob {
			max.class = can.class
			max.prob = can.prob
		}
		classes--
		if classes == 0 {
			break
		}
	}
	close(out)

	return max.class
}

type Counter struct {
	classes  map[string]int
	features map[string]map[string]int
}

func NewCounter() *Counter {
	return &Counter{
		classes:  make(map[string]int),
		features: make(map[string]map[string]int),
	}
}

func (c *Counter) Add(class string, features []string) {
	c.classes[class]++
	for _, f := range features {
		if _, ok := c.features[f]; ok {
			c.features[f][class]++
		} else {
			c.features[f] = make(map[string]int)
		}
	}
}

// returns counted (class, feature) tuples
func (c *Counter) Pairs(class, feature string) (int, bool) {
	if _, ok := c.features[feature]; !ok {
		return 0, false
	} else {
		return c.features[feature][class], true
	}
}

func (c *Counter) Classes() map[string]int {
	return c.classes
}

func (c *Counter) NumberOfTuples() int {
	count := 0
	for _, v := range c.classes {
		count += v
	}

	return count
}

func readC50DataSet(path string, procWords func([]byte) []string) <-chan Row {
	ch := make(chan Row)
	dir, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	dirs, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for i := range dirs {
		if dirs[i].IsDir() {
			wg.Add(1)
			go func(ch chan<- Row, dir string) {
				defer func() {
					wg.Done()
				}()

				var (
					f     *os.File
					files []os.FileInfo
					err   error
					src   []byte
				)

				if f, err = os.Open(dir); err == nil {
					if files, err = f.Readdir(-1); err == nil {
						if len(files) != 0 {
							for i := range files {
								if src, err = ioutil.ReadFile(dir + "/" + files[i].Name()); err == nil {
									ch <- Row{
										Class:    dir[strings.LastIndex(dir, "/")+1:],
										Features: procWords(src),
									}
								} else {
									break
								}
							}
						} else {
							err = errors.New("no files in dir")
						}
					}

				}

				if err != nil {
					panic(err)
				}
			}(ch, path+"/"+dirs[i].Name())
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
