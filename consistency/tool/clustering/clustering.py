import numpy as np
import pandas as pd
from sklearn.cluster import dbscan
from strsimpy.sorensen_dice import SorensenDice
from sklearn.cluster import dbscan
from collections import Counter
import warnings
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    word_group_pb2 as wg,
)

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
        assert (
            self.words != None and len(self.words) >= 3
        ), "No clusters formed. Not enough words detected."

        # We use a Sorensen-Dice (F-1) similairty measurement algorithm
        # using bi-grams (groups of 2 letters for each word).
        def extract_indices_dice(x, y):
            dice = SorensenDice(2)
            i, j = int(x[0]), int(y[0])
            return dice.distance(self.words[i], self.words[j])

        data = np.arange(len(self.words)).reshape(-1, 1)

        try:
            _, labels = dbscan(
                data,
                metric=extract_indices_dice,
                eps=0.3,
                min_samples=2,
                algorithm="brute",
            )
        except Exception as e:
            print(e, " Word Clustering Failed!")
            return None

        # if all labels are -1, DBSCAN detected no possible clusters.
        if np.count_nonzero(labels == -1) == labels.size:
            warnings.warn("There were no clusters detected. All words are unique.")

        return labels

    def create_word_groups(self):

        labels = np.array(self.cluster())

        temp_dict = {}
        if not labels.size or not self.words:
            return None
        for j in range(len(self.words)):
            word_label = labels[j]

            if word_label in temp_dict:
                temp_dict[word_label].append(self.words[j])

            else:
                temp_dict[word_label] = [self.words[j]]

        word_groups = []

        def find_clusterID(similar_list):
            if len(similar_list) == 0 or similar_list == None:
                return None
            counts = Counter(similar_list)
            max_count = counts.most_common(1)[0][1]
            most_freq_words = [
                value for value, count in counts.most_common() if count == max_count
            ]
            most_freq_words.sort()

            return most_freq_words[0]

        for k in temp_dict.keys():
            word_group = wg.WordGroup()
            word_group.kind = "WordGroup"
            map = dict(Counter(temp_dict[k]))
            for word, frequency in map.items():
                word_group.word_frequency[word] = frequency

            if k == -1:
                word_group.id = "NOISE_WORDS"

            else:
                similar_words = temp_dict[k]
                id = find_clusterID(similar_words)

                if id is None:
                    continue
                word_group.id = id
            word_groups.append(word_group)

        temp_dict.clear()
        return word_groups
