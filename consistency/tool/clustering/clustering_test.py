import unittest
import json
from mock import patch
from parameterized import parameterized
from clustering import ClusterWords

class TestClusterWords(unittest.TestCase):

    with open('/home/gelaw/work-stuff/gocode/src/registry-experimental/consistency/tool/clustering/clustering_test.json', 'r') as myfile:
            data=myfile.read()
    obj = json.loads(data)
    names = []
    expected_values = []

    for name in obj:
        names.append(name)
        expected_values.append(obj[name][0]["clusters"])



    @parameterized.expand([
        (names[0],{'abandon': ['abandon', 'abandoning'], 'Abort': ['Abort', 'abort', 'aborted'], 
        'About': ['About', 'about'], 'Above': ['Above', 'Abstain', 'Abuse', 'Accepted', 'Accessibility', 'absentee', 'abusive'],
        'Absence': ['Absence', 'Absent'], 'Absolute': ['Absolute', 'absolutely'],
        'Accelerator': ['Accelerator', 'Accelerators', 'accelerator', 'accelerators'],
        'accept': ['accept', 'acceptable', 'accepts'], 'Access': ['Access', 'Accessed', 'Accesses', 'access']})
    ])

    @patch.object(ClusterWords, 'get_words')
    def test_clustering(self, name, expected, mock_get_words):

        # PATCH
        # Construct mock_response
        
        mock_get_words.return_value = ["abandon", "abandoning", "Abort", "abort", "aborted", "About", "about", "Above",
            "Absence","Absent", "absentee", "Absolute", "absolutely", "Abstain", "Abuse",
            "abusive", "Accelerator", "accelerator", "accelerators", "Accelerators",
            "accept", "acceptable", "Accepted", "accepts", "Access", "access", "Accessed",
            "Accesses", "Accessibility"]
 
        #CALL
        clustr = ClusterWords(stub = "stub")
        actual = clustr.create_word_groups()
        print(actual)
        # ASSERT
        
        self.assertDictEqual(dict(sorted(actual.items())), dict(sorted(expected.items())))
 
 
if __name__ == '__main__':
    unittest.main()