import tsv
import sys
import argparse
import random
import math

# slope, intercept, r_value, p_value, std_err = stats.linregress(x,y)
# y = slope * x + intercept
from scipy import stats

def check_data_integrity(data1, data2):
   x = [ r['ID'] for r in data1 ]
   y = [ r['ID'] for r in data2 ]
   if x != y:
      print("Order of genomes in both files are not the same.")
      print(x)
      print(y)
      sys.exit()


def split_data(datax, catx, datay, caty, k):
   x = [float(r[catx]) for r in datax]
   y = [float(r[caty]) for r in datay]
   assert(len(x) == len(y))
   train_x, test_x, train_y, test_y, = [], [], [], []
   train_idx = random.sample(xrange(len(x)), k)
   # print train_idx
   for i in range(len(x)):
      if i in train_idx:
         train_x.append(x[i])
         train_y.append(y[i])
      else:
         test_x.append(x[i])
         test_y.append(y[i])
   return train_x, test_x, train_y, test_y


def rmsd(x,y):
   return math.sqrt(sum((x[i]-y[i])**2 for i in range(len(x)))/len(x)) if x else 0

def test_prediction(slope, intercept, x, y):
   prediction = [ slope*i + intercept for i in x ]
   return rmsd(prediction, y)

if __name__ == '__main__':
   parser = argparse.ArgumentParser(description='Analyze genome complexity and alignment performance.')
   parser.add_argument('complexity', help='file containing complexity values of genomes')
   parser.add_argument('performance', help='file containing performance values of aligner')
   args = vars(parser.parse_args())
   complexities = tsv.Read(args['complexity'], '\t')
   performances = tsv.Read(args['performance'], '\t')

   check_data_integrity(complexities, performances)

   k = 15
   M = 200
   total_r, total_err = 0, 0
   for i in range(M):
      train_comp, test_comp, train_perf, test_perf = \
         split_data(complexities, 'D', performances, 'Rec-400', k)
      slope, intercept, r_value, p_value, std_err = stats.linregress(train_comp, train_perf)
      perf_err = test_prediction(slope, intercept, test_comp, test_perf)
      total_r, total_err = total_r + r_value, total_err + perf_err
      # print ("%.4f\t%.4f" % (r_value, perf_err))

   print ("Average r, error:\t%.4f, %.4f" % (total_r/M, total_err/M))




   # complexity_type = [ k for k in complexities.keys() if k!='ID' ]
   # performance_type = [ k for k in performances.keys() if k!='ID' ]
   # for c in complexity_type:
   #    for p in performance_type:
   #       c_val = [ float(r[c]) for r in complexities ]
   #       p_val = [ float(r[p]) for r in performances ]
   #       slope, intercept, r_value, p_value, std_err = stats.linregress(c_val, p_val)
   #       print "%s %s\t%f\t%f\t%f" % (c, p, r_value, p_value, std_err)

