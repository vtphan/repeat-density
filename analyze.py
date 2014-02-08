import tsv
import sys
import argparse
import random
import math

# slope, intercept, r_value, p_value, std_err = stats.linregress(x,y)
# y = slope * x + intercept
from scipy import stats


if __name__ == '__main__':
   parser = argparse.ArgumentParser(description='Analyze genome complexity and alignment performance.')
   parser.add_argument('complexity', help='file containing complexity values of genomes')
   parser.add_argument('performance', help='file containing performance values of aligner')
   args = vars(parser.parse_args())
   complexities = tsv.Read(args['complexity'], '\t')
   performances = tsv.Read(args['performance'], '\t')
   if [ r['ID'] for r in complexities ] != [ r['ID'] for r in performances ]:
      print("Order of genomes in both files are not the same.")
      sys.exit()

   complexity_type = [ k for k in complexities.keys() if k!='ID' ]
   performance_type = [ k for k in performances.keys() if k!='ID' ]

   for c in complexity_type:
      for p in performance_type:
         c_val = [ float(r[c]) for r in complexities ]
         p_val = [ float(r[p]) for r in performances ]
         slope, intercept, r_value, p_value, std_err = stats.linregress(c_val, p_val)
         print "%s %s\t%f" % (c, p, r_value)
