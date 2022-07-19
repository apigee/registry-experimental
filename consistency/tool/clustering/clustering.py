import numpy as np
import pandas as pd
from sklearn.cluster import dbscan
from strsimpy.sorensen_dice import SorensenDice
from sklearn.cluster import dbscan
from collections import Counter
import warnings
class ClusterWords:
    def __init__(self, stub, words):
        self.stub = stub
        self.words = words

    def clean_words(self):
        words = self.words
        valid_words = []

        if words is None:
            return None
        for word in words:

            # Each word needs to be at least of length 2 to form Dice bigrams. 
            if type(word) == str and "." not in word and len(word) > 2:
                valid_words.append(word)

        words_length = len(valid_words)

        # We need a minimum of 3 words to form a cluster.
        if words_length < 3:
            print(words_length, " words found. Forming clusters not possible.")

        self.words = valid_words

    def cluster(self):
        assert self.words != None and len(self.words) >= 3, "No clusters formed. Not enough words detected." 
        
        # We use a Sorensen-Dice (F-1) similairty measurement algorithm
        # using bi-grams (groups of 2 letters for each word).
        def extract_indices_dice(x, y):
            dice = SorensenDice(2)
            i, j = int(x[0]), int(y[0])
            return dice.distance(self.words[i], self.words[j])
    
        data = np.arange(len(self.words)).reshape(-1, 1)
    
        try:
            db = dbscan(data, metric=extract_indices_dice, eps=0.3, min_samples=2, algorithm='brute')
        except Exception as e:
            print(e, " Word Clustering Failed!")
            return None
    
        #if all labels are -1, DBSCAN detected no possible clusters. 
        if np.count_nonzero(db[1] == -1) == db[1].size:
            warnings.warn("There were no clusters detected. All words are unique.")
    
        return db[1]





    


