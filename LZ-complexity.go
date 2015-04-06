/*
 * =============================================================================
 *
 *       Filename:  LZ-complexity.go
 *
 *    Description:  Compute the Lempel-Ziv complexity of binary strings. 
                  Take a file as input (fasta/fa or one string/line). 
 *        Author:  Quang Tran @ University of Memphis  
 *
 * =============================================================================
 */

package genomecomplexity

import (
    //"fmt"
    "os"
    "bufio"
    "bytes"
    "math"
    "io/ioutil"
)
/*
func main(){
   if len(os.Args) != 2 {
      panic("must provide sequence file.")
   }
   
   seq := ReadSequence(os.Args[1])
   //fmt.Println(string(seq))
   //fmt.Println(string(reverse(seq)))
   if len(seq)>0 {      
      c78 := LZ78(seq)
      fmt.Print(c78, "\t")

      //c76 := LZ76(seq)
      //fmt.Print(c76, "\t")

      rev78 := LZ78(reverse(seq))
      //rev76 := LZ76(reverse(seq))
      fmt.Print(c78+rev78, "\t")
      //fmt.Print(c76+rev76, "\t")

      //fmt.Print(Kolmogorov(seq), "\t")

      nom := float64(len(seq))/math.Log2(float64(len(seq)))
            
      // Normalize
      fmt.Print(float64(c78)/(nom), "\t")
      //fmt.Print(float64(c76)/(nom), "\t")
      fmt.Println(float64(c78+rev78)/(nom), "\t")
      //fmt.Println(float64(c76+rev76)/(nom))
   }
}
*/
func Reverse(s []byte) []byte {
  rev := make([]byte, len(s))
  for i, elem := range s {
    if rune(elem) == 'A' {
      rev[i] = 'T'
    } else if rune(elem) == 'T' {
      rev[i] = 'A'
    } else if rune(elem) == 'C' {
      rev[i] = 'G'
    } else if rune(elem) == 'G' {
      rev[i] = 'C'
    } else {
      rev[i] = elem
    }
  }
  return rev
}

func Kolmogorov(s []byte) float64 {
      c := float64(LZ76(s)) * math.Log2(float64(len(s)))
   return c
}

/* LZ76 implemented following "Easily calculable measure for the complexity of spatiotemporal patterns" F. Kaspar and H. G. Schuster. */
func LZ76(s []byte) int {
   c := 1
   l := 1
   i := 0
   k := 1
   n := len(s)
   kmax := 1
   stop := 0
   for (stop ==0) {
      if (s[i+k-1] != s[l+k-1]) {
         if (k > kmax) {
        kmax=k;
      }
      i++;
      if (i==l) {
         c++;
         l += kmax;
         if (l+1>n) {
            stop =1;
         } else {
            i=0;
            k=1;
            kmax=1;
         }
      } else {
         k=1;
      }
    } else {
      k++;
      if (l+k > n) {
         c++;
         stop =1;
      }
    }
  }
  return c;
}

func LZ78(s []byte) int {
   dict := make(map[string]bool)
   block := ""
   for i, l:=0, 0; i<len(s); i++ {
      block += string(s[i])
      dict[block] = true
      if len(dict) > l {
         l = len(dict)
         //fmt.Printf("%s.",block)
         block = ""
      }
   }
   if block != "" {
      return len(dict)+1
   }
   return len(dict)
}
/*
func LZ78(seq []byte) int {
   m := make(map[string]bool)
   m[string(seq[0])] = true
   c := 1
   i := 1
   for (i<len(seq)) {
      temp := string(seq[i])
      _, ok := m[temp]
      if ok {
         for {
            i = i + 1
            if (i==len(seq)) {break}
            temp = temp + string(seq[i])
            _, f := m[temp]
            if !f {
               c++
               m[temp] = true
               break
            }
         }
      } else {
         c++
         m[temp] = true
      }
      i = i + 1
   }    
   return c
}
*/
func ReadSequence(file string) []byte{
   f, err := os.Open(file)
   if err != nil {
      panic(err)
   }
   defer f.Close()
   byte_array := make([]byte, 0)
   Ns := []byte("N")
   None := []byte("")
   if file[len(file)-6:] == ".fasta" || file[len(file)-3:] == ".fa" {
      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
         line := scanner.Bytes()
         if len(line)>0 && line[0] != '>' {
            byte_array = append(byte_array, bytes.Replace(bytes.Trim(line,"\n\r "), Ns, None, -1)...)
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