#!/usr/bin/python

import yaml

if __name__ == "__main__":
	f = open("test.yml", "r")
	d = yaml.load(f)
	for e in d:
		print e 
