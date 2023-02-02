import sys
from generateBoardUtility import *
from globals import *
import time




states = {
  ""
}
while 1:
  inp = sys.stdin.readline().strip("\n").split(" ")
  case = inp.pop(0)
  
  if case == "gb":

    batchName = inp[0]
    gridSize = inp[1]
    batchSize = inp[2]
    cropAmt = inp[3]
    generateBoard(batchName, gridSize, batchSize, cropAmt)
    
  # elif "cb":
  #   cropBatch(cropAmt, batches)
  # elif "ab":
  #   addBatch(cropAmt, batches)
  # elif "db":
  #   deleteBatch(batches)
  # elif "sb":
  #   saveBatches(batches)
  elif case == "pb":
    batchName = inp[0]
    printBatch(batchName)

  elif case == "q":
    break

