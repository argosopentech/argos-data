#!/bin/env python

# Reads text from stdin and writes to stdout
# with one sentence per line.

# Requires pySBD
# https://spacy.io/universe/project/python-sentence-boundary-disambiguation

import sys

import pysbd

sentences = list()

segmenter = pysbd.Segmenter(language="en", clean=False)

for line in sys.stdin:
    disambiguated = segmenter.segment(line)
    sentences.extend(disambiguated)

for sentence in sentences:
    print(sentence.strip())
