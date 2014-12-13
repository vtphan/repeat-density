complexity.go computes various measures of genome complexity, including
the I complexity, D complexity, D_k complexity and R_k complexity.  The program
is written in the Go programming language.

complexity.go will remove Ns from the sequence (in FASAT format). Long stretches of N's
must be removed since they would affect incorrectly the complexity of the sequence.

## Usage

genome-complexity.go:
	Compute various complexities of genomes, include repeat density,
	distinct substring density, and I complexity.

Install genome-complexity.go:

```
    go get github.com/vtphan/sequence-complexity
```

Compute D, I, D_25, D_50, D_75, D_100, D_125, D_150, D_175, D_200:

```
    go run usage/compute_complexity.go usage/CP003835.fasta
```

If you want to compute other D_k, modify usage/compute_complexity.go accordingly.