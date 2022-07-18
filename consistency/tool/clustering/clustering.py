
from cmath import nan
import numpy as np
import pandas as pd
from sklearn.cluster import dbscan
from strsimpy.sorensen_dice import SorensenDice
from sklearn.cluster import dbscan
from word_extraction import ExtractWords
from collections import Counter
import consistency.rpc.google.cloud.apigeeregistry.applications.v1alpha1.consistency.word_group_pb2 as wg

from google.protobuf.json_format import ParseDict
class ClusterWords:
    def __init__(self, stub):
      self.stub = stub

    def get_words(self):
        extrct = ExtractWords(stub="stub")
        words = extrct.get_vocabs()

        if words is None:
            return None
        return words

    def clean_words(self):
        words = self.get_words()

        if words is None:
            return None
        for word in words:
            if type(word) != str or len(word) < 2 or word != word:
                words.remove(word)
        words_length = len(words)

        if words_length < 3:
            print("Only ", words_length, " words found. Forming clusters not possible.")
            return None

        return words

    def cluster(self):

        words = self.clean_words()

        if words is None:
            return words

        def extract_indices_dice(x, y):

            dice = SorensenDice(2)
            nonlocal words
            i, j = int(x[0]), int(y[0])  
            return dice.distance(words[i], words[j])

        data = np.arange(len(words)).reshape(-1, 1)

        try:
            db = dbscan(data, metric=extract_indices_dice, eps=0.3, min_samples=2, algorithm='brute')
        except Exception as e:
            print(e, " Word Clustering Failed!")
            return None

        return words, db[1]
    def create_word_groups(self):

        words, labels = self.cluster()

        _word_groups  = {}
        for j in range(len(words)):
            word_label = labels[j]

            if word_label in _word_groups:
                _word_groups[word_label].append(words[j])
                        
            else:
                _word_groups[word_label] = [words[j]]

        word_groups = {}

        for k in _word_groups.keys():
            similar_words = _word_groups[k]
            if len(Counter(similar_words).most_common())  > 1:
                word, count =Counter(similar_words).most_common()[0] 
                if count > 1:
                    word_groups[word] = similar_words
                else:
                    similar_words.sort()
                    word_groups[similar_words[0]] = similar_words

        return word_groups

    # def vocabulary_upload(self):
    #     _word_groups = self.create_word_groups()
    #     wordGroup = wg.WordGroup()
    #     for key, val in _word_groups.items():
    #         word_group = {}
    #         word_group["id"] = key
    #         word_group["kind"] = "kind"
    #         word_group["word_frequency"] = dict(Counter(val))
    #         ParseDict(word_group,  wordGroup)

    #         #to do 
    #         # upload the parsed wordGroup to registry 
    #         # test by patching ExtractVocabs
    #         # 


        



    


