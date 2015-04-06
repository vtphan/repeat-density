/*
Compute D, I, D_25, D_50, D_75, D_100, D_125, D_150, D_175, D_200, LZ78 of an input sequence
in FASTA format.
*/
package main

import (
    "fmt"
    "os"
    "math"
    comp "../../sequence-complexity"
)

func main(){
    if len(os.Args) != 2 {
        panic("must provide sequence file.")
    }
    seq := comp.ReadSequence(os.Args[1])
    idx := new(comp.Index)
    idx.Build(os.Args[1])
    fmt.Print(idx.D(), "\t")
    fmt.Print(idx.I(), "\t")
    for k:=25; k<=200; k+=25 {
        fmt.Print(idx.Dk(k), "\t")
    }
    // LZ-complexity
    if len(seq)>0 {      
      c78 := comp.LZ78(seq)
      fmt.Print(c78, "\t")

      rev78 := comp.LZ78(comp.Reverse(seq))
      fmt.Print(c78+rev78, "\t")

      nom := float64(len(seq))/math.Log2(float64(len(seq)))
            
      // Normalize
      fmt.Print(float64(c78)/(nom), "\t")
      fmt.Println(float64(c78+rev78)/(nom), "\t")
   }
}