print("\033c")
from collections import Counter
import unittest
import json
from mock import patch
from word_extraction import ExtractWords

class FakeEntry():
    def __init__(self, word, count):
        self.word = word
        self.count = count

class FakeVocab:
    def __init__(self, schemas_list, properties_list, operations_list, parameters_list):
        self.schemas_list= schemas_list
        self.properties_list = properties_list
        self.operations_list = operations_list
        self.parameters_list = parameters_list
        self.schemas = []
        self.properties = []
        self.operations = []
        self.parameters = []
        


    def parse_lists(self):
        schemas_dict = dict(Counter(self.schemas_list))
        properties_dict = dict(Counter(self.properties_list))
        operations_dict = dict(Counter(self.operations_list))
        parameters_dict = dict(Counter(self.parameters_list))

        for key, value in schemas_dict.items():
            entry = FakeEntry(key, value)
            self.schemas.append(entry)

        for key, value in properties_dict.items():
            entry = FakeEntry(key, value)
            self.properties.append(entry)

        for key, value in operations_dict.items():
            entry = FakeEntry(key, value)
            self.operations.append(entry)

        for key, value in parameters_dict.items():
            entry = FakeEntry(key, value)
            self.parameters.append(entry)



class TestExtractWords(unittest.TestCase):

    def get_fake_vocabs(self):
        fake_vocab_file = open('consistency/tool/word_extraction_test.json')
        tests = json.load(fake_vocab_file)
        fake_vocab_file.close()
        vocabs = []
        ids = []
        all_combined = []
        for test in tests["test_data"]:
            id = test["id"]
            schemas = test["schemas"]
            properties = test["properties"]
            operations = test["operations"]
            parameters = test["parameters"]

            ids.append(id)

            print(len(all_combined))
            for word in schemas:
                all_combined.append(word)
            for word in properties:
                all_combined.append(word)
            for word in operations:
                all_combined.append(word)
            for word in parameters:
                all_combined.append(word)

            print(len(all_combined))
            fakevocab = FakeVocab(schemas, properties, operations, parameters)
            fakevocab.parse_lists()

            vocabs.append(fakevocab)
            

        return vocabs, ids, all_combined
 
    @patch.object(ExtractWords, 'extract_vocabs')
    def test_fake_vocabs(self, mock_extract_vocabs):

        longMessage = True 
        fake_vocabs, ids, combined_words = self.get_fake_vocabs()
        mock_extract_vocabs.return_value = fake_vocabs

        extrct = ExtractWords(stub = "stub", linearize=True)
        actual = extrct.get_vocabs()

        actual.sort()
        combined_words.sort()

        #print(len(combined_words))

        self.assertEquals(actual, combined_words, "failed!")


        # extrct = ExtractWords(stub = "stub", linearize=False)
        # actual = extrct.get_vocabs()

        # actual.sort()
        # combined_words.sort()

        # actual = dict(Counter(actual))
        # combined_words = dict(Counter(combined_words))

        # self.assertEquals(actual, combined_words, "failed!")







if __name__ == '__main__':
    unittest.main()