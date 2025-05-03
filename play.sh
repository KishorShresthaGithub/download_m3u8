for f in *; do
    if [[ "$f" == *.html || "$f" == *.m3u8 ]]; then
        ./main.exe -i "$f" 
        echo -en "\007"    # Run the command in the background
    fi
done