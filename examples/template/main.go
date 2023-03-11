package main

import (
	"github.com/go-msvc/config"
	_ "github.com/go-msvc/utils/api"
	"github.com/go-msvc/utils/ms"
)

type template struct{}

func main() {
	config.AddSource("config.json", config.File("./config.json"))
	ms := ms.New(
		ms.WithOper("greet", operGreet),
		ms.WithOper("wave", operWave),
		ms.WithOper("add_group", addGroup),
		ms.WithOper("get_group", getGroup),
		ms.WithOper("upd_group", updGroup),
		ms.WithOper("del_group", delGroup),
		ms.WithOper("find_groups", findGroup),
	)
	ms.Configure()
	ms.Serve()
}
