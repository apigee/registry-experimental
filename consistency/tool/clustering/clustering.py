from tkinter.tix import StdButtonBox
import numpy as np
import pandas as pd
from sklearn.cluster import dbscan
from strsimpy.sorensen_dice import SorensenDice
from sklearn.cluster import dbscan
from consistency.tool.extraction.word_extraction import ExtractWords
from collections import Counter
import consistency.rpc.google.cloud.apigeeregistry.applications.v1alpha1.consistency.word_group_pb2 as wg 
from google.protobuf.json_format import ParseDict
class ClusterWords:
    def __init__(self, stub, useSemanticSimilarity = False):
      self.useSemanticSimilarity = useSemanticSimilarity
      self.stub = stub


    def clean_words(self):
        stub = self.stub 
        extrct = ExtractWords(stub)
        words = extrct.get_vocabs()
        for word in words:
            if not isinstance(word, str) or len(word)<2:
                words.remove(word)
        return words

    def cluster(self):
        # place holder for when we explore semantic similarity
        useSemanticSimilarity = self.useSemanticSimilarity

        dice = SorensenDice(2)

        # extract indices
        def extract_indices_dice(x, y):
            i, j = int(x[0]), int(y[0])     
            return dice.distance(words[i], words[j])

        words = self.clean_words()
        words = words.to_numpy()
        words = np.arange(len(words)).reshape(-1, 1)
        db = dbscan(words, metric=extract_indices_dice, eps=.3, min_samples=2, algorithm='brute')
        labels = db[1]

        return words, labels

    def create_word_groups(self):
        words, labels = self.cluster()
        temp_word_groups  = {}
        for j in range(len(words)):
            word_label = labels[1][j]

            if word_label in temp_word_groups:
                temp_word_groups[word_label].append(words[j])
                
            else:
                temp_word_groups[word_label] = [words[j]]

        word_groups = {}
        for k in temp_word_groups.keys():
            value = temp_word_groups[k]
            if len(Counter(value).most_common())  > 1:
                word, count =Counter(value).most_common()[0] 
                if count > 1:
                    word_groups[word] = value
                else:
                    value.sort()
                    word_groups[value[0]] = value

        temp_word_groups.clear()

        return word_groups

    def vocabulary_upload(self):
        word_groups = self.create_word_groups()
        wordGroup = wg.WordGroup()
        for key, val in word_groups.items():
            word_group = {}
            word_group["id"] = key
            word_group["kind"] = "kind"
            word_group["word_frequency"] = dict(Counter(val))
            ParseDict(word_group,  wordGroup)

            #to do 
            # upload the parsed wordGroup to registry 
            # test by patching ExtractVocabs
            # 


        



    


