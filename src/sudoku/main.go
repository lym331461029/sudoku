package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	inputfile  string = ".\\Input.json"
	outputfile string = "Output.json"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/sudoku", solveSuduku)

	router.Run(":9090")
}

type ResponseContent struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Sdks []*Sudoku `json:"result"`
}

func ResLogicError(errCode int, errMsg string, resCont *ResponseContent, c *gin.Context) {
	resCont.Code = errCode
	resCont.Msg = errMsg
	c.JSON(200, resCont)
}

func solveSuduku(c *gin.Context) {
	sdk := &Sudoku{}
	respCont := &ResponseContent{}
	respCont.Code = 0
	respCont.Msg = "success"

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(sdk)
	if err != nil {
		if err != nil {
			log.Println(err)
			ResLogicError(-1, err.Error(), respCont, c)
			return
		}
	}
	c.Request.Body.Close()

	rels := make(chan *Sudoku, 100)
	go func() {
		sdk.GenerateSudoku(rels)
		close(rels)
		log.Println("计算已经完成...")
	}()

	for relSudoku := range rels {
		respCont.Sdks = append(respCont.Sdks, relSudoku)
	}
	bytes, err := json.MarshalIndent(respCont, " ", " ")
	if err != nil {
		if err != nil {
			log.Println(err)
			ResLogicError(-1, err.Error(), respCont, c)
			return
		}
	}

	c.String(200, string(bytes))
}

/*
func main() {
	var sdk Sudoku
	sdk.ReadJsonInit(inputfile)

	rels := make(chan *Sudoku, 100)

	go func() {
		sdk.GenerateSudoku(rels)
		close(rels)
		fmt.Println("计算已经完成...")
	}()

	for relSudoku := range rels {
		encolder := json.NewEncoder(os.Stdout)
		encolder.Encode(relSudoku)
		//relSudoku.WriteJsonOut(outputfile)
	}
}
*/
