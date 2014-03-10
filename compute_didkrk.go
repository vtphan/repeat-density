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
    fmt.Printf("Genome\tI\tD\tD_100\tD_200\tD_400\tR_100\tR_200\tR_400\n")
    fmt.Printf("%s\t%f\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n",
        os.Args[1], idx.I(), idx.D(),
        idx.Dk(100), idx.Dk(200), idx.Dk(400), idx.Rk(100), idx.Rk(200), idx.Rk(400))
}
