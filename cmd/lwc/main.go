package main

import (
	"github.com/timdp/lwc/internal/app/lwc"
)

var version = "master"
var date = ""

func main() {
	lwc.Run(version, date)
}
