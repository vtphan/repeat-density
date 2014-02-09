# Author: Vinhthuy Phan, 2014
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


def train_and_test(complexity_data, performance_data, training_size, rounds):
   total_r, total_err = 0, 0
   for i in range(rounds):
      train_comp, test_comp, train_perf, test_perf = \
         split_data(complexity_data, 'D', performance_data, 'Rec-400', training_size)
      slope, intercept, r_value, p_value, std_err = stats.linregress(train_comp, train_perf)
      perf_err = test_prediction(slope, intercept, test_comp, test_perf)
      total_r, total_err = total_r + r_value, total_err + perf_err
      # print ("%.4f\t%.4f" % (r_value, perf_err))
   return total_r/rounds, total_err/rounds


if __name__ == '__main__':
   TRAIN_FRAC = 0.5
   ITER = 100

   parser = argparse.ArgumentParser(description='Analyze genome complexity and alignment performance.')
   parser.add_argument('complexity', help='file containing complexity values of genomes')
   parser.add_argument('performances', nargs='+', help='file(s) containing performance values of aligner')
   args = vars(parser.parse_args())
   complexity_data = tsv.Read(args['complexity'], '\t')
   training_size = int(len(complexity_data) * TRAIN_FRAC)
   print ("Sample size\t%d\nTraining size\t%d\nIteration\t%d\nData\tMean_R\tMean_Error" %
      (len(complexity_data), TRAIN_FRAC, ITER))

   for performance in args['performances']:
      performance_data = tsv.Read(performance, '\t')
      check_data_integrity(complexity_data, performance_data)
      average_R, average_err = \
         train_and_test(complexity_data, performance_data, training_size, ITER)
      print("%s\t%.4f\t%.4f" % (performance, average_R, average_err))


   # complexity_type = [ k for k in complexities.keys() if k!='ID' ]
   # performance_type = [ k for k in performances.keys() if k!='ID' ]
   # for c in complexity_type:
   #    for p in performance_type:
   #       c_val = [ float(r[c]) for r in complexities ]
   #       p_val = [ float(r[p]) for r in performances ]
   #       slope, intercept, r_value, p_value, std_err = stats.linregress(c_val, p_val)
   #       print "%s %s\t%f\t%f\t%f" % (c, p, r_value, p_value, std_err)

