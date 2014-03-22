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
    fmt.Println(idx.D())
    fmt.Println(idx.I())
    for i:=12; i<150; i+=5 {
        fmt.Println(i,"\t",idx.Ik(i))
    }
}
