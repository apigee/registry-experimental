import unittest
import json
from parameterized import parameterized
from clustering import ClusterWords
from mock import patch
from google.protobuf.json_format import ParseDict
from google.cloud.apigeeregistry.applications.v1alpha1.consistency import (
    word_group_pb2 as wg,
)


class TestClusterWords(unittest.TestCase):
    @parameterized.expand(
        ["simple-test", "int-test", "short-test", "None-test", "dot-test"]
    )
    def test_cleaning(self, name):

        # PATCH
        # Construct mock_response
        with open("clustering_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        mock_words = obj["clean-words"][name]["words"]
        expected = obj["clean-words"][name]["cleaned-words"]
        # CALL
        clustr = ClusterWords(stub="stub", words=mock_words)
        clustr.clean_words()

        # ASSERT
        self.assertEqual(clustr.words, expected)

    # #Simple Assertions test
    @parameterized.expand(["less-than-3-words", "null-words"])
    def test_cluster_assertions(self, name):

        # PATCH
        # Construct mock_response
        with open("clustering_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        mock_words = obj["cluster-words"][name]["words"]

        # CALL
        clustr = ClusterWords(stub="stub", words=mock_words)

        # ASSERT
        with self.assertRaises(AssertionError):
            clustr.cluster()

    # Simple clustering labels test
    @parameterized.expand(["simple-test", "nonconvergence-test"])
    def test_cluster_simple(self, name):

        # PATCH
        # Construct mock_response
        with open("clustering_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)
        mock_words = obj["cluster-words"][name]["words"]

        # CALL
        clustr = ClusterWords(stub="stub", words=mock_words)
        labels = clustr.cluster()
        expected_labels = obj["cluster-words"][name]["labels"]

        # ASSERT
        self.assertListEqual(list(labels), expected_labels)

    # ID forming labels test
    @parameterized.expand(["noise-words", "duplicate-words", "unique-most-freq"])
    @patch.object(ClusterWords, "cluster")
    def test_cluster_simple(self, name, mock_cluster):

        # PATCH
        # Construct mock_response
        with open("clustering_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)

        mock_words = obj["form-ids"][name]["words"]
        mock_cluster.return_value = obj["form-ids"][name]["labels"]

        # CALL
        clustr = ClusterWords(stub="stub", words=mock_words)

        actual = clustr.create_word_groups()
        expected = []
        parsable_list = obj["form-ids"][name]["parsable-list"]
        for word_group in parsable_list:
            wrd_grp = wg.WordGroup()
            expected.append(ParseDict(word_group, wrd_grp))
        # ASSERT
        self.assertListEqual(actual, expected)

        # ID forming labels test

    @parameterized.expand(["none-words"])
    @patch.object(ClusterWords, "cluster")
    def test_cluster_simple(self, name, mock_cluster):

        # PATCH
        # Construct mock_response
        with open("clustering_test.json", "r") as myfile:
            data = myfile.read()
        obj = json.loads(data)

        mock_words = obj["form-ids"][name]["words"]
        mock_cluster.return_value = obj["form-ids"][name]["labels"]

        # CALL
        clustr = ClusterWords(stub="stub", words=mock_words)

        actual = clustr.create_word_groups()
        # ASSERT
        self.assertEqual(actual, None)


if __name__ == "__main__":
    unittest.main()
