package main

import (
	"flag"
	"fmt"
	"strconv"
)

type intslice []int

func (i *intslice) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *intslice) Set(value string) error {
	fmt.Printf("%s\n", value)
	tmp, err := strconv.Atoi(value)
	if err != nil {
		*i = append(*i, -1)
	} else {
		*i = append(*i, tmp)
	}
	return nil
}

func main() {
	var enterprises intslice

	flag.Var(&enterprises, "e", "Enterprise ID. This flag is repeatable.")
	flag.Var()
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	} else {
		fmt.Println("Here are the values in 'enterprises'")
		for i := 0; i < len(enterprises); i++ {
			fmt.Printf("%d\n", enterprises[i])
		}
	}
}
