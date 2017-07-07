package builder

import (
	`io/ioutil`
	`strings`
)

func BuildHelper(dir, packname string) {
	path := dir + `/helper.go`
	content := strings.Replace(helper, `{PACKAGE}`, packname, -1)
	ioutil.WriteFile(path, []byte(content), 0600)

}
