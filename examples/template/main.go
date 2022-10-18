package main

import (
	"github.com/go-msvc/utils/ms"
)

type template struct{}

func main() {
	ms := ms.New("template",
		ms.WithConfig("template", template{}, Config{}),
		ms.WithOper("greet", operGreet),
		ms.WithOper("wave", operWave),
	)
	ms.Configure()
	ms.Serve()
}
