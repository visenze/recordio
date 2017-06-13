import unittest

import recordio
import cPickle as pickle

class TestStringMethods(unittest.TestCase):
    def test_write_read(self):
        w = recordio.writer("/tmp/record_0")
        w.write(pickle.dumps("1"))
        w.write("2")
        w.write("")
        w.close()
        w = recordio.writer("/tmp/record_1")
        w.write("3")
        w.write("4")
        w.write("")
        w.close()

        r = recordio.reader("/tmp/record_*")
        self.assertEqual(pickle.loads(r.read()), "1")
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
