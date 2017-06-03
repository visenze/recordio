import unittest

import recordio


class TestStringMethods(unittest.TestCase):
    def test_upper(self):
        w = recordio.writer("/tmp/record_0")
        w.write("1")
        w.write("2")
        w.write("")
        w.close()
        w = recordio.writer("/tmp/record_1")
        w.write("3")
        w.write("4")
        w.write("")
        w.close()

        r = recordio.reader("/tmp/record_*")
        self.assertEqual(r.read(), "1")
        self.assertEqual(r.read(), "2")
        self.assertEqual(r.read(), "")
        self.assertEqual(r.read(), "3")
        self.assertEqual(r.read(), "4")
        self.assertEqual(r.read(), "")
        self.assertEqual(r.read(), None)
        self.assertEqual(r.read(), None)
        r.close()


if __name__ == '__main__':
    unittest.main()
