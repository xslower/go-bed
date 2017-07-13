package test

import (
	// "io/ioutil"
	"strconv"
	"testing"

	"github.com/xslower/go-die/orm"
)

var (
	gIds   = []int{}
	gMTest *MytestModel
)

func init() {
	err := ormInit(`config.json`)
	if err != nil {
		echo(err)
		return
	}
	gMTest = NewMytestModel()
}

func gIdsDel(i int) {
	gIds = append(gIds[:i], gIds[i+1:]...)
}

func TestInsert(t *testing.T) {
	fm := map[string]interface{}{`name`: `this is TestInsert`}
	f := orm.NewFieldsWrite().InitWithMap(fm)
	id := gMTest.Insert(f)
	gIds = append(gIds, id)
}

func TestInsertOrUpdate(t *testing.T) {

}

func TestMultiInsert(t *testing.T) {
	fm := map[string]interface{}{`name`: `this is TestMultiInsert1`}
	f := orm.NewFieldsMultiInsert().InitWithMap(fm)
	f.AddSlice([]interface{}{`这是多条插入测试2`})
	id := gMTest.MultiInsert(f)
	gIds = append(gIds, id, id+1)
}

func TestUpdate(t *testing.T) {
	fm := map[string]interface{}{`name`: `修改后的名称`}
	f := orm.NewFieldsWrite().InitWithMap(fm)
	cond := orm.NewConds().Equal(`id`, gIds[0])
	num := gMTest.Update(f, cond)
	if num != 1 {
		t.Error(`update failed`)
	}
}

func TestDelete(t *testing.T) {
	cond := orm.NewConds().Equal(`id`, gIds[1])
	num := gMTest.Delete(cond)
	if num != 1 {
		t.Error(`delete failed`)
	} else {
		gIdsDel(1)
	}
}

func TestSelect(t *testing.T) {
	gMTest.Select(nil, nil, nil)
	tts := gMTest.Result()
	if len(tts) == 0 {
		t.Error(`select failed`)
	}
	for i, tt := range tts {
		echo(i, *tt)
	}
}

func TestSave(t *testing.T) {
	tt := &Mytest{Name: `this is TestSave`}
	id := gMTest.Save(tt)
	gIds = append(gIds, id)
}

func TestSaveMulti(t *testing.T) {
	tts := []*Mytest{}
	for i := 0; i < 3; i++ {
		name := `this is TestSaveMulti X` + strconv.Itoa(i+1)
		tt := &Mytest{Name: name}
		tts = append(tts, tt)
	}
	id := gMTest.SaveMulti(gMTest.ToIRows(&tts))
	gIds = append(gIds, id, id+1, id+2)
}

func TestSaveMultiEmpty(t *testing.T) {
	tts := []*Mytest{}
	id := gMTest.SaveMulti(gMTest.ToIRows(&tts))
	echo(id)
}

func TestRemove(t *testing.T) {
	idx := len(gIds) - 1
	tt := &Mytest{Id: gIds[idx]}
	num := gMTest.Remove(tt)
	if num != 1 {
		t.Error(`remove failed`)
	} else {
		gIdsDel(idx)
	}
}

func TestRead(t *testing.T) {
	tt := &Mytest{Id: gIds[0]}
	gMTest.Read(tt)
	if tt.Name == `` {
		t.Error(`read failed`)
	}
	echo(tt)
}

func TestReadMulti(t *testing.T) {
	tts := &([]*Mytest{})
	gMTest.ReadMulti(gMTest.ToIRows(tts), nil)
	if len(*tts) == 0 {
		t.Error(`read multi failed`)
	}
	// tts = gMTest.Result()
	for i, tt := range *tts {
		echo(i, *tt)
	}
}

func TestTableName(t *testing.T) {

}
