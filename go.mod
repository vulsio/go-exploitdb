module github.com/mozqnet/go-exploitdb

go 1.12

replace (
	github.com/genuinetools/reg => github.com/tomoyamachi/reg v0.16.1-0.20190706172545-2a2250fd7c00
	gopkg.in/mattn/go-colorable.v0 => github.com/mattn/go-colorable v0.1.0
	gopkg.in/mattn/go-isatty.v0 => github.com/mattn/go-isatty v0.0.6
)

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/cheggaaa/pb v2.0.7+incompatible
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/gocarina/gocsv v0.0.0-20190821091544-020a928c6f4e
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/jinzhu/gorm v1.9.10
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.8.1
	github.com/russross/blackfriday v1.5.2
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	gopkg.in/VividCortex/ewma.v1 v1.1.1 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/cheggaaa/pb.v2 v2.0.7 // indirect
	gopkg.in/fatih/color.v1 v1.7.0 // indirect
	gopkg.in/mattn/go-colorable.v0 v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/mattn/go-isatty.v0 v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/mattn/go-runewidth.v0 v0.0.4 // indirect
)
