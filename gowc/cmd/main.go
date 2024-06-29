package main

import (
	"fmt"
	"gowc/pkg/wc"
	"gowc/pkg/wcflag"
)

func main() {

	fmt.Println("Running gowc ...")

	engine := wc.InitEngine(wcflag.ParseFlag())
	engine.Count()

}
