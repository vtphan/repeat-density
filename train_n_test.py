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


def train_and_test(complexity_data, performance_data, x, y, training_size, rounds):
   total_r, total_err = 0, 0
   for i in range(rounds):
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


if __name__ == '__main__':
   ITER = 100
   comparisons = dict (
      D = ['Prec-100','Rec-100','Prec-200','Rec-200','Prec-400','Rec-400'],
      D100 = ['Prec-100', 'Rec-100'],
      D200 = ['Prec-200', 'Rec-200'],
      D400 = ['Prec-400', 'Rec-400'],
      R100 = ['Prec-100', 'Rec-100'],
      R200 = ['Prec-200', 'Rec-200'],
      R400 = ['Prec-400', 'Rec-400'],
   )
   parser = argparse.ArgumentParser(description='Train and predict short-read alignment performance using different complexity measures.')
   parser.add_argument('complexity', help='file containing complexity values of genomes')
   parser.add_argument('aligners', nargs='+', help='file(s) containing performance values of aligner')
   parser.add_argument('TRAIN_FRAC', default=0.5, nargs='?', type=float, help='fraction of data used for training')

   # parser.add_argument('x', help='I, D, D100, D200, D400, R100, R200, R400')
   # parser.add_argument('y', help='Prec-100, Prec-200, Prec-400, Rec-100, Rec-200, Rec-400')
   args = vars(parser.parse_args())

   complexity_data = tsv.Read(args['complexity'], '\t')
   TRAIN_FRAC = args['TRAIN_FRAC']
   training_size = int(len(complexity_data) * TRAIN_FRAC)
   # x = args['x']
   # y = args['y']

   print ("Sample size\t%d\nTraining size\t%d (%.2f * %d)\nIteration\t%d\nData\tMean_R\tMean_Error" % (len(complexity_data), training_size, TRAIN_FRAC, len(complexity_data), ITER))

   for aligner in args['aligners']:
      performance_data = tsv.Read(aligner, '\t')
      check_data_integrity(complexity_data, performance_data)
      R, err, label = [], [], []
      for x, ys in comparisons.items():
         for y in ys:
            average_R, average_err = \
               train_and_test(complexity_data, performance_data, x, y, training_size, ITER)
            R.append(average_R)
            err.append(average_err)
            label.append('%s_%s' % (x[:2], y.split('-')[0][0] + y.split('-')[1][0]))
      print("%s average R (row 1) and average error (row 2)\n%s" % (aligner, '\t'.join(label)))
      print("%s" % '\t'.join([str(round(i,4)) for i in R]))
      print("%s" % '\t'.join([str(round(i,4)) for i in err]))

   # complexity_type = [ k for k in complexities.keys() if k!='ID' ]
   # performance_type = [ k for k in performances.keys() if k!='ID' ]
   # for c in complexity_type:
   #    for p in performance_type:
   #       c_val = [ float(r[c]) for r in complexities ]
   #       p_val = [ float(r[p]) for r in performances ]
   #       slope, intercept, r_value, p_value, std_err = stats.linregress(c_val, p_val)
   #       print "%s %s\t%f\t%f\t%f" % (c, p, r_value, p_value, std_err)

