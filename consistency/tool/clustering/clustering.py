
from cmath import nan
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





    


