module github.com/gogufo/gufo-server

go 1.16

require (
	github.com/gogufo/gufodao v1.2.1
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/urfave/cli/v2 v2.2.0
)

replace github.com/johnfercher/maroto => github.com/sucsessyan/maroto v0.1.0

replace github.com/johnfercher/maroto/pkg/color => github.com/sucsessyan/maroto/pkg/color v0.1.0
