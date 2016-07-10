#!/bin/sh

set -e

echo "Syncing likes 💕"

docker-compose run --rm botmachine syncicecream

echo "Making an ice cream 🍦"

NOW=$(date +"%d_%m_%Y")
IMAGE_PATH_SVG=/data/$NOW.svg
IMAGE_PATH_PNG=/data/$NOW.png

function run {
 docker-compose run --rm botmachine sh -c "$1"
}

run "makeicecream generate > $IMAGE_PATH_SVG"
run "rsvg-convert -d 1200 -p 1200 -w 1080 -h 1080 $IMAGE_PATH_SVG > $IMAGE_PATH_PNG"

echo "Tweeting 🐦"

ICECREAM_ID=$(docker-compose run --rm botmachine makeicecream id | tr -d "['\r\n']")

run "tweeticecream $ICECREAM_ID $IMAGE_PATH_PNG"
