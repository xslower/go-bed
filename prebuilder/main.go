package main

import (
	`flag`
	`fmt`
	`greentea/prebuilder/builder`
	`greentea/utils`
	`os`
)

var (
	file       = flag.String(`d`, ``, `-d /project/dir/`)
	skip_model = flag.Bool(`skip`, false, `-skip`)
	model_name = flag.String(`m`, `model.go`, `-m model.go`)
	pack_name  = flag.String(`p`, `main`, `-p main`)
)

func main() {
	flag.Parse()
	dir := *file
	if dir == `` {
		fmt.Println(`must specify the project directory. usage: -d \project\dir\`)
		return
	}
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		fmt.Println(`the dir is not a directory or path is error`)
		return
	}
	if !*skip_model {
		model := dir + `/` + *model_name
		if !utils.PathExist(model) {
			fmt.Println(`the model file path:[` + model + `] is not exist`)
			return
		}
		builder.ParseFile(model)
	}
	builder.BuildHelper(dir, *pack_name)
}
