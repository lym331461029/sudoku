package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	Start1 int8 = 1
	End1   int8 = 3
	Start2 int8 = 5
	End2   int8 = 7
)

type Point struct {
	X int8
	Y int8
}

type Sudoku struct {
	Input       [9][9]SudokuElem
	ProblemType byte
	AreaMap     map[int8][]Point
}

func (s Sudoku) String() string {
	var str string
	str = fmt.Sprint("\n")
	for i := 0; i < len(s.Input); i++ {
		str = str + fmt.Sprintf("%v\n", s.Input[i])
	}
	str = str + "\n"
	return str
}

func (suduku *Sudoku) UnmarshalJSON(data []byte) error {
	ms := &struct {
		Input       [9][9]int `json:"input"`
		ProblemType string    `json:"problemType"`
	}{}
	err := json.Unmarshal(data, ms)
	if err != nil {
		return err
	}

	if ms.ProblemType != "A" {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				suduku.Input[i][j].SetValue(int8(ms.Input[i][j]))
				if suduku.Input[i][j].GetValue() == 0 {
					suduku.Input[i][j].PushAllToCache()
				}
			}
		}
		suduku.ProblemType = ms.ProblemType[0]
	} else {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				areaNo := ms.Input[i][j] / 10
				value := ms.Input[i][j] % 10
				suduku.Input[i][j].SetValue(int8(value))
				suduku.Input[i][j].SetArea(int8(areaNo))
				if suduku.Input[i][j].GetValue() == 0 {
					suduku.Input[i][j].PushAllToCache()
				}

				if suduku.AreaMap == nil {
					suduku.AreaMap = map[int8][]Point{}
				}

				if _, ok := suduku.AreaMap[int8(areaNo)]; ok {
					suduku.AreaMap[int8(areaNo)] = append(suduku.AreaMap[int8(areaNo)], Point{X: int8(i), Y: int8(j)})
				} else {
					suduku.AreaMap[int8(areaNo)] = []Point{Point{X: int8(i), Y: int8(j)}}
				}
			}
		}
		suduku.ProblemType = ms.ProblemType[0]
	}
	return nil
}

func (suduku Sudoku) MarshalJSON() ([]byte, error) {
	/*
		sudokuinit := struct {
			Input [9][9]uint8 `json:"data"`
			//ProblemType string      `json:"problemType"`
		}{}

		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				sudokuinit.Input[i][j] = uint8(suduku.Input[i][j].GetValue())
			}
		}
		//sudokuinit.ProblemType = string(suduku.ProblemType)
		return json.MarshalIndent(sudokuinit, "", "\t")
	*/

	buf := bytes.NewBuffer(nil)
	buf.WriteString("[\n")
	for i := 0; i < 9; i++ {
		buf.WriteString("\t\"")
		for j := 0; j <= 7; j++ {
			buf.WriteString(strconv.FormatInt(int64(suduku.Input[i][j].GetValue()), 10))
			buf.WriteString(",")
		}

		buf.WriteString(strconv.FormatInt(int64(suduku.Input[i][8].GetValue()), 10))
		if i < 8 {
			buf.WriteString("\",\n")
		} else {
			buf.WriteString("\"\n")
		}
	}
	buf.WriteString("]")

	log.Print(buf.String())
	return buf.Bytes(), nil

}

func (suduku *Sudoku) ReadJsonInit(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	bio := bufio.NewReader(f)

	data, err := ioutil.ReadAll(bio)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, suduku); err != nil {
		return err
	}

	return nil
}

func (suduku Sudoku) WriteJsonOut(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(suduku)
}

func generateRestrictFunc(XS, XE, YS, YE int8) func(suduku *Sudoku, x, y int8) {
	return func(suduku *Sudoku, x, y int8) {
		if x >= XS && x <= XE && y >= YS && y <= YE {
			for i := XS; i <= XE; i++ {
				for j := YS; j <= YE; j++ {
					if x == i && y == j {
						continue
					}
					if suduku.Input[i][j].GetValue() > 0 {
						suduku.Input[x][y].RemoveFromCache(suduku.Input[i][j].GetValue())
					}
				}
			}
		}
	}
}

//因为性能原因，该函数用来替换generateRestrictFunc
func (suduku *Sudoku) internalRestrict(XS, XE, YS, YE int8, x, y int8) {
	if x >= XS && x <= XE && y >= YS && y <= YE {
		for i := XS; i <= XE; i++ {
			for j := YS; j <= YE; j++ {
				if x == i && y == j {
					continue
				}
				if suduku.Input[i][j].GetValue() > 0 {
					suduku.Input[x][y].RemoveFromCache(suduku.Input[i][j].GetValue())
				}
			}
		}
	}
}

func (suduku *Sudoku) NineRestrict(x, y int8) int8 {
	var _XS int8 = x - x%3
	var _YS int8 = y - y%3

	//generateRestrictFunc(_XS, _XS+2, _YS, _YS+2)(suduku, x, y)
	suduku.internalRestrict(_XS, _XS+2, _YS, _YS+2, x, y)
	return suduku.Input[x][y].CacheNum()
}

func (suduku *Sudoku) RowRestrict(x, y int8) int8 {
	//generateRestrictFunc(x, x, 0, 8)(suduku, x, y)
	suduku.internalRestrict(x, x, 0, 8, x, y)
	return suduku.Input[x][y].CacheNum()
}

func (suduku *Sudoku) ColRestrict(x, y int8) int8 {
	//generateRestrictFunc(0, 8, y, y)(suduku, x, y)
	suduku.internalRestrict(0, 8, y, y, x, y)
	return suduku.Input[x][y].CacheNum()
}

func (suduku Sudoku) GetCandidateNum(x, y int8) int8 {
	if suduku.Input[x][y].GetValue() > 0 {
		return int8(1)
	}
	return suduku.Input[x][y].CacheNum()
}

//X数独限定
func (suduku *Sudoku) XRestrict(x, y int8) int8 {
	if x == y {
		for i := 0; i < 9; i++ {
			if x == int8(i) {
				continue
			}
			if suduku.Input[i][i].GetValue() > 0 {
				suduku.Input[x][y].RemoveFromCache(suduku.Input[i][i].GetValue())
			}
		}
	}

	if x+y == 8 {
		for i := 0; i < 9; i++ {
			if x == int8(i) {
				continue
			}
			if suduku.Input[i][8-i].GetValue() > 0 {
				suduku.Input[x][y].RemoveFromCache(suduku.Input[i][8-i].GetValue())
			}
		}
	}

	return suduku.Input[x][y].CacheNum()
}

//百分比数独限定
func (suduku *Sudoku) PercentumRestrict(x, y int8) int8 {
	//generateRestrictFunc(Start1, End1, Start1, End1)(suduku, x, y)
	//generateRestrictFunc(Start2, End2, Start2, End2)(suduku, x, y)
	suduku.internalRestrict(Start1, End1, Start1, End1, x, y)
	suduku.internalRestrict(Start2, End2, Start2, End2, x, y)

	if x+y == 8 {
		for i := 0; i < 9; i++ {
			if x == int8(i) {
				continue
			}
			if suduku.Input[i][8-i].GetValue() > 0 {
				suduku.Input[x][y].RemoveFromCache(suduku.Input[i][8-i].GetValue())
			}
		}
	}
	return suduku.Input[x][y].CacheNum()
}

//超数独限定
func (suduku *Sudoku) SuperRestrict(x, y int8) int8 {
	//generateRestrictFunc(Start1, End1, Start1, End1)(suduku, x, y)
	//generateRestrictFunc(Start2, End2, Start2, End2)(suduku, x, y)
	//generateRestrictFunc(Start1, End1, Start2, End2)(suduku, x, y)
	//generateRestrictFunc(Start2, End2, Start1, End1)(suduku, x, y)

	suduku.internalRestrict(Start1, End1, Start1, End1, x, y)
	suduku.internalRestrict(Start2, End2, Start2, End2, x, y)
	suduku.internalRestrict(Start1, End1, Start2, End2, x, y)
	suduku.internalRestrict(Start2, End2, Start1, End1, x, y)
	return suduku.Input[x][y].CacheNum()
}

//颜色数独限定
func (suduku *Sudoku) ColorRestrict(x, y int8) int8 {
	var _TpX, _TpY, i, j int8
	_TpX, _TpY = x%3, y%3

	for i = 0; i < 9; i += 3 {
		for j = 0; j < 9; j += 3 {
			if i+_TpX == x && j+_TpY == y {
				continue
			}

			if suduku.Input[i+_TpX][j+_TpY].GetValue() > 0 {
				suduku.Input[x][y].RemoveFromCache(
					suduku.Input[i+_TpX][j+_TpY].GetValue())
			}
		}
	}
	return suduku.Input[x][y].CacheNum()
}

//区域限制
func (suduku *Sudoku) AreaRestrict(x, y int8) int8 {
	an := suduku.Input[x][y].GetArea()
	aresPoints := suduku.AreaMap[an]

	for _, p := range aresPoints {
		if suduku.Input[p.X][p.Y].GetValue() > 0 {
			suduku.Input[x][y].RemoveFromCache(suduku.Input[p.X][p.Y].GetValue())
		}
	}
	return suduku.Input[x][y].CacheNum()
}

func (suduku *Sudoku) GenerateSudoku(rels chan *Sudoku) bool {
	var MinX, MinY, MinC int8
	var MaxC int8
	var next bool = true
	var tpCand int8

	for next {
		next = false
		MinC = 9
		MaxC = 1
		var i, j int8
		for i = 0; i < 9; i++ {
			for j = 0; j < 9; j++ {
				var CandidateNum int8 = suduku.GetCandidateNum(i, j)
				if CandidateNum > 1 {
					var tpNum int8
					tpNum = suduku.NineRestrict(i, j)
					if tpNum > 0 {
						tpNum = suduku.RowRestrict(i, j)
					}
					if tpNum > 0 {
						tpNum = suduku.ColRestrict(i, j)
					}
					if tpNum > 0 && suduku.ProblemType == 'X' {
						tpNum = suduku.XRestrict(i, j)
					}
					if tpNum > 0 && suduku.ProblemType == 'U' {
						tpNum = suduku.SuperRestrict(i, j)
					}

					if tpNum > 0 && suduku.ProblemType == 'P' {
						tpNum = suduku.PercentumRestrict(i, j)
					}

					if tpNum > 0 && suduku.ProblemType == 'C' {
						tpNum = suduku.ColorRestrict(i, j)
					}

					if tpNum > 0 && suduku.ProblemType == 'A' {
						tpNum = suduku.AreaRestrict(i, j)
					}

					if tpNum == 0 {
						return false
					}

					if tpNum == 1 {
						suduku.Input[i][j].SetValue(suduku.Input[i][j].PopCacheFront())
					}

					tpCand = suduku.GetCandidateNum(i, j)
					if tpCand < CandidateNum {
						next = true
					}

					if tpCand < MinC && tpCand > 1 {
						MinC = tpCand
						MinX = i
						MinY = j
					}

					if tpCand > MaxC {
						MaxC = tpCand
					}
				}
			}
		}
	}

	if MaxC > 1 {
		for suduku.Input[MinX][MinY].CacheNum() > 0 {
			_TpSudoku := *suduku
			_TpSudoku.Input[MinX][MinY].SetValue(_TpSudoku.Input[MinX][MinY].PopCacheFront())
			_TpSudoku.Input[MinX][MinY].RemoveAllCache()

			//problemes <- &_TpSudoku
			_TpSudoku.GenerateSudoku(rels)
			suduku.Input[MinX][MinY].PopCacheFront()
		}
	} else if MaxC == 1 && MinC == 9 {
		rels <- suduku
		return true
	}
	return false
}
