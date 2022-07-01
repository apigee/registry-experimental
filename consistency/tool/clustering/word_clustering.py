import word_extracting
import grpc
import numpy as np
import pandas as pd
from sklearn.cluster import dbscan
from strsimpy.sorensen_dice import SorensenDice
from sklearn.cluster import dbscan
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from metrics import vocabulary_pb2
from google.cloud.apigeeregistry.v1 import registry_service_pb2_grpc
from google.cloud.apigeeregistry.v1 import registry_service_pb2
from metrics import vocabulary_pb2
from collections import Counter



class ClusterWords:

    def __init__(self, words, linearized = True):
       self.words = words
       self.linearized = linearized
       self.eps = .3 
       self.dbscan_min_samples = 2
       self.dice = SorensenDice(2)
       
    def linearize_words(self):
        linearized = self.linearized
        words = self.words
        data = []
        if not linearized:
            for key in words:
                for _ in range(words[key]):
                    data.append(key)
        else:
            data = words
        return data

    def extract_indices_dice(self, x, y, data):
        i, j = int(x[0]), int(y[0])     # extract indices
        return self.dice.distance(data[i], data[j])

    def extract_words (self):
        data = self.linearize_words()
        data = np.arange(len(data)).reshape(-1, 1)


        db = dbscan(data, metric=self.extract_indices_dice, eps=self.eps, min_samples = self.dbscan_min_samples, algorithm='brute')
        labels = db[1]
        counts = Counter(labels)
        clusters = dict()

        for i, n in enumerate(labels):
            if counts[n]>1:
                clusters.setdefault(n,[]).append(i)

        for key in clusters:
            similar_words = list(clusters[key])
            cluster_counts = Counter(similar_words)
            new_key = max(similar_words, key=cluster_counts.get)
            clusters[new_key] = clusters.pop(key)

        return clusters


