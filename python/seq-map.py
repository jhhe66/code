#!/usr/bin/python
#-*- coding: utf-8 -*-

import yaml

if __name__ == "__main__":
	f = open("seq-map.yml", "r")
	d = yaml.load(f)
	print d
