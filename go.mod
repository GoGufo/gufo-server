module github.com/gogufo/gufo-server

go 1.16

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/gogufo/gufodao v1.2.3
	github.com/gomodule/redigo v1.8.5 // indirect
	github.com/lib/pq v1.10.3 // indirect
	github.com/microcosm-cc/bluemonday v1.0.15
	github.com/spf13/viper v1.9.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/net v0.0.0-20211014222326-fd004c51d1d6 // indirect
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/driver/postgres v1.1.2 // indirect
	gorm.io/gorm v1.21.16 // indirect
)

replace github.com/johnfercher/maroto => github.com/sucsessyan/maroto v0.1.0

replace github.com/johnfercher/maroto/pkg/color => github.com/sucsessyan/maroto/pkg/color v0.1.0
