## Compute Various Measures of Complexity

genome-complexity.go:
	Compute various complexities of genomes, include repeat density,
	distinct substring density, and I complexity.

Install genome-complexity.go:
   go get github.com/vtphan/sequence-complexity

~~~ go
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
~~~


## Cross Validation

    usage: cross_validation.py [-h]
                               complexity training_portion dir performance_keys
                               [performance_keys ...]

    Train and predict short-read alignment performance using different complexity
    measures.

    positional arguments:
      complexity        file containing complexity values of genomes
      training_portion  fraction of data used for training
      dir               directory containing text files storing aligner
                        performance
      performance_keys  Prec-100, Rec-100, Prec-75, Rec-75, Prec-50, Rec-50, ...

    optional arguments:
      -h, --help        show this help message and exit


Example: using 40% of data for training, 60% for testing

    python cross_validation.py complexity.txt 100 Prec-100 Rec-100
