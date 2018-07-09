#!/usr/bin/python 
import os, sys
def Test1(rootDir): 
	list_dirs = os.walk(rootDir) 
	for root, dirs, files in list_dirs: 
		for d in dirs: 
			print os.path.join(root, d)      
		for f in files: 
			print os.path.join(root, f) 

def Test2(rootDir): 
	for lists in os.listdir(rootDir): 
		path = os.path.join(rootDir, lists) 
		print path 
		if os.path.isdir(path): 
			Test2(path) 


if __name__ == "__main__":
	Test1(sys.argv[1])
	print "=============================="
	Test2(sys.argv[1])
