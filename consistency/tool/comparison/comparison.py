class Comparison:
    def __init__(self, stub, spec_vocabs):
        self.stub = stub
        self.spec_vocabs = spec_vocabs

    def get_word_group_ids (self):
        # TODO: retrieve the WordGroups we have stored from registry. 
        # return them in an appropriate format for easy comparison 
        return
    def find_word_group (self):
        # TODO: using the list of IDs returned by get_word_group_ids, 
        # try to find the closest wordgroup for each spec vocab.
        # if successful, save it in the format vocab:id  else vocab :None
        # return a list of these
        return

    def find_word_group (self):
        # TODO: using the map returned by find_word_group,
        # if the vocab has an associated ID, form a Variation object 
        # add this object into the current_variations list
        # if not, check if the vocab is a NOISE_WORD, if not, add it the 
        # unique_terms list object. If there is a matching word, form a variation object 
        # from the match -> means creating a wordgroup instance at the dot. 
        # return the current variations object
        return

    def find_past_variations (self):
        # TODO: check the timestamp for the uploaded spec(s) and compare
        #  against the timestamp of the specs we used in the clustering step. 
        # if different the variations we last computed are outdated as such are no longer current. 
        # we recompute the comparison to get current variations
        return