#!/usr/bin/python

import sys
from zlib import crc32


if __name__ == "__main__":

	appkey_fd = open(sys.argv[1], 'r')

	for line in appkey_fd:
		print crc32(line.strip()) & 0xFFFFFFFF
