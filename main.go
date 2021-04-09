package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
	r.Static("/wwwroot/css", APP_PATH+"wwwroot/css")
	r.Static("/wwwroot/js", APP_PATH+"wwwroot/js")
	r.Static("/wwwroot/picture", APP_PATH+"wwwroot/picture")

	r.GET("/expression", func(c *gin.Context) {
		typeGene := c.Query("typeGene")

		humanResult := []map[string]interface{}{}
		mouseResultLncRNA := [][]map[string]interface{}{}
		headersRnaLncRNA := retrieveHeader("LncRNA Expression Counts")
		mouseInfoLncRna := []map[string]interface{}{}
		mouseResultMRNA := [][]map[string]interface{}{}
		headersRnaMRna := retrieveHeader("mRNA Expression Counts")
		mouseInfoMRna := []map[string]interface{}{}
		mouseInfoTopLncRna := false
		mouseInfoTopMRna := false
		if typeGene != "" {
			humanResult = getTypeHumanResult(typeGene)
			mouseResultLncRNA = retrieveRows("LncRNA Expression Counts", typeGene, headersRnaLncRNA)
			mouseInfoLncRna = getMouseInfo(typeGene, "LncRNA")
			mouseResultMRNA = retrieveRows("mRNA Expression Counts", typeGene, headersRnaMRna)
			mouseInfoMRna = getMouseInfo(typeGene, "mRNA")
		}

		rowsType := [][]interface{}{}
		for i := 0; i < len(mouseResultLncRNA); i++ {
			var result []interface{}
			result = append(result, "LncRNA")
			if len(mouseInfoLncRna) > 0 {
				result = append(result, mouseInfoLncRna[0]["symbol"])
				result = append(result, mouseInfoLncRna[0]["log_fc"])
				result = append(result, mouseInfoLncRna[0]["p_value"])
				result = append(result, mouseInfoLncRna[0]["adj_p_value"])
				result = append(result, mouseInfoLncRna[0]["name"])

				if mouseInfoLncRna[0]["top"].(int64) == 1 {
					mouseInfoTopLncRna = true
				}
			} else {
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
			}

			for j := 0; j < len(mouseResultLncRNA[i]); j++ {
				item := mouseResultLncRNA[i][j]["value"]
				result = append(result, item)
			}

			rowsType = append(rowsType, result)
		}

		for i := 0; i < len(mouseResultMRNA); i++ {
			var result []interface{}
			result = append(result, "mRNA")
			if len(mouseInfoMRna) > 0 {
				result = append(result, mouseInfoMRna[0]["symbol"])
				result = append(result, mouseInfoMRna[0]["log_fc"])
				result = append(result, mouseInfoMRna[0]["p_value"])
				result = append(result, mouseInfoMRna[0]["adj_p_value"])
				result = append(result, mouseInfoMRna[0]["name"])

				if mouseInfoMRna[0]["top"].(int64) == 1 {
					mouseInfoTopMRna = true
				}
			} else {
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
			}

			for j := 0; j < len(mouseResultMRNA[i]); j++ {
				item := mouseResultMRNA[i][j]["value"]
				result = append(result, item)
			}

			rowsType = append(rowsType, result)
		}

		// Human
		sheets := getHumanSheets()

		cancerType := c.Query("selectedCancerType")
		expression := c.Query("selectedExpressionList")

		expressionList := []map[string]interface{}{}
		expressionResult := []map[string]interface{}{}
		displayExpression := false
		if cancerType != "" && cancerType != "None" {
			expressionList = getExpressions(cancerType)

			if expression != "" && expression != "None" {
				displayExpression = true
				expressionResult = getExpressionsResult(cancerType, expression)
			}
		}

		if (typeGene != "") && (len(humanResult) > 0) {
			displayExpression = true
			for i := 0; i < len(humanResult); i++ {
				expressionResult = append(expressionResult, humanResult[i])
			}
		}

		// Mouse
		lncExpressions := getLncExpressions()
		mRnaExpressions := getMRnaExpressions()

		queriedSheet := c.Query("selectedSheet")
		queriedRna := c.Query("selectedRna")

		rowsRna := [][]map[string]interface{}{}
		headersRna := []map[string]interface{}{}
		mouseInfoRna := []map[string]interface{}{}

		if queriedSheet == "LncRNA" {
			if queriedRna != "None" && queriedRna != "" {
				headersRna = retrieveHeader("LncRNA Expression Counts")
				rowsRna = retrieveRows("LncRNA Expression Counts", queriedRna, headersRna)
				mouseInfoRna = getMouseInfo(queriedRna, "LncRNA")
			}
		} else if queriedSheet == "mRNA" {
			if queriedRna != "None" && queriedRna != "" {
				headersRna = retrieveHeader("mRNA Expression Counts")
				rowsRna = retrieveRows("mRNA Expression Counts", queriedRna, headersRna)
				mouseInfoRna = getMouseInfo(queriedRna, "mRNA")
			}
		}

		var rows [][]interface{}
		for i := 0; i < len(rowsRna); i++ {
			var result []interface{}
			result = append(result, queriedSheet)
			if len(mouseInfoRna) > 0 {
				result = append(result, mouseInfoRna[0]["symbol"])
				result = append(result, mouseInfoRna[0]["log_fc"])
				result = append(result, mouseInfoRna[0]["p_value"])
				result = append(result, mouseInfoRna[0]["adj_p_value"])
				result = append(result, mouseInfoRna[0]["name"])

				if queriedSheet == "LncRNA" {
					if mouseInfoRna[0]["top"].(int64) == 1 {
						mouseInfoTopLncRna = true
					}
				} else if queriedSheet == "mRNA" {
					if mouseInfoRna[0]["top"].(int64) == 1 {
						mouseInfoTopMRna = true
					}
				}

				fmt.Printf("%t %t\n", mouseInfoTopLncRna, mouseInfoTopMRna)
			} else {
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
				result = append(result, nil)
			}

			for j := 0; j < len(rowsRna[i]); j++ {
				item := rowsRna[i][j]["value"]
				result = append(result, item)
			}

			rows = append(rows, result)
		}

		if typeGene != "" && (len(rowsType) > 0) {
			for i := 0; i < len(rowsType); i++ {
				rows = append(rows, rowsType[i])
			}
		}

		c.HTML(http.StatusOK, "/views/index.tmpl", gin.H{
			"Sheets":            sheets,
			"DisplayExpression": displayExpression,
			"ExpressionList":    expressionList,
			"ExpressionResult":  expressionResult,

			"LncExpressions":  lncExpressions,
			"MRnaExpressions": mRnaExpressions,
			"TopLnc":          mouseInfoTopLncRna,
			"TopMRna":         mouseInfoTopMRna,
			"QueriedSheet":    queriedSheet,
			"QueriedRna":      queriedRna,
			"HeaderRna":       headersRna,
			"RowsRna":         rowsRna,
			"MouseData":       rows,
			"MouseInfoRna":    mouseInfoRna,
			"DisplayRnaOut":   (queriedRna != "None" && queriedRna != "") || (len(rows) > 0),
		})
	})

	r.POST("/expression", func(c *gin.Context) {
		selectedCancerType := c.PostForm("selectCancerType")
		selectedExpressionList := c.PostForm(("selectExpressionList"))

		q := url.Values{}
		q.Set("selectedCancerType", selectedCancerType)
		q.Set("selectedExpressionList", selectedExpressionList)

		selectedExpression := c.PostForm("selectExpression")
		selectedLnc := c.PostForm("selectLnc")
		selectedMRna := c.PostForm("selectMRna")

		q.Set("selectedSheet", selectedExpression)

		if selectedExpression == "LncRNA" {
			q.Set("selectedRna", selectedLnc)
		} else if selectedExpression == "mRNA" {
			q.Set("selectedRna", selectedMRna)
		}

		typeGene := c.PostForm("typeGene")
		q.Set("typeGene", typeGene)

		location := url.URL{Path: "/expression", RawQuery: q.Encode()}
		c.Redirect(http.StatusFound, location.RequestURI())
	})

	r.GET("/human", func(c *gin.Context) {

		sheets := getHumanSheets()

		sheet := c.Query("selectedSheet")
		expression := c.Query("selectedExpressionList")

		expressionList := []map[string]interface{}{}
		expressionResult := []map[string]interface{}{}
		displayExpression := false
		if sheet != "" && sheet != "None" {
			expressionList = getExpressions(sheet)

			if expression != "" && expression != "None" {
				displayExpression = true
				expressionResult = getExpressionsResult(sheet, expression)
			}
		}

		c.HTML(http.StatusOK, "/views/human.tmpl", gin.H{
			"Sheets":            sheets,
			"DisplayExpression": displayExpression,
			"ExpressionList":    expressionList,
			"ExpressionResult":  expressionResult,
		})
	})

	r.POST("/human", func(c *gin.Context) {
		selectedSheet := c.PostForm("selectSheet")
		selectedExpression := c.PostForm(("selectExpressionList"))

		q := url.Values{}
		q.Set("selectedSheet", selectedSheet)
		q.Set("selectedExpressionList", selectedExpression)

		location := url.URL{Path: "/human", RawQuery: q.Encode()}
		c.Redirect(http.StatusFound, location.RequestURI())
	})

	type HumanSheet struct {
		Sheet string `form:"sheet"`
	}

	r.GET("/human_sheet", func(c *gin.Context) {
		var humanSheet HumanSheet
		c.Bind(&humanSheet)

		q := url.Values{}
		q.Set("selectedSheet", humanSheet.Sheet)

		expressionList := []map[string]interface{}{}
		if humanSheet.Sheet != "" {
			expressionList = getExpressions(humanSheet.Sheet)
		}

		c.JSON(http.StatusOK, expressionList)
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/views/about.tmpl", gin.H{})
	})

	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/views/file_upload.tmpl", gin.H{})
	})

	r.POST("/upload", func(c *gin.Context) {
		err := populateRnaDb()
		if err != nil {
			fmt.Println(err)
		}

		err = populatePearsonDb()
		if err != nil {
			fmt.Println(err)
		}

		err = populateMouseInfo()
		if err != nil {
			fmt.Println(err)
		}

		err = populateHuman()
		if err != nil {
			fmt.Println(err)
		}

		c.Redirect(http.StatusFound, "/")
	})

	r.Run(HOST)
}
