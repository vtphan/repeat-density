/*
Usage: go  run  repeat-density.go  genome-file

This program computes the density of repeats in a genome.  It compute P(k | S)
	for various values of k.

Given a k-mer k and a genome S, P(k | S) = count(k) / |S|, where
	count(k) = frequency of k in S and reverse complement of S.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

func analyze_kmer(S string, T string, k int, result chan string) {
	repeats := make(map[string]int)
	total := 0
	for i := 0; i < len(S)-k+1; i++ {
		kmer := S[i : i+k]
		repeats[kmer] += 1
		if repeats[kmer] > 1 {
			if repeats[kmer] == 2 {
				total++
			}
			total++
		}

		kmer_rc := T[i : i+k]
		repeats[kmer_rc] += 1
		if repeats[kmer_rc] > 1 {
			if repeats[kmer_rc] == 2 {
				total++
			}
			total++
		}
	}
	result <- fmt.Sprintf("%d\t%d\t%f", k, total, float64(total)/float64(len(S)))
}

func complement(r rune) rune {
	switch r {
	case 'A', 'a':
		return 'T'
	case 'T', 't':
		return 'A'
	case 'G', 'g':
		return 'C'
	case 'C', 'c':
		return 'G'
	default:
		return r
	}
}

func reverse_complement(S string) string {
	runes := []rune(S)
	for i, j := 0, len(runes)-1; i <= j; i, j = i+1, j-1 {
		runes[i], runes[j] = complement(runes[j]), complement(runes[i])
	}
	return string(runes)
}

func main() {
	if len(os.Args) != 2 {
		panic("must provide sequence file.")
	}

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	Seq := string(b)
	Seq_rc := reverse_complement(Seq)
	lengths := []int{35, 51, 76, 100, 200, 400, 600, 800}
	result := make(chan string)

	runtime.GOMAXPROCS(4)

	for _, v := range lengths {
		go analyze_kmer(Seq, Seq_rc, v, result)
	}
	fmt.Println(os.Args[1])
	for i := 0; i < len(lengths); i++ {
		fmt.Printf("%s\n", <-result)
	}
}
