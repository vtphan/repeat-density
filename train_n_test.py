# Author: Vinhthuy Phan, 2014
import tsv
import sys
import argparse
import random
import math
import os

# slope, intercept, r_value, p_value, std_err = stats.linregress(x,y)
# y = slope * x + intercept
from scipy import stats

IGNORE = []

BIAS_TRAIN_DATA = []

def check_data_integrity(data1, data2):
   x = [ r['ID'] for r in data1 ]
   y = [ r['ID'] for r in data2 ]

   if x != y:
      print("Order of genomes in both files are not the same.")
      print(x)
      print(y)
      sys.exit()


def split_data(datax, catx, datay, caty, k):
   x = [float(r[catx]) for r in datax if r['ID'] not in IGNORE]
   y = [float(r[caty]) for r in datay if r['ID'] not in IGNORE]
   assert(len(x) == len(y))
   train_x, test_x, train_y, test_y, = [], [], [], []
   train_idx = random.sample(xrange(len(x)), k)
   for i in range(len(x)):
      if i in train_idx:
         train_x.append(x[i])
         train_y.append(y[i])
      else:
         test_x.append(x[i])
         test_y.append(y[i])
   return train_x, test_x, train_y, test_y


def biased_split(datax, catx, datay, caty):
   train_x = [float(r[catx]) for r in datax if r['ID'] in BIAS_TRAIN_DATA]
   test_x = [float(r[catx]) for r in datax if r['ID'] not in BIAS_TRAIN_DATA]
   train_y = [float(r[caty]) for r in datay if r['ID'] in BIAS_TRAIN_DATA]
   test_y = [float(r[caty]) for r in datay if r['ID'] not in BIAS_TRAIN_DATA]
   return train_x, test_x, train_y, test_y


def error(x,y,dim=1):
   if dim==1: # mean absolute error
      return sum(abs(x[i]-y[i]) for i in range(len(x)))/float(len(x)) if x else 0
   else:  # mean square error
      return math.sqrt(sum((x[i]-y[i])**2 for i in range(len(x)))/len(x)) if x else 0

def test_prediction(slope, intercept, x, y):
   prediction = [ slope*i + intercept for i in x ]
   return error(prediction, y, 1)


def train_and_test(complexity_data, performance_data, x, y, training_size, rounds):
   total_r, total_err = 0, 0
   for i in range(rounds):
      if BIAS_TRAIN_DATA:
         train_comp, test_comp, train_perf, test_perf = \
            biased_split(complexity_data, x, performance_data, y)
      else:
         train_comp, test_comp, train_perf, test_perf = \
            split_data(complexity_data, x, performance_data, y, training_size)

      # train on training data set
      slope, intercept, r_value, p_value, std_err = stats.linregress(train_comp, train_perf)
      total_r += r_value

      # use linear model to predict on testing data set
      perf_err = test_prediction(slope, intercept, test_comp, test_perf)
      total_err += perf_err
      # print ("%.4f\t%.4f" % (r_value, perf_err))

   return total_r/rounds, total_err/rounds



def run(args, complexity_keys, perf_key, training_size, ITER):
   print ("Performance:\t%s" % perf_key)
   print("\t\t%s" % '\t'.join(complexity_keys))

   for aligner in os.listdir(args['dir']):
      performance_data = tsv.Read(os.path.join(args['dir'], aligner), '\t')
      check_data_integrity(complexity_data, performance_data)
      R, err = [], []
      for x in complexity_keys:
         y = perf_key
         average_R, average_err = \
            train_and_test(complexity_data, performance_data, x, y, training_size, ITER)
         R.append(average_R)
         err.append(average_err)
      print("%s" % aligner)
      print("mean_R  \t%s" % '\t'.join([str(round(i,2)) for i in R]))
      print("mean_err\t%s" % '\t'.join([str(round(i,4)) for i in err]))



if __name__ == '__main__':
   parser = argparse.ArgumentParser(description='Train and predict short-read alignment performance using different complexity measures.')
   parser.add_argument('complexity', help='file containing complexity values of genomes')
   parser.add_argument('training_portion', type=float, help='fraction of data used for training')
   parser.add_argument('dir', help='directory containing text files storing aligner performance')
   parser.add_argument('performance_keys', nargs='+', help='Prec-100, Rec-100, Prec-75, Rec-75, Prec-50, Rec-50, ...')
   args = vars(parser.parse_args())

   complexity_data = tsv.Read(args['complexity'], '\t')
   TRAIN_FRAC = args['training_portion']
   training_size = int((len(complexity_data) - len(IGNORE)) * TRAIN_FRAC) if not BIAS_TRAIN_DATA else len(BIAS_TRAIN_DATA)
   ITER = 100 if (not BIAS_TRAIN_DATA) and (training_size < len(complexity_data)) else 1

   print ("Sample size\t%d\nTraining size\t%d (%.2f * (%d-%d))\nIteration\t%d" %
      (len(complexity_data), training_size, TRAIN_FRAC, len(complexity_data),len(IGNORE), ITER))


   complexity_keys = ['D12', 'D25', 'D50', 'D100', 'D200', 'D400', 'R12', 'R25', 'R50', 'R100', 'R200','R400', 'D', 'I' ]

   for perf_key in args['performance_keys']:
      run(args, complexity_keys, perf_key, training_size, ITER)
      print("\n")



