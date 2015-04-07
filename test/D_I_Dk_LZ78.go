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
    
    // D, I, D_25, D_50, D_75, D_100, D_125, D_150, D_175, D_200
    var idx comp.Index
    idx.Build(os.Args[1])
    fmt.Print(idx.D(), "\t")
    fmt.Print(idx.I(), "\t")
    for k:=25; k<=200; k+=25 {
        fmt.Print(idx.Dk(k), "\t")
    }

    for k:=25; k<=200; k+=25 {
        fmt.Print(idx.Rk(k), "\t")
    }

    // LZ-complexity
    seq := comp.ReadSequence(os.Args[1])
    if len(seq)>0 {      
      c78 := comp.LZ78(seq)
      fmt.Print(c78, "\t")

      norm := NormLZ78(len(seq))
      // Normalize
      fmt.Print(float64(c78)/(norm), "\t")
   }
}

func NormLZ78(n int) float64 {
  var rs, i, c float64
  rs = 0
  i = 0
  c = 0
  for {
    i++
    f := math.Pow(4,i)
    c += i * f
    if (c>float64(n)) { break } else { rs += f }
  }  
  rs += (float64(n)-rs)/i
  return rs
}