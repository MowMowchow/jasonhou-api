from globals import *
import random
import heapq
import time


def findWords(board: list[list[str]]) -> list[str]:
  R, C = len(board), len(board[0]) 
  moves = [[0, 1], [1, 0], [-1, 0], [0, -1], [1, 1], [-1, -1], [-1, 1], [1, -1]]
  ans = {}

  def dfs(x, y, t, vis):
    if not (x, y) in vis:
      vis[x, y] = 1

      if 1 in t:
        ans[t[2]] = 1
      for mx, my in moves:
        nx, ny = x+mx, y+my
        if 0 <= nx < C and 0 <= ny < R:
          if board[ny][nx] in t and (nx, ny) not in vis:
            dfs(nx, ny, t[board[ny][nx]], dict(vis))
    

  for x in range(C):
    for y in range(R):
      if board[y][x] in trie:
        dfs(x, y, trie[board[y][x]], {})
  
  return len(ans.keys())


def generateBoard(batchName, gridSize, batchSize, cropAmt):
  startTime = time.time()
  gridSize = int(gridSize)
  batchSize = int(batchSize)
  tempBoardList = []
  if batchName in batchList:
    print("BATCH NAME ALREADY EXISTS")
    return
  batchList[batchName] = {"boards":[]}

  for batch in range(batchSize):
    board = ["".join([alphabet[random.getrandbits(5)%alphabetLen] for x in range(gridSize)]) for y in range(gridSize)]
    wordAmt = findWords(board)
    heapq.heappush(tempBoardList, [wordAmt, board]) 
  
  raw = False if "%" in cropAmt else True 

  if raw:
    cropAmt = int(cropAmt)
    if (batchSize < cropAmt):
      print("BAD INPUT")
      return

    saveAmt = batchSize-cropAmt

    for i in range(cropAmt):
      heapq.heappop(tempBoardList)

    for i in range(saveAmt):
      batchList[batchName]["boards"].append(heapq.heappop(tempBoardList)[1])
  
  else:
    cropAmt = int(cropAmt.replace("%", ""))
    if cropAmt > 100 or cropAmt < 0:
      print("BAD INPUT")
      return

    remAmt = int(batchSize*(.001*cropAmt))
    saveAmt = batchSize-remAmt

    for i in range(remAmt):
      heapq.heappop(tempBoardList)
    
    for i in range(saveAmt):
      batchList[batchName]["boards"].append(heapq.heappop(tempBoardList)[1])


  print("DONE GENERATING BATCHES IN: ", time.time()-startTime, "seconds")
  

def cropBatch(cropAmt, batches):
  raw = False if "%" in CropAmt else True 

  for batchName in batches:
    if not batchName in batchList:
      print("BATCH:", batch, "DNE")
      return

    boards = heapq.heapify(batchList[batchName]["boards"])
    batchList[batchName]["boards"] = {}
    batchSize = len(batchList)

    if raw:
       cropAmt = int(cropAmt)
       if (batchSize < cropAmt):
         print("BAD INPUT")
         return

       saveAmt = batchSize-cropAmt

       for i in range(cropAmt):
         heapq.heappop(boards)

       for i in range(saveAmt):
         batchList[batchName]["boards"].append(heapq.heappop(boards)[1])
  
    else:
      cropAmt = int(cropAmt.replace("%", ""))
      if cropAmt > 100 or cropAmt < 0:
        print("BAD INPUT")
        return

      boards = batchList[batchName]["boards"]
      batchList[batchName]["boards"] = {}
      batchSize = len(batchList)

      remAmt = int(batchSize*(.001*cropAmt))
      saveAmt = batchSize-remAmt

      for i in range(remAmt):
        heapq.heappop(boards)
    
      for i in range(saveAmt):
        batchList[batchName]["boards"].append(heapq.heappop(boards)[1])

  print("DONE!")


def addBatch(cropAmt, batches):
  return

def deleteBatch(batches):
  return


def saveBatches(batches):
  return


def getWordCount(board):
  return findWords(board)


def printBatch(batch):
  if not batch in batchList:
    print("COULD NOT FIND BATCH:", batch)
    return

  batchLen = len(batchList[batch]["boards"])

  maxCnt = 0 
  for i in range(batchLen):
    boardWordCount = getWordCount(batchList[batch]["boards"][i])
    maxCnt = max(maxCnt, boardWordCount)
    print("BOARD", i, "| TOTAL WORDS:", boardWordCount, ":")
    for row in batchList[batch]["boards"][i]:
      print(row)
    print()
    print()

  print()
  print("TOTAL BOARDS:", batchLen, "| MAX WORD COUNT:", maxCnt)
  print("DONE")

  return
