#! /bin/bash

temp=`mktemp add-bookmark.XXX.log`
bm insert "$2" "$1" >"$temp" 2>&1
if [ $? -ne 0 ]
then
    cat "$temp"
fi
rm -f "$temp"

# update ~/.netsboat/config
# make sure that bm is in your path
# bookmark-cmd "~/path/to/newsboat-savebookmark.sh"
# bookmark-autopilot yet
# bookmark=interactive no
