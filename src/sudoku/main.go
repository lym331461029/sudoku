// sudokutest project main.go

///数独计算测试程序
package main

import (
	"time"
)

const (
	inputfile  string = ".\\Input.json"
	outputfile string = "Output.json"
)

func main() {
	var sdk Sudoku
	sdk.ReadJsonInit(inputfile)

	rels := make(chan *Sudoku, 100)
	problems := make(chan *Sudoku, 10)
	flagxx := make(chan bool)

	for i := 0; i < 4; i++ {
		go func() {
			for sudokuinit := range problems {
				sudokuinit.GenerateSudoku(rels, problems)
			}
		}()
	}

	go func() {
		finish := time.After(2 * time.Second)
		endflag := false
		for {
			if !endflag {
				select {
				case rel := <-rels:
					rel.WriteJsonOut(outputfile)
				case <-finish:
					endflag = true
					break
				}
			}
			if endflag {
				break
			}
		}
		flagxx <- true
	}()

	problems <- &sdk

	for _ = range flagxx {
		close(problems)
		close(rels)
		return
	}
}
