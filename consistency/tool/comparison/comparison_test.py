import unittest
import json
from mock import patch
from comparison import Comparison
from parameterized import parameterized
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    word_group_pb2 as wg,
)
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    consistency_report_pb2 as cr,
)
from google.protobuf.json_format import ParseDict
from strsimpy import SorensenDice
class TestComparison(unittest.TestCase):
    # @parameterized.expand(
    #     [
    #         "simple",
    #         "none-wordgroups",
    #         "none-words",
    #         "both-null",
    #     ]  # , "unique-terms", "both-null"]
    # )
    # def test_find_word_groups(
    #     self,
    #     name,
    # ):

    #     # PATCH
    #     # Construct mock_response
    #     with open("comparison_test.json", "r") as myfile:
    #         data = myfile.read()
    #     obj = json.loads(data)
    #     test_suite = obj[name]
    #     input_wordgroups = []
    #     if test_suite["wordgroups"] == None or test_suite["words"] == None:
    #         return None
    #     for i in test_suite["wordgroups"]:
    #         wrd_grp = wg.WordGroup()
    #         if i != None:
    #             input_wordgroups.append(ParseDict(i, wrd_grp))
    #         else:
    #             input_wordgroups.append(None)

    #     # CALL
    #     cmprsn = Comparison(
    #         stub="stub", new_words=test_suite["words"], word_groups=input_wordgroups
    #     )
    #     actual = cmprsn.find_closest_word_groups()
    #     expected = {}
    #     for word, comparison_info in test_suite["expected"].items():
    #         wrd_grp = wg.WordGroup()
    #         ParseDict(comparison_info[0], wrd_grp)
    #         expected[word] = [wrd_grp, comparison_info[1]]

    #     # ASSERT
    #     self.assertDictEqual(actual, expected)

    # @parameterized.expand(["unique-terms"])  # , "unique-terms", "both-null"]
    # def test_find_word_groups_unique(
    #     self,
    #     name,
    # ):

    #     # PATCH
    #     # Construct mock_response
    #     with open("comparison_test.json", "r") as myfile:
    #         data = myfile.read()
    #     obj = json.loads(data)
    #     test_suite = obj[name]
    #     input_wordgroups = []
    #     if test_suite["wordgroups"] == None or test_suite["words"] == None:
    #         return None
    #     for i in test_suite["wordgroups"]:
    #         wrd_grp = wg.WordGroup()
    #         if i != None:
    #             input_wordgroups.append(ParseDict(i, wrd_grp))
    #         else:
    #             input_wordgroups.append(None)

    #     # CALL
    #     cmprsn = Comparison(
    #         stub="stub", new_words=test_suite["words"], word_groups=input_wordgroups
    #     )
    #     actual = cmprsn.find_closest_word_groups()
    #     expected = test_suite["expected"]
    #     self.assertDictEqual(actual, expected)


    # Compparison Report test
    @parameterized.expand(["simple-report-test"])
    @patch.object(Comparison, "find_closest_word_groups")
    def test_report_simple(self, name, mock_find_closest_word_groups):
        

        dice = SorensenDice(2)
        # PATCH
        # Construct mock_response
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        test_suite = obj[name]
        words = test_suite["words"]
        wordgroups = test_suite["wordgroups"]
        noisegroup = test_suite["noisegroup"]

    

        # CALL

  
        # ASSERT
        #self.assertEqual(actual, expected)

if __name__ == "__main__":
    unittest.main()
