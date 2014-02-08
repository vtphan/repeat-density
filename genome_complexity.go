/*
Author: Vinhthuy Phan, Shanshan Gao
Copyright 2014
*/
package main

import (
    "fmt"
    "sort"
    "os"
    "bufio"
    "bytes"
    "math"
    // "flag"
)

type Index struct{
    data [] byte
    sa []int
    lcp []int
}

func (x *Index) Len() int           { return len(x.sa) }
func (x *Index) Less(i, j int) bool { return bytes.Compare(x.at(i), x.at(j)) < 0 }
func (x *Index) Swap(i, j int)      { x.sa[i], x.sa[j] = x.sa[j], x.sa[i] }
func (a *Index) at(i int) []byte    { return a.data[a.sa[i]:] }


func (idx *Index) build(filename string) {
    idx.data = fastaRead(filename)
    idx.sa = make([]int, len(idx.data))
    idx.lcp = make([]int, len(idx.data)-1)
    for i := 0; i < len(idx.data); i++ {
        idx.sa[i] = i
    }
    sort.Sort(idx)
    for i := 1; i < len(idx.data); i++ {
        idx.lcp[i-1] = idx.lcp_len(i)
    }
}

// length of longest common prefix of data[SA[m]:] and data[SA[m-1]:]
func (idx *Index) lcp_len(m int) int{
    L, i, j := len(idx.data), idx.sa[m], idx.sa[m-1]
    for i<L && j<L && idx.data[i]==idx.data[j] {
        i++
        j++
    }
    return j - idx.sa[m-1]
}

func (idx Index) suffix_len(m int) int{
    return len(idx.data) - idx.sa[m]
}

// D = rate of distinct substrings
func (idx Index) D() float64{
    numDistSub := uint64(idx.suffix_len(0))
    for i := 1; i < len(idx.data); i++ {
        numDistSub += uint64(idx.suffix_len(i) - idx.lcp[i-1])
    }
    return 2.0 * (float64(numDistSub)/float64(len(idx.data)))/ float64(len(idx.data) + 1)
}

// Dk = rate of distinct k-mers
func (idx Index) Dk(k int) float64{
    var c uint64 = 0
    for i := 1; i < len(idx.data); i++ {
        if idx.lcp[i-1] < k && idx.suffix_len(i) >= k {
            c++
        }
    }
    return float64(c)/float64(len(idx.data) - k + 1)
}

func (idx Index) Block(m int, k int) int{
    for i := m; i < len(idx.data)-1; i++ {
        if idx.lcp[i] < k {
            return (i - 1)
        }
    }
    return len(idx.data) - 2
}

// Rk = k-repeat density
func (idx Index) Rk(k int) float64{
    var c uint64 = 0
    i := 0
    for i < len(idx.data)-1 {
        // fmt.Println(i, idx.lcp[i], idx.Block(i,k), c)
        if idx.lcp[i] >= k {
            c += uint64(idx.Block(i, k) - i + 2)
            i = idx.Block(i, k) + 1
        } else {
            i++
        }
    }
    return float64(c)/float64(len(idx.data) - k + 1)
}

// I complexity (Becher & Heiber, 2012)
func (idx *Index) I() float64 {
    var sum float64 = 0
    for _, v := range idx.lcp {
        sum += (math.Log(float64(v+2)) - math.Log(float64(v+1))) / math.Log(4.0)
    }
    return sum
}

func fastaRead(sequence_file string) []byte {
    f,err := os.Open(sequence_file)
    if err != nil{
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }

    defer f.Close()
    br := bufio.NewReader(f)
    byte_array := bytes.Buffer{}

    _ , err = br.ReadString('\n')
    for {
        line , isPrefix, err := br.ReadLine()
        if err != nil || isPrefix{
            break
        } else {
            byte_array.Write([]byte(line))
        }
    }
    byte_array.Write([]byte("$"))
    input := []byte(byte_array.String())
    return input
}

func main(){
    if len(os.Args) != 2 {
        panic("must provide sequence file.")
    }
    idx := new(Index)
    idx.build(os.Args[1])
    fmt.Printf("I\tD\tD_100\tD_200\tD_400\tR_100\tR_200\tR_400\n")
    fmt.Printf("%f\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n", idx.I(), idx.D(),
        idx.Dk(100), idx.Dk(200), idx.Dk(400), idx.Rk(100), idx.Rk(200), idx.Rk(400))
}
