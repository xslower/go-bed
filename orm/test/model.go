package test

//@db test
//@table T1:mytest_+tvplay{category}
type Mytest struct {
	Category   string
	Id         int
	Name       string
	CreateTime string
}

//@db test
//@table T2:php_test_+01{id%100}
type PhpTest struct {
	Id   int
	Name string
}

//@table T3:php_test_+01{crc32(name)%100}
type User struct {
	Id   int
	Name string
}
