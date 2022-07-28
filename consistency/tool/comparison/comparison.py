from strsimpy.sorensen_dice import SorensenDice
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    consistency_report_pb2 as cr,
)


class Comparison:
    def __init__(self, stub, new_words, word_groups, noise_words):
        self.stub = stub
        self.new_words = new_words
        self.word_groups = word_groups
        self.noise_words = noise_words

    # format: {word : [wordgroup, distance]}
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
            comparsion_info = [self.noise_words, 1]
            for word_group in self.word_groups:
                distance = dice.distance(word, word_group.id)

                # our dice maximal threshold for considering words to be close is eps = 0.3
                if distance < comparsion_info[1] and distance < 0.3:
                    comparsion_info = [word_group, distance]
            return comparsion_info

        closest_word_groups = {}
        for word in self.new_words:
            comparison_info = find_closest_id(word)
            closest_word_groups[word] = comparison_info
        return closest_word_groups

    def generate_consistency_report(self):

        if self.word_groups == None:
            return None

        closest_word_groups = self.find_closest_word_groups()

        if closest_word_groups == None or len(closest_word_groups) < 1:
            return None
        report = cr.ConsistencyReport()
        current_variations = []
        unique_words = []
        for word in closest_word_groups:

            # Construct variation
            # Check that closest_word_group is not a noise cluster
            if closest_word_groups[word][1] != 1:
                variation = cr.ConsistencyReport.Variation()
                variation.term = word
                variation.cluster.CopyFrom(closest_word_groups[word][0])
                current_variations.append(variation)

            # Construct unique terms
            if closest_word_groups[word][0] != None and word not in closest_word_groups[word][0].word_frequency:
               unique_words.append(word)


        report.id = "consistency-report"
        report.kind = "ConsistencyReport"
        report.current_variations.extend(current_variations)
        report.unique_terms.extend(unique_words)

        return report
