import unittest
import json
from mock import patch
from parameterized import parameterized
from clustering import ClusterWords

class TestClusterWords(unittest.TestCase):

    @parameterized.expand([
        ('simple', ['abandon', 'abandoning', 'Abort', 'abort', 'aborted', 'About', 'about', 'Above',
                                        'Absence', 'Absent', 'absentee', 'Absolute', 'absolutely', 'Abstain', 'Abuse',
                                        'abusive', 'Accelerator', 'accelerator', 'accelerators', 'Accelerators',
                                        'accept', 'acceptable', 'Accepted', 'accepts', 'Access', 'access', 'Accessed',
                                        'Accesses', 'Accessibility'])
    ])

    @patch.object(ClusterWords, 'get_words')
    def test_clustering(self, name, expected, mock_get_words):

        # PATCH
        # Construct mock_response
        
        mock_get_words.return_value = ['abandon', 'abandoning', 'Abort', 'abort', 'aborted', 'About', 'about', 'Above',
                                        'Absence', 'Absent', 'absentee', 'Absolute', 'absolutely', 'Abstain', 'Abuse',
                                        'abusive', 'Accelerator', 'accelerator', 'accelerators', 'Accelerators',
                                        'accept', 'acceptable', 'Accepted', 'accepts', 'Access', 'access', 'Accessed',
                                        'Accesses', 'Accessibility', 'a']
 
        #CALL
        clustr = ClusterWords(stub = "stub")
        actual = clustr.clean_words()
        print(actual)
        # ASSERT
        self.assertListEqual(actual, expected)
    # @parameterized.expand([
    #     ('simple', ["ab",  "ab"])
    # ])

    # @patch.object(ClusterWords, 'get_words')
    # def test_clean_Words(self, name, expected, mock_get_words):

    #     # PATCH
    #     # Construct mock_response
        
    #     mock_get_words.return_value = ["ab", "ab", "a"]
 
    #     #CALL
    #     clustr = ClusterWords(stub = "stub")
    #     actual = clustr.clean_words()
    #     # ASSERT
    #     self.assertListEqual(actual, expected)
 
 
if __name__ == '__main__':
    unittest.main()