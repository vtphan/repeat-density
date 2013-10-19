package main

import (
    "fmt"
    "sort"
    "os"
    "bufio"
    "runtime"
    "bytes"
    // "flag"
)

type index Index

func (x *index) Len() int           { return len(x.sa) }
func (x *index) Less(i, j int) bool { return bytes.Compare(x.at(i), x.at(j)) < 0 }
func (x *index) Swap(i, j int)      { x.sa[i], x.sa[j] = x.sa[j], x.sa[i] }
func (a *index) at(i int) []byte    { return a.data[a.sa[i]:] }

type Index struct{
    data [] byte
    sa []int    //suffix array or LCP array for data
}

func suffixarray(data []byte) *Index{
    sa := make([]int, len(data))
    for i := range sa {
        sa[i] = i
    }
    x := &Index{data, sa}
    sort.Sort((*index)(x))
    return x
}

func (I *Index) Lcparray() *Index{
	la := make([]int, len(I.data) - 1)
	for i := 1; i < len(I.data); i++ {
		la[i-1] = len(I.LCP(i))
    }
    y := &Index{I.data, la}
    return y
}

func (I *Index) DistinctSub() float64{
    var numDistSub uint64
    var result float64
    numDistSub = uint64(I.SuffixLen(0))
    for i := 1; i < len(I.data); i++ {
        numDistSub = numDistSub + uint64(I.SuffixLen(i)) - uint64(len(I.LCP(i)))
    }
    result = float64(numDistSub)/float64((len(I.data) * (len(I.data) + 1))/2.0)
    fmt.Printf("%d\t%f\n", numDistSub, result)
    return result
}

func (I *Index) DistinctSubWithLen(length int) float64{
    var numDistSubWithLen uint64
    var result float64
    for i := 1; i < len(I.data); i++ {
        if len(I.LCP(i)) < length && I.SuffixLen(i) >= length {
            numDistSubWithLen++
        }
    }
    result = float64(numDistSubWithLen)/float64(len(I.data) - length + 1)
    fmt.Printf("%d\t%d\t%f\n", length, numDistSubWithLen, result)
    return result
}

func (I *Index) Krepeats(length int) float64{
    var numRepeatsWithLen uint64
    var result float64
    i := 0
    for i < len(I.data)-1 {
        if I.sa[i] >= length {
            numRepeatsWithLen = numRepeatsWithLen + uint64(I.Block(i, length) - i + 2)
            i = I.Block(i, length) + 1
        } else {
            i++
        }
    }
    result = float64(numRepeatsWithLen)/float64(len(I.data) - length + 1)
    fmt.Printf("%d\t%d\t%f\n", length, numRepeatsWithLen, result)
    return result
}

func (I *Index) SuffixLen(m int) int{
        stringLength := len(I.data)
        l := stringLength - I.sa[m]
        return l
}

func (I *Index) LCP(m int) []byte{
        j := I.sa[m-1]
        for i := I.sa[m]; i < len(I.data); i++ {
                if (j < len(I.data)) && (I.data[i] == I.data[j]) {
                        j++
                } else {
                        break
                }
        }
        LCP := I.data[I.sa[m-1] : j]
        return LCP
}

func (I *Index) Block(m int, length int) int{
	for i := m; i < len(I.data)-1; i++ {
        if I.sa[i] < length {
            return (i - 1)
        }
    }
    return len(I.data) - 2
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

    line , err := br.ReadString('\n')

    fmt.Println(line)

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
    // var sequence_file = flag.String("-s", "", "sequence file")
    // flag.Parse()
    stri := fastaRead(os.Args[1])

    runtime.GOMAXPROCS(4)
    suffixarray := suffixarray(stri)

    suffixarray.DistinctSub()

    suffixarray.DistinctSubWithLen(100)
    suffixarray.DistinctSubWithLen(200)
    suffixarray.DistinctSubWithLen(400)

    lcparray := suffixarray.Lcparray()
    lcparray.Krepeats(100)
    lcparray.Krepeats(200)
    lcparray.Krepeats(400)
}