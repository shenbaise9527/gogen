# gogen
a cli tool to generate golang code, mainly use [go-zero](https://github.com/tal-tech/go-zero) and [gorm](https://gorm.io/)

## model
generate model code, only support mysql
``` bash
$ ./gogen model datasource -h
NAME:
   gogen model datasource - generate model from datasource

USAGE:
   gogen model datasource [command options] [arguments...]

OPTIONS:
   --url value      data soucre of database, mysql: "root:password@tcp(127.0.0.1:3306)/database"
   --table value    the tables in the database,support for comma separation
   --dir value      the target dir
   --cache value    generate code with cache [optional]
   --tracing value  generate code with tracing [optional]
   --help, -h       show help (default: false)
```