/*
Author: Vinhthuy Phan, Shanshan Gao
Copyright 2014
Measures of complexity: I, Ik, D, Dk, Rk
*/
package genomecomplexity

import (
    // "fmt"
    "os"
    "bufio"
    "bytes"
    "math"
    "io/ioutil"
)

type Index struct{
    data []byte
    sa []int
    lcp []int
}

func (idx *Index) Build(filename string) {
   idx.data = ReadSequence(filename)
   idx.lcp = make([]int, len(idx.data)-1)  // lcp[i] stores length of lcp of sa[i] and sa[i+1]
   idx.sa = qsufsort(idx.data)
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

func ReadSequence(file string) []byte{
   f, err := os.Open(file)
   if err != nil {
      panic(err)
   }
   defer f.Close()
   byte_array := make([]byte, 0)

   if file[len(file)-6:] == ".fasta" {
      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
         line := scanner.Bytes()
         if len(line)>0 && line[0] != '>' {
            byte_array = append(byte_array, bytes.Trim(line,"\n\r ")...)
         }
      }
   } else {
      byte_array, err = ioutil.ReadFile(file)
      if err != nil {
         panic(err)
      }
   }
   return byte_array
}