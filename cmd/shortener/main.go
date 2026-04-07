package main

import (
	"context"
	"fmt"
	"shortener/internal/codec"
	"shortener/internal/generator"
)

func main() {
	gen := generator.NewMemoryGenerator()
	ctx := context.Background()
	for range 100 {
		i, _ := gen.Next(ctx)
		fmt.Println(codec.Encode(uint64(i)))
	}

}
