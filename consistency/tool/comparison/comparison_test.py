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
    @parameterized.expand(["simple", "unique-terms"])
    def test_find_word_groups(self, name):
        # PATCH
        # Construct inputs
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        test_suite = json.loads(data)[name]
        input_wordgroups = []
        for i in test_suite["wordgroups"]:
            wrd_grp = wg.WordGroup()
            input_wordgroups.append(ParseDict(i, wrd_grp))

        noise_group = wg.WordGroup()
        ParseDict(test_suite["noisewordgroup"], noise_group)

        # CALL
        cmprsn = Comparison(
            stub="stub",
            new_words=test_suite["words"],
            word_groups=input_wordgroups,
            noise_words=noise_group,
        )
        actual = cmprsn.find_closest_word_groups()

        expected = {}
        for word, comparison_info in test_suite["expected"].items():
            wrd_grp = wg.WordGroup()
            ParseDict(comparison_info[0], wrd_grp)
            expected[word] = [wrd_grp, comparison_info[1]]

        # ASSERT
        self.assertDictEqual(actual, expected)

    @parameterized.expand(
        [
            "none-wordgroups",
            "none-words",
            "both-null",
        ]
    )
    def test_find_word_groups_none(self, name):

        # PATCH
        # Construct inputs
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        test_suite = json.loads(data)[name]
        input_wordgroups = []
        if test_suite["wordgroups"] is not None:
            input_wordgroups.append(
                ParseDict(i, wg.WordGroup()) for i in test_suite["wordgroups"]
            )

        # CALL
        cmprsn = Comparison(
            stub="stub",
            new_words=test_suite["words"],
            word_groups=input_wordgroups,
            noise_words=None,
        )
        actual = cmprsn.find_closest_word_groups()
        self.assertIsNone(actual)

    # Comparison Report simple tests
    @parameterized.expand(["report-test-with-unqiue", "report-test-with-no-unqiue"])
    @patch.object(Comparison, "find_closest_word_groups")
    def test_report_simple(self, name, mock_find_closest_word_groups):

        dice = SorensenDice(2)
        # PATCH
        # Construct mock_response
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        test_suite = json.loads(data)[name]
        words = test_suite["words"]
        wordgroups_unparsed = test_suite["wordgroups"]
        wordgroups = []
        for wordgroup in wordgroups_unparsed:
            wrd_grp = wg.WordGroup()
            wordgroups.append(ParseDict(wordgroup, wrd_grp))

        noise_words = ParseDict(test_suite["noisegroup"], wg.WordGroup())
        report_unparsed = test_suite["expected_report"]
        expected = cr.ConsistencyReport()
        ParseDict(report_unparsed, expected)

        mock_find_closest_word_groups.return_value = {
            words[0]: [wordgroups[0], dice.distance(words[0], wordgroups[0].id)]
        }

        # CALL
        rprt = Comparison(
            stub="stub",
            new_words=words,
            word_groups=wordgroups,
            noise_words=noise_words,
        )
        actual = rprt.generate_consistency_report()

        # ASSERT
        self.assertEqual(actual, expected)

    # Comparison Report test with none words and wordgroups
    @parameterized.expand(
        ["report-test-with-none-words", "report-test-with-none-wordgroups"]
    )
    @patch.object(Comparison, "find_closest_word_groups")
    def test_report_none(self, name, mock_find_closest_word_groups):

        # PATCH
        # Construct mock_response
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        test_suite = obj[name]
        words = test_suite["words"]
        wordgroups_unparsed = test_suite["wordgroups"]
        wordgroups = None
        if wordgroups_unparsed != None:
            wordgroups = []
            for wordgroup in wordgroups_unparsed:
                wrd_grp = wg.WordGroup()
                wordgroups.append(ParseDict(wordgroup, wrd_grp))

        noise_words = ParseDict(test_suite["noisegroup"], wg.WordGroup())
        mock_find_closest_word_groups.return_value = None

        # CALL
        rprt = Comparison(
            stub="stub",
            new_words=words,
            word_groups=wordgroups,
            noise_words=noise_words,
        )
        actual = rprt.generate_consistency_report()

        # ASSERT
        self.assertIsNone(actual)

    # Comparison Report none noise_words test
    @parameterized.expand(["report-test-with-none-noise-words"])
    @patch.object(Comparison, "find_closest_word_groups")
    def test_report_simple(self, name, mock_find_closest_word_groups):

        dice = SorensenDice(2)
        # PATCH
        # Construct mock_response
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        test_suite = json.loads(data)[name]
        words = test_suite["words"]
        wordgroups_unparsed = test_suite["wordgroups"]
        wordgroups = []
        for wordgroup in wordgroups_unparsed:
            wrd_grp = wg.WordGroup()
            wordgroups.append(ParseDict(wordgroup, wrd_grp))

        noise_words = None
        report_unparsed = test_suite["expected_report"]
        expected = cr.ConsistencyReport()
        ParseDict(report_unparsed, expected)

        mock_find_closest_word_groups.return_value = {
            words[0]: [wordgroups[0], dice.distance(words[0], wordgroups[0].id)]
        }

        # CALL
        rprt = Comparison(
            stub="stub",
            new_words=words,
            word_groups=wordgroups,
            noise_words=noise_words,
        )
        actual = rprt.generate_consistency_report()

        # ASSERT
        self.assertEqual(actual, expected)

    # Comparison report with unique and existing.
    @parameterized.expand(["report-test-unqiue-existing"])
    @patch.object(Comparison, "find_closest_word_groups")
    def test_report_simple(self, name, mock_find_closest_word_groups):

        dice = SorensenDice(2)
        # PATCH
        # Construct mock_response
        with open("comparison_test.json", "r") as myfile:
            data = myfile.read()
        test_suite = json.loads(data)[name]
        words = test_suite["words"]
        wordgroups_unparsed = test_suite["wordgroups"]
        wordgroups = []
        for wordgroup in wordgroups_unparsed:
            wrd_grp = wg.WordGroup()
            wordgroups.append(ParseDict(wordgroup, wrd_grp))

        report_unparsed = test_suite["expected_report"]
        expected = cr.ConsistencyReport()
        ParseDict(report_unparsed, expected)
        closest_word_groups = {}
        closest_word_groups[words[0]] = [
            wordgroups[0],
            dice.distance(words[0], wordgroups[0].id),
        ]

        noise_words = ParseDict(test_suite["noisegroup"], wg.WordGroup())
        closest_word_groups[words[1]] = [
            noise_words,
            dice.distance(words[1], noise_words.id),
        ]
        mock_find_closest_word_groups.return_value = closest_word_groups

        # CALL
        rprt = Comparison(
            stub="stub",
            new_words=words,
            word_groups=wordgroups,
            noise_words=noise_words,
        )
        actual = rprt.generate_consistency_report()

        # ASSERT
        self.assertEqual(actual, expected)

if __name__ == "__main__":
    unittest.main()
