#!/usr/bin/python

from __future__ import print_function
import GeoIP

gi = GeoIP.open("GeoIP.dat", GeoIP.GEOIP_STANDARD)

print(gi.country_code_by_name("yahoo.com"))
print(gi.country_code_by_name("ifeng.com"))
