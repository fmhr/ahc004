package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

var sc = bufio.NewScanner(os.Stdin)
var buff []byte

func nextInt() int {
	sc.Scan()
	i, err := strconv.Atoi(sc.Text())
	if err != nil {
		panic(err)
	}
	return i
}
func nextString() string {
	sc.Scan()
	return sc.Text()
}

func init() {
	sc.Split(bufio.ScanWords)
	sc.Buffer(buff, bufio.MaxScanTokenSize*1024)
	log.SetFlags(log.Lshortfile)
}

// timer
var StartTime time.Time

// https://golang.org/pkg/runtime/pprof/
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

var local = flag.Bool("local", false, "if local")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...
	StartTime = time.Now()
	input()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

type Matrix [20][20]byte

func (m *Matrix) Reset() {
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			m[i][j] = '.'
		}
	}
}

func (m *Matrix) SetRow(sub string, y, x, d int) {
	n := len(sub)
	var nx int = x
	for i := 0; i < n; i++ {
		m[y][nx] = sub[i]
		nx++
		if nx > 19 {
			break
		}
	}
}

func (m Matrix) Output() (out [20]string) {
	for i := 0; i < 20; i++ {
		fmt.Println(string(m[i][:]))
	}
	return
}

// N は20に固定されている
func input() {
	N := nextInt()
	_ = N
	M := nextInt()
	s := make([]string, M)
	for i := 0; i < M; i++ {
		s[i] = nextString()
	}
	//solver(s)
	//solverB(s)
	solverB_random(s)
}

// 方針1 短い単語を優先して 横にいれて残りは偶然一致することを願う
// このとき横一列を回すとき、横に入れた単語は常に一致しているのでスコアは変更されない
func solver(sub []string) {
	sort.Slice(sub, func(i, j int) bool {
		return len(sub[i]) < len(sub[j])
	})
	var dna Matrix
	dna.Reset()
	y := 0
	x := 0
	for i := 0; i < len(sub); i++ {
		if x > 19 {
			x = 0
			y++
		}
		if y > 19 {
			break
		}
		dna.SetRow(sub[i], y, x, 0)
		x += len(sub[i])
	}
	str := dna.Output()
	for i := 0; i < 20; i++ {
		fmt.Println(str[i])
	}
	score := score(dna, sub)
	log.Printf("score=%d\n", int(score))
}

func searchPrefix(used []bool, str string, sub []string) string {
	if len(str) > 20 {
		return str
	}
	for i := 1; i < len(str); i++ {
		for j := 0; j < len(sub); j++ {
			if !used[j] {
				if strings.HasPrefix(sub[j], str[i:]) {
					if len(str)+len(sub[j][len(str)-i:]) > 20 {
						continue
					}
					used[j] = true
					if !strings.Contains(str, sub[j]) {
						// log.Println(i, str, sub[j])
						// log.Println(len(str)-i, sub[j][len(str)-i:])
						str += sub[j][len(str)-i:]
					}
					return searchPrefix(used, str, sub)
				}
			}
		}
	}
	return str
}

func usedCheck(sub []string, dna Matrix, used []bool, y int) []bool {
	for i := 0; i < y; i++ {
		str := string(dna[i][:]) + string(dna[i][:])
		for j := 0; j < len(sub); j++ {
			if !used[j] {
				if strings.Contains(str, sub[j]) {
					used[j] = true
				}
			}
		}
	}
	for i := 0; i < 20; i++ {
		var column [20]byte
		for j := 0; j < 20; j++ {
			column[j] = dna[j][i]
		}
		s := string(column[:]) + string(column[:])
		for j := 0; j < len(sub); j++ {
			if !used[j] {
				if strings.Contains(s, sub[j]) {
					used[j] = true
				}
			}
		}
	}

	return used
}

func solverB(sub []string) {
	var dna Matrix
	//sort.Strings(sub)
	used := make([]bool, len(sub))
	x := 0
	str := ""
	for i := 0; i < 20; i++ {
		for used[x] {
			x = rand.Intn(len(sub))
		}
		str = sub[x]
		used[x] = true
		s := searchPrefix(used, str, sub)
		for j := 0; j < 20; j++ {
			if len(s) > j {
				dna[i][j] = s[j]
			} else {
				dna[i][j] = '.'
			}
		}
	}
	dna.Output()
	score := score(dna, sub)
	log.Printf("score=%d\n", int(score))
}

func solverB_solo(sub []string) (int, Matrix) {
	ABCDEFGH := "ABCDEFGH"
	var dna Matrix
	//sort.Strings(sub)
	used := make([]bool, len(sub))
	x := 0
	str := ""
	for i := 0; i < 20; i++ {
		if i > 5 && i%2 == 1 {
			used = usedCheck(sub, dna, used, i)
		}
		for used[x] {
			x = rand.Intn(len(sub))
		}
		str = sub[x]
		used[x] = true
		s := searchPrefix(used, str, sub)
		for j := 0; j < 20; j++ {
			if len(s) > j {
				dna[i][j] = s[j]
			} else {
				dna[i][j] = ABCDEFGH[rand.Intn(8)]
			}
		}
	}
	return score(dna, sub), dna
}

func solverB_random(sub []string) {
	var timeLimit time.Duration = 2800
	var timeout bool

	go func() {
		time.Sleep(timeLimit * time.Millisecond)
		timeout = true
	}()
	var loop int
	max_score := 0
	var max_dna Matrix
	for !timeout {
		loop++
		score, dna := solverB_solo(sub)
		if score > max_score {
			max_score = score
			max_dna = dna
		}
	}
	max_dna.Output()
	exTime := time.Since(StartTime)
	log.Printf("score=%d loop=%d time=%v\n", score(max_dna, sub), loop, exTime)
	//solver_vertical_slot(sub, max_dna)
}

func solver_vertical_slot(sub []string, dna Matrix) {
	max_score := score(dna, sub)
	var max_dna Matrix
	max_dna = dna
	loop := 0
	var timeLimit time.Duration = 200
	var timeout bool
	go func() {
		time.Sleep(timeLimit * time.Millisecond)
		timeout = true
	}()
	for !timeout {
		loop++
		y := rand.Intn(20)
		s := score_right_shift(dna, sub, y)
		if s >= max_score {
			max_score = s
			v0 := dna[y][19]
			for i := 19; i > 0; i-- {
				dna[y][i] = dna[y][i-1]
			}
			dna[y][0] = v0
			max_dna = dna
		}
	}
	max_dna.Output()
	log.Printf("score=%d loop=%d\n", max_score, loop)
}

func score_right_shift(dna Matrix, sub []string, y int) int {
	hit := make([]bool, len(sub))
	hits := 0
	M := len(sub)
	for i := 0; i < 20; i++ {
		s := string(dna[i][:]) + string(dna[i][:])
		for j := 0; j < len(sub); j++ {
			if !hit[j] {
				if strings.Contains(s, sub[j]) {
					hit[j] = true
					hits++
				}
			}
		}
	}

	for i := 0; i < 20; i++ {
		var column [20]byte
		for j := 0; j < 20; j++ {
			i2 := i
			if j == y {
				i2 = (i + 1) % 20
			}
			column[j] = dna[j][i2]
		}
		s := string(column[:]) + string(column[:])
		for j := 0; j < len(sub); j++ {
			if !hit[j] {
				if strings.Contains(s, sub[j]) {
					hit[j] = true
					hits++
				}
			}
		}
	}
	score := 0.0
	if hits == M {
		cntDot := 0.0
		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				if dna[i][j] == '.' {
					cntDot++
				}
			}
		}
		score = math.Round((1e8 * (2 * 20 * 20 / (2*20*20 - cntDot))))
	} else {
		//log.Println(hits, M)
		score = math.Round(1e8 * (float64(hits) / float64(M)))
	}
	//log.Printf("score=%d\n", int(score))
	return int(score)
}

func score(dna Matrix, sub []string) int {
	hit := make([]bool, len(sub))
	hits := 0
	M := len(sub)
	for i := 0; i < 20; i++ {
		s := string(dna[i][:]) + string(dna[i][:])
		for j := 0; j < len(sub); j++ {
			if !hit[j] {
				if strings.Contains(s, sub[j]) {
					hit[j] = true
					hits++
				}
			}
		}
	}

	for i := 0; i < 20; i++ {
		var column [20]byte
		for j := 0; j < 20; j++ {
			column[j] = dna[j][i]
		}
		s := string(column[:]) + string(column[:])
		for j := 0; j < len(sub); j++ {
			if !hit[j] {
				if strings.Contains(s, sub[j]) {
					hit[j] = true
					hits++
				}
			}
		}
	}
	score := 0.0
	if hits == M {
		cntDot := 0.0
		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				if dna[i][j] == '.' {
					cntDot++
				}
			}
		}
		score = math.Round((1e8 * (2 * 20 * 20 / (2*20*20 - cntDot))))
	} else {
		//log.Println(hits, M)
		score = math.Round(1e8 * (float64(hits) / float64(M)))
	}
	//log.Printf("score=%d\n", int(score))
	return int(score)
}
