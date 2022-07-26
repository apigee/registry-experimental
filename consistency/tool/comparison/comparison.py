from strsimpy.sorensen_dice import SorensenDice
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    consistency_report_pb2 as cr,
)


class Comparison:
    def __init__(self, stub, new_words, word_groups=None, noise_words_index = None):
        self.stub = stub
        self.new_words = new_words
        self.word_groups=word_groups
        self.noise_words_index  = noise_words_index 

    # format: {word : (cluster_id, distance)}
    def find_closest_word_groups(self):

        # the maximum possible score for dice distance is 1, signifying complete dissimilarity.
        if (
            self.word_groups == None
            or len(self.word_groups) < 1
            or self.new_words == None
            or len(self.new_words) < 1
        ):
            return None

        def find_closest_id(word):
            dice = SorensenDice(2)
            comparsion_info = ["POSSIBLE_UNIQUE_WORD", 1]
            for word_group in self.word_groups:
                distance = dice.distance(word, word_group.id)

                # our dice maximal threshold for considering words to be close is ep3 = 0.3
                if distance < comparsion_info[1] and distance < 0.3:
                    comparsion_info = [word_group, distance]
            return comparsion_info

        closest_word_groups = {}
        for word in self.new_words:
            comparison_info = find_closest_id(word)
            closest_word_groups[word] = comparison_info
        return closest_word_groups

    def find_word_group(self):

        
        if self.word_groups == None:
            return None
        unique_words = []
        variations = []
        closest_word_groups = self.find_closest_word_groups()
        
        for word in closest_word_groups:
            if closest_word_groups[word] == ["UNIQUE_WORD", 1] and word not in self.word_groups[self.noise_words_index]:
                unique_words.append(word)
            variation = cr.ConsistencyReport.Variation()
            variation.term = word
            variation.cluster = 0
            

    #     return

    def find_past_variations(self):
        # TODO: check the timestamp for the uploaded spec(s) and compare
        #  against the timestamp of the specs we used in the clustering step.
        # if different the variations we last computed are outdated as such are no longer current.
        # we recompute the comparison to get current variations
        return
