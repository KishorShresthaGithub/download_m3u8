#!/bin/sh

# Prepare files with no host

m3u8_files=$(ls | grep ".m3u8")
# finding max (2) files that have http urls | filtering :1 (count ) | replacing to get the file name
playlists=$(rg -e "^https?:\/\/" -m 2 -c | grep ":1" | sed -e "s/:1//g")

#get m3u8 files
# grep files with only one link

for f in $playlists; do
    url_list=$(awk '!/^#/ {print $0}' $f)
    master=$(echo "$url_list" | head -n 1)
    url=$(dirname $master)

    filename=$(basename $f .m3u8)
    echo $filename

    echo "$url_list" | tail -n +2 | awk "NF"| awk -v url="$url" '{print url "/" $0 }' >"$filename.txt"
done
