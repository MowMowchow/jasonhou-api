def buildTrie(words, trie):
    def tDfs(word, ind, t):
      if ind == len(word):
        return
      if not word[ind] in t:
        t[word[ind]] = {}
      if len(word)-1 == ind:
        t[word[ind]][1] = 1 # mark that word exists
        t[word[ind]][2] = word # mark that word exists
      tDfs(word, ind+1, t[word[ind]])

    for word in words:
      tDfs(word, 0, trie)


alphabet = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z']
alphabetLen = len(alphabet)
f = open("./scripts/fordle/wordList.txt", "r")
wordList = f.read().split("\n")
f.close()
trie = {}
buildTrie(wordList, trie)
batchList = {}