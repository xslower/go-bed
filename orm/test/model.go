package test

//table T1:my{test:test}_{category:d}
//@db test
//@table mytest
type Mytest struct {
	Test       string
	Category   string
	Id         int
	Name       string
	CreateTime string
}

//@db test
//@table T2:php_test_{id%100:001}
// type PhpTest struct {
// 	Id   int
// 	Name string
// }

// //@table T3:php_test_+01{crc32(name)%100}
// type User struct {
// 	Id   int
// 	Name string
// }
