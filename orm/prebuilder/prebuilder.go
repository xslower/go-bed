package main

import (
	// `bufio`
	// `io`
	`greentea/orm`
	`io/ioutil`
	`os`
	`os/exec`
	`regexp`
	`strings`
)

var (
	importDefine = `
import (
	"errors"
	"strconv"
	"greentea/orm"
	"hash/crc32"
)
`
	subModelDefine = `
func New{MODEL}() *{MODEL} {
	mt := &{MODEL}{}
	mt.BaseModel.IModel = mt
	mt.result = []*{ROW}{}
	return mt
}

type {MODEL} struct {
	orm.BaseModel
	result []*{ROW}
}

func (this *{MODEL}) DbName() string {
	return "{DB}"
}
func (this *{MODEL}) TableName(elems ...orm.ISqlElem) string {
	{CODING}
}
func (this *{MODEL}) PartitionKey() string {
	return "{PARTKEY}"
}

func (this *{MODEL}) Result() []*{ROW}{
	return this.result
}
func (this *{MODEL}) CreateRow() orm.IRow {
	irow := &{ROW}{}
	this.result = append(this.result, irow)
	return irow
}
func (this *{MODEL}) ResetResult() {
	this.result = []*{ROW}{}
}`
	rowSetterDefine = `
func (this *{ROW}) Set(key string, val []byte) error {
	var err error
	switch key {
	{CASE}
	default:
		err = errors.New("No such column [" + key + "]")
	}
	return err
}`
	rowGetterDefine = `
func (this *{ROW}) Get(key string) (interface{}) {
	switch key {
	{CASE}
	default:	
		return nil
	}
}`
	rowColumnsDefine = `
func (this *{ROW}) Columns() []string {
	return []string{{COLUMNS}}
}`

	initDefine = `
func ormStart(dbConfig map[string]string) {
	orm.Start(dbConfig)
	{MODEL}
}
func packageHolder(){
	_ = crc32.ChecksumIEEE([]byte("a"))
}`
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
	regxStt = regexp.MustCompile(`type\s*([^\s]*)\s*struct`)
	regxFld = regexp.MustCompile("\t([^\\s]+)\\s+([^\\s]+)(\\s*[`\"](.*)[`\"])?")
	regxPkg = regexp.MustCompile(`package ([^\s]*)`)
	regxDb = regexp.MustCompile(`@db ([^\s]*)`)
	regxTblBlk = regexp.MustCompile(`@table ([^\s]*)`)
	regxTblDef = regexp.MustCompile(`T(\d*):(.*)\+(.*)\{(.*)\}`)
	registerModel = ``
}

func ParseHeader(fileHeader string) {
	pkg := regxPkg.FindStringSubmatch(fileHeader)
	outputSlice = append(outputSlice, `package `+pkg[1], importDefine)
}

func ParseTableDef(tblDef string) (string, string) {
	pre := `if len(elems) == 0 {
		return "{PREFIX}{DEFAULT}"
	}
	elem := elems[0]
	ifc, _ := elem.Get("{FIELD}")`
	t1 := `
	{FIELD}, _ := orm.InterfaceToString(ifc)
	if {FIELD} == "" {
		return "{PREFIX}{DEFAULT}"
	}
	return "{PREFIX}" + {FIELD}`
	t2 := `
	{FIELD}, _ := orm.InterfaceToInt(ifc)
	if {FIELD} == 0 {
		return "{PREFIX}{DEFAULT}"
	}
	part := {FIELD} % {PART_NUMBER}
	prefix := "{PREFIX}"
	if part < 10 {
		prefix += "0"
	}
	return prefix + strconv.FormatUint(uint64(part), 10)`
	t3 := `
	{FIELD}, _ := orm.InterfaceToString(ifc)
	if {FIELD} == "" {
		return "{PREFIX}{DEFAULT}"
	}
	hash := {HASH}
	part := hash % {PART_NUMBER}
	prefix := "{PREFIX}"
	if part < 10 {
		prefix += "0"
	}
	return prefix + strconv.FormatUint(uint64(part), 10)`

	fragment := regxTblDef.FindStringSubmatch(tblDef)
	if len(fragment) < 1 {
		return `return "` + tblDef + `"`, ``
	}
	typ := fragment[1]
	prefix := fragment[2]
	defVal := fragment[3]
	inner := fragment[4]
	field := inner
	part_number := ``
	hash := ``

	code := pre
	switch typ {
	case `1`:
		code += t1
	case `2`:
		code += t2
		pos := strings.Index(inner, `%`)
		field = inner[:pos]
		part_number = inner[pos+1:]
	case `3`:
		code += t3
		regx := regexp.MustCompile(`(.*)\((\w*)\)%(\d*)`)
		fgmt := regx.FindStringSubmatch(inner)
		field = fgmt[2]
		part_number = fgmt[3]
		if fgmt[1] == `crc32` {
			hash = `crc32.ChecksumIEEE([]byte(` + field + `))`
		}
	default:
		echo(`table partition define error`)
	}
	code = strings.Replace(code, `{FIELD}`, field, -1)
	code = strings.Replace(code, `{PREFIX}`, prefix, -1)
	code = strings.Replace(code, `{DEFAULT}`, defVal, -1)
	code = strings.Replace(code, `{PART_NUMBER}`, part_number, -1)
	code = strings.Replace(code, `{HASH}`, hash, -1)
	return code, field
}

func AssembleModel(sttName, sttHeader string) {
	sttRowName := sttName
	sttModelName := sttRowName + `Model`
	dbDef := regxDb.FindStringSubmatch(sttHeader)
	dbName := ``        //default is empty
	if len(dbDef) > 1 { //have find the db define
		dbName = dbDef[1]
	}
	tblDef := regxTblBlk.FindStringSubmatch(sttHeader)
	tblCode := ``
	partKey := ``
	if len(tblDef) > 1 {
		tblCode, partKey = ParseTableDef(tblDef[1])
	} else { //did not find the table name define
		tblCode = orm.ToUnderline(sttRowName)
	}
	modelCode := subModelDefine
	modelCode = strings.Replace(modelCode, `{MODEL}`, sttModelName, -1)
	modelCode = strings.Replace(modelCode, `{ROW}`, sttRowName, -1)
	modelCode = strings.Replace(modelCode, `{DB}`, dbName, -1)
	modelCode = strings.Replace(modelCode, `{CODING}`, tblCode, -1)
	modelCode = strings.Replace(modelCode, `{PARTKEY}`, partKey, -1)

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

func ParseFile(file string) {
	content, _ := ioutil.ReadFile(file)
	sttName := regxStt.FindAllStringSubmatch(string(content), -1)
	sttBody := regxStt.Split(string(content), -1)
	sttNum := len(sttBody)
	if sttNum < 2 {
		echo(`There is no struct definition in the file:` + file)
		return
	}
	ParseHeader(sttBody[0])

	for i := 1; i < len(sttBody); i++ {

		mRow := sttName[i-1][1]
		AssembleModel(mRow, sttBody[i-1])

		ImplementIRow(mRow, sttBody[i])

	}
	initStr := strings.Replace(initDefine, `{MODEL}`, registerModel, 1)
	outputSlice = append(outputSlice, initStr)

	WriteToFile(file)

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
