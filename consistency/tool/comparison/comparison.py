from strsimpy.sorensen_dice import SorensenDice
import sys


class Comparison:
    def __init__(self, stub, new_words):
        self.stub = stub
        self.new_words = new_words

    def get_word_group_ids(self):
        # TODO: retrieve the WordGroups we have stored from registry.
        # return them in an appropriate format for easy comparison
        return []

    # format: {word : (cluster_id, distance)}
    def find_word_groups(self):
        try:
            ids = self.get_word_group_ids()
        except Exception as e:
            print(e, "Getting word group IDs failed.")

        # the maximum possible score for dice distance is 1, signifying complete dissimilarity.
        if (
            ids == None
            or len(ids) < 1
            or self.new_words == None
            or len(self.new_words) < 1
        ):
            return None

        def find_closest_id(word):
            dice = SorensenDice(2)
            comparsion_info = ["UNIQUE_WORD", 1]
            for id in ids:
                distance = dice.distance(word, id)

                # our dice maximal threshold for considering words to be close is ep3 = 0.3
                if distance < comparsion_info[1] and distance < 0.3:
                    comparsion_info = [id, distance]
            return comparsion_info

        closest_word_groups = {}
        for word in self.new_words:
            comparison_info = find_closest_id(word)
            closest_word_groups[word] = comparison_info
        return closest_word_groups

    def find_word_group(self):
        # TODO: using the map returned by find_word_group,
        # if the vocab has an associated ID, form a Variation object
        # add this object into the current_variations list
        # if not, check if the vocab is a NOISE_WORD, if not, add it the
        # unique_terms list object. If there is a matching word, form a variation object
        # from the match -> means creating a wordgroup instance at the dot.
        # return the current variations object
        return

    def find_past_variations(self):
        # TODO: check the timestamp for the uploaded spec(s) and compare
        #  against the timestamp of the specs we used in the clustering step.
        # if different the variations we last computed are outdated as such are no longer current.
        # we recompute the comparison to get current variations
        return
