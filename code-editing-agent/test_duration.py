import unittest

from duration import to_seconds


class ToSecondsTest(unittest.TestCase):
    def test_hours_and_minutes(self):
        self.assertEqual(to_seconds("1h30m"), 5400)

    def test_seconds_only(self):
        self.assertEqual(to_seconds("45s"), 45)

    def test_hours_only(self):
        self.assertEqual(to_seconds("2h"), 7200)

    def test_hours_and_seconds(self):
        self.assertEqual(to_seconds("1h5s"), 3605)

    def test_invalid_text_raises(self):
        with self.assertRaises(ValueError):
            to_seconds("abc")

    def test_empty_raises(self):
        with self.assertRaises(ValueError):
            to_seconds("")


if __name__ == "__main__":
    unittest.main()
