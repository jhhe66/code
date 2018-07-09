#!/usr/bin/python
#-*- coding: utf-8 -*-


import xxhash, sys

if __name__ == "__main__":
	app_fd = open(sys.argv[1], "r")	

	for line in app_fd:
		print xxhash.xxh32(line.strip()).intdigest()
		#print line.strip()
