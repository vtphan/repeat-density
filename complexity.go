/*
Author: Vinhthuy Phan, Shanshan Gao
Copyright 2014
Measures of complexity: I, Ik, D, Dk, Rk
*/
package genomecomplexity

import (
    "fmt"
    "sort"
    "os"
    "bufio"
    "bytes"
    "math"
)

type Index struct{
    data []byte
    sa []int
    lcp []int
}

func (x *Index) Len() int           { return len(x.sa) }
func (x *Index) Less(i, j int) bool { return bytes.Compare(x.at(i), x.at(j)) < 0 }
func (x *Index) Swap(i, j int)      { x.sa[i], x.sa[j] = x.sa[j], x.sa[i] }
func (a *Index) at(i int) []byte    { return a.data[a.sa[i]:] }


func (idx *Index) Build(filename string) {
    idx.data = fastaRead(filename)
    idx.sa = make([]int, len(idx.data))
    idx.lcp = make([]int, len(idx.data)-1)  // lcp[i] stores length of lcp of sa[i] and sa[i+1]
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

// D = rate of distinct substrings
func (idx Index) D() float64{
    c := uint64(len(idx.data) - idx.sa[0])
    for i := 1; i < len(idx.data); i++ {
        c += uint64(len(idx.data) - idx.sa[i] - idx.lcp[i-1])
    }
    return 2.0 * (float64(c)/float64(len(idx.data)))/ float64(len(idx.data) + 1)
}

// Dk = rate of distinct k-mers
func (idx Index) Dk(k int) float64{
    var c uint64 = 0
    if idx.sa[0] <= len(idx.data)-k {
      c++
      // fmt.Println(string(idx.data[idx.sa[0] : idx.sa[0]+k]))
    }
    for i := 1; i < len(idx.data); i++ {
        if idx.lcp[i-1] < k && idx.sa[i] <= len(idx.data)-k {
            c++
            // fmt.Println(string(idx.data[idx.sa[i] : idx.sa[i]+k]))
        }
    }
    // return float64(c)
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

func (idx Index) Ik(k int) float64{
   var sum float64 = 0
   for i := 1; i < len(idx.data); i++ {
     if idx.lcp[i-1] < k && len(idx.data)-idx.sa[i] >= k {
        sum += (math.Log(float64(idx.lcp[i-1]+2)) - math.Log(float64(idx.lcp[i-1]+1))) / math.Log(4.0)
     }
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
    // do not add $ at the end of the text
    // byte_array.Write([]byte("$"))
    input := []byte(byte_array.String())
    return input
}

