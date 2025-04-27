#!/bin/sh

## usage ./script.sh -i input.m3u8

# input m3u8 path

skip=0

while getopts "si:" arg; do
    case $arg in
    i) file=$OPTARG ;;
    s) skip=1 ;;
    esac
done

foldername=$(echo $file | awk -F'.' '{print $1}' -)
# create input name and create directory using that name
mkdir -p $foldername
cp "$file" "$foldername/$file" 

# prepare links for download using aria2c
# check if single line or not
lineCount=$(awk 'END{print NR}' $file)

if [ $lineCount -eq 1 ]; then
    # Single line: extract links and save to a file
    awk -v RS=" " -v FS="/" '/^https?:\/\// {print $0; print "    out="$NF}' "$file" >"$foldername/$foldername.links.txt"
    awk -v RS=" " -v FS="/" '/^https?:\/\// {print "file \047" $NF "\047"}' "$file" >"$foldername/$foldername.file.txt"
else
    # Download Link generation step
    # get all lines that do not start with # and export them
    awk -F'/' '/^https?:\/\// { split($NF,a,"?"); print $0 "\n  out=" a[1] }' "$file" >"$foldername/$foldername.links.txt"
    # File process input generation step
    awk -F'/' '/^https?:\/\// { split($NF,a,"?"); print "file \047" a[1] "\047" }' "$file" >"$foldername/$foldername.file.txt"
fi

cd $foldername

if [ "$skip" -eq 0 ]; then
    aria2c -U "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0" \
        --header="Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8" \
        --header="Accept-Language: en-US,en;q=0.5" \
        --header="Accept-Encoding: gzip, deflate, br" \
        --header="Connection: keep-alive" \
        --header="Upgrade-Insecure-Requests: 1" \
        --header="Sec-Fetch-Dest: document" \
        --header="Sec-Fetch-Mode: navigate" \
        --header="Sec-Fetch-Site: none" \
        --header="Sec-Fetch-User: ?1" \
        -j 10 -s 10 -x 10 -c -i "$foldername.links.txt"
fi

isPng=$(ls | head -n 1)

if file "$isPng" | grep -q "PNG image data"; then

    mkdir -p backup
    for f in *; do
        cp "$f" "backup/$f"
        tail -c +9 "$f" >temp && mv temp "$f"
    done
fi

# after download finish merge the files using ffmpeg
ffmpeg -f concat -safe 0 -i "$foldername.file.txt" -c copy -crf 22 $foldername.mp4

mv $foldername.mp4 ../$foldername.mp4

cd ..

rm "$file"

echo "Task Complete Press Enter to continue..."
#read # Waits for the user to press Enter
