package builder

import (
	// `bufio`
	// `io`
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/xslower/go-die/orm"
)

var (
	regxStt    *regexp.Regexp
	regxFld    *regexp.Regexp
	regxPkg    *regexp.Regexp
	regxDb     *regexp.Regexp
	regxTblBlk *regexp.Regexp
	regxTblDef *regexp.Regexp

	outputSlice   = []string{}
	registerModel string
)

func init() {
	// ptnAll := `[.\s` + "\n\r" + `]*`
	regxStt = regexp.MustCompile("\n(?://@db (.*)\n)?(?://@table (.*)\n)?" + `type\s*([^\s]*)\s*struct\s*{(?s)([^\}]*)}`)
	regxFld = regexp.MustCompile("\n\t([^\\s]+)\\s+([^\\s]+)(\\s*[`\"](.*)[`\"])?")
	regxPkg = regexp.MustCompile(`package ([^\s]*)`)
	// regxDb = regexp.MustCompile(`@db ([^\s]*)`)
	// regxTblBlk = regexp.MustCompile(`@table ([^\s]*)`)
	regxTblDef = regexp.MustCompile(`T(\d*):(.+)`)
	registerModel = ``
}

func ParseHeader(fileHeader string) {
	pkg := regxPkg.FindStringSubmatch(fileHeader)
	outputSlice = append(outputSlice, `package `+pkg[1], importDefine)
}

func ParseFile(file string) {
	content, _ := ioutil.ReadFile(file)
	// echo(string(content))
	sttDef := regxStt.FindAllStringSubmatch(string(content), -1)
	// echoStrSlice(sttDef...)
	// sttBody := []string{}
	// return
	sttNum := len(sttDef)
	if sttNum < 1 {
		fmt.Println(`There is no struct definition in the file:` + file)
		return
	}
	ParseHeader(string(content))

	for i := 0; i < sttNum; i++ {
		mRow := sttDef[i][3]
		AssembleModel(mRow, sttDef[i][1], sttDef[i][2])
		ImplementIRow(mRow, sttDef[i][4])

	}
	initStr := strings.Replace(initDefine, `{MODEL}`, registerModel, 1)
	outputSlice = append(outputSlice, initStr)

	WriteToFile(file)

}

func ParseTableDef(tblDef string) (string, string) {
	tblDef = strings.Trim(tblDef, " \n\r")
	if tblDef[0] != 'T' {
		return `return "` + tblDef + `"`, ``
	}
	typ := string(tblDef[1])
	tbnDefault := tblDef[3:]
	tbn_define := tbnDefault
	fields := []string{}
	defVals := []string{}
	partNum := ``
	code := tbnCodeDefine
	mode := ``
	left := strings.Index(tbnDefault, `{`)
	deal := func() {
		right := strings.Index(tbnDefault, `}`)
		fldDef := tbnDefault[left+1 : right]
		arrFldDef := strings.Split(fldDef, `:`)
		field := arrFldDef[0]
		value := arrFldDef[1]
		if typ == `2` {
			part_def := strings.Split(arrFldDef[0], `%`)
			field = part_def[0]
			partNum = part_def[1]
		}
		fields = append(fields, field)
		defVals = append(defVals, arrFldDef[1])
		tbn_define = strings.Replace(tbn_define, fldDef, field, -1)
		tbnDefault = strings.Replace(tbnDefault, `{`+fldDef+`}`, value, -1)
		left = strings.Index(tbnDefault, `{`)
	}
	if typ == `1` {
		mode = t1Define
		for left >= 0 {
			deal()
		}
	} else if typ == `2` {
		mode = t2Define
		deal()
	}

	partkey := strings.Join(fields, `", "`)
	strVal := strings.Join(defVals, `", "`)

	code = strings.Replace(code, `{MODE}`, mode, -1)
	code = strings.Replace(code, `{TBN_DEFAULT}`, tbnDefault, -1)
	code = strings.Replace(code, `{TBN_DEFINE}`, tbn_define, -1)
	code = strings.Replace(code, `{VAL_DEFAULT}`, strVal, -1)
	code = strings.Replace(code, `{PART_NUMBER}`, partNum, -1)
	return code, partkey
}

func AssembleModel(sttName, dbDef, tblDef string) {
	sttRowName := sttName
	sttModelName := sttRowName + `Model`
	dbName := ``        //default is empty
	if len(dbDef) > 1 { //have find the db define
		dbName = strings.Trim(dbDef, "\n\r")
	}
	tblCode := ``
	partKey := ``
	if len(tblDef) > 1 {
		tblCode, partKey = ParseTableDef(tblDef)
	} else { //did not find the table name define
		tblCode = orm.ToUnderline(sttRowName)
	}
	modelCode := subModelDefine
	modelCode = strings.Replace(modelCode, `{MODEL}`, sttModelName, -1)
	modelCode = strings.Replace(modelCode, `{ROW}`, sttRowName, -1)
	modelCode = strings.Replace(modelCode, `{DB}`, dbName, -1)
	modelCode = strings.Replace(modelCode, `{CODING}`, tblCode, -1)
	modelCode = strings.Replace(modelCode, `{PARTKEY}`, partKey, -1)
	// echo(modelCode)
	registerModel += `orm.RegisterModel(New` + sttModelName + `())` + "\n"

	outputSlice = append(outputSlice, modelCode)
}

func ImplementIRow(sttName, sttContent string) {
	fieldsDefine := regxFld.FindAllStringSubmatch(sttContent, -1)

	fSetCase := ``
	fGetCase := ``
	columns := []string{}
	for _, fd := range fieldsDefine {
		fname := fd[1]
		if fname == `//` {
			continue
		}
		ftype := fd[2]
		fnameInDb := orm.ToUnderline(fname)
		columns = append(columns, fnameInDb)
		if fd[4] != `` {
			fnameInDb = fd[4]
		}
		fSetCase += byteToTypeInCode(fnameInDb, fname, ftype)
		fGetCase += `case "` + fnameInDb + `":` + "\n" + `return this.` + fname + "\n"
	}
	setterStr := strings.Replace(rowSetterDefine, `{CASE}`, fSetCase, -1)
	setterStr = strings.Replace(setterStr, `{ROW}`, sttName, -1)
	getterStr := strings.Replace(rowGetterDefine, `{CASE}`, fGetCase, -1)
	getterStr = strings.Replace(getterStr, `{ROW}`, sttName, -1)
	columnStr := strings.Replace(rowColumnsDefine, `{COLUMNS}`, "`"+strings.Join(columns, "`,`")+"`", -1)
	columnStr = strings.Replace(columnStr, `{ROW}`, sttName, -1)

	outputSlice = append(outputSlice, setterStr, getterStr, columnStr)
}

func WriteToFile(file string) {
	targetFile := strings.Replace(file, `.go`, `_implements.go`, 1)
	content := []byte(strings.Join(outputSlice, "\n"))
	ioutil.WriteFile(targetFile, content, 0600)

	gofmt := exec.Command(`gofmt`, targetFile)
	gofmt.Stdout, _ = os.OpenFile(targetFile, os.O_WRONLY, 0)
	gofmt.Run()
}

func byteToTypeInCode(nameInDb, fdName, typ string) string {
	caseStr := `case "` + nameInDb + `":` + "\n"
	assignValue := ``
	if typ == `int` {
		assignValue = `this.` + fdName + `, err = strconv.Atoi(string(val))`
	} else if strings.HasPrefix(typ, `int`) {
		assignValue = `tmpVal, err := strconv.ParseInt(string(val), 10, ` + strings.TrimPrefix(typ, `int`) + `)` + "\n"
		assignValue += `this.` + fdName + ` = ` + typ + `(tmpVal)`
	} else if typ == `uint` {
		assignValue = `tmpVal, err := strconv.ParseUint(string(val), 10, 0)` + "\n"
		assignValue += `this.` + fdName + ` = ` + typ + `(tmpVal)`
	} else if strings.HasPrefix(typ, `int`) {
		assignValue = `tmpVal, err := strconv.ParseUint(string(val), 10, ` + strings.TrimPrefix(typ, `uint`) + `)` + "\n"
		assignValue += `this.` + fdName + ` = ` + typ + `(tmpVal)`
	} else if typ == `[]byte` {
		assignValue = `this.` + fdName + ` = val`
	} else if typ == `string` {
		assignValue = `this.` + fdName + ` = string(val)`
	} else if typ == `time.Time` {
		//not suuport yet
	}
	return caseStr + assignValue + "\n"
}
