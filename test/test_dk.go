package main

import (
    "fmt"
    "os"
    comp "github.com/vtphan/sequence-complexity"
)

func main(){
    if len(os.Args) != 2 {
        panic("must provide sequence file.")
    }
    idx := new(comp.Index)
    idx.Build(os.Args[1])
    fmt.Printf("%s\t%f\t%f\t%f\t%f\n", os.Args[1], idx.Dk(1), idx.Dk(2), idx.Dk(3),idx.Dk(4))
}
