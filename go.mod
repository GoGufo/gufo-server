module github.com/gogufo/gufo-server

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gogufo/gufodao v1.1.0 // indirect
	github.com/gomodule/redigo v1.8.4
	github.com/jinzhu/gorm v1.9.16
	github.com/microcosm-cc/bluemonday v1.0.4
	github.com/spf13/viper v1.7.1
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

replace github.com/johnfercher/maroto => github.com/sucsessyan/maroto v0.1.0

replace github.com/johnfercher/maroto/pkg/color => github.com/sucsessyan/maroto/pkg/color v0.1.0
