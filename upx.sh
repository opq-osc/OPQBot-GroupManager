#!/bin/bash

echo $1 $2
if [ "$1" != "android" ]; then
        upx $2;
fi
exit 0;