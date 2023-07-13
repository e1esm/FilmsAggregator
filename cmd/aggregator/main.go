package main

import (
	"fmt"
	"github.com/e1esm/FilmsAggregator/utils/config"
)

func main() {
	cfg := config.NewConfig()
	fmt.Println(cfg.Aggregator.Name)
}
