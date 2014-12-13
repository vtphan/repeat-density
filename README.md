This program computes various measures of genome complexity, including
the I complexity, D complexity, D_k complexity and R_k complexity.  The program
is written in the Go programming language.

## Usage

genome-complexity.go:
	Compute various complexities of genomes, include repeat density,
	distinct substring density, and I complexity.

Install genome-complexity.go:

```
    go get github.com/vtphan/sequence-complexity
```

Example:

```
    go run test/compute_complexity.go test/CP003835.fasta
```