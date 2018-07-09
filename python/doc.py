#!/usr/bin/python
#-*- coding: utf-8 -*-

import yaml

if __name__ == "__main__":
	f = open("doc.yml", "r")
	d = yaml.load_all(f)
	for e in d:
		print e
