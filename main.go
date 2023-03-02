package main

import (
	cmd "github.com/anousoneFS/clean-architecture/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
