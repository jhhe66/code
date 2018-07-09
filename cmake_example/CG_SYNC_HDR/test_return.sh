#!/bin/sh

cg_sync_hdr 0x0A40 | tee 0x0A40.log

if [ $? -eq 0 ]
then
	echo "success."
else
	echo "error."
fi

exit 0

