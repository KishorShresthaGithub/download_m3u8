for f in *; do
    if [[ "$f" == *.html || "$f" == *.m3u8 ]]; then
        ./download_m3u8.sh -i "$f" 
        echo -en "\007"    # Run the command in the background
    fi
done