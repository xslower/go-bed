# GoORM
a high performance and lack use reflect and assertion golang framework.
now there is only have orm
orm just support normal CRUD and table partition
and orm just support mysql now.

todo list:
[cache] cache support
[orm] auto cache data for db select.
[model] distribute strategy then the orm will support distributed database and cache
[orm] relate table search
[router] router and dispatcher
[view] support bootstrap and jquery

orm usage:
write you model definition in one go file and use prebuilder to generate a implements file.

cd prebuilder
go build
./prebuilder -m path/to/you/model.go
then u can write you logic

if u want partition your table, see example model.go in test directory

CRUD
see example in test directory

