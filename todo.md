## Algorithm things to do

[] download m3u8 files

[] download all files in m3u8 using aria2c
options to use during aria2c
parallel: -j  
 splits : -s
connection per: -x
directory -d
input -i
continue -c

[] create files.txt using m3u8 file

```sh
    sed -n "/^[^#]/s|.\*/|file |p" file.m3u8 > basenames_with_file.txt

    awk -F'/' '!/^#/ {split($NF, a, "?"); print "file " a[1]}' urls.txt > output.txt

```

[] use ffmpeg to merge all files

    ffmpeg -f concat -i test.txt  -c:v hevc_nvenc -pix_fmt yuv420p output.mp4

openssl enc -aes-128-cbc -d -in 322890.ts -out decrypted.ts -K $(xxd -p ../enc.ts | tr -d '\n') -iv a56d08bab69c1776afb9a639817c9a15 | tail -c +17 > fixed_decrypted.ts

## get iv
head -c 16 ../enc.ts > iv.bin

## remaining
tail -c +17 ../enc.ts > enc_no_iv.ts

## try this as well
00000000000000000000000000000000

for $file in *.ts; openssl enc -aes-128-cbc -d -in $file -out ../decrypt/$file -K $(xxd -p ../enc.ts | tr -d '\n') -iv a56d08bab69c1776afb9a639817c9a15; done 


## batch download
for f in *.m3u8; do ./download_m3u8.sh -i $f & \ ; done


```s
# Step 1: Clear previous list
> file_list.txt

# Step 2: Process files and build concat list
while read -r line; do
    if [[ $line == file* ]]; then
        filename=$(echo "$line" | awk '{print $2}' | sed "s/['\"]//g")
        base="${filename%.*}"
        tsfile="${base}.ts"

        if ffmpeg -v error -i "$tsfile" -f null - 2>/dev/null; then
            echo "file '$tsfile'" >> file_list.txt
            echo "[OK] Added $tsfile"
        else
            echo "[SKIP] $tsfile is missing or corrupted."
        fi
    fi
done < ipzz624.playlist.txt
```