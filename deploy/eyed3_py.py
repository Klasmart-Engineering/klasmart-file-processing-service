#!/usr/bin/python3
import eyed3
import sys

if len(sys.argv) != 2:
    print("1 arg is need to spec the mp3 file")
    exit(1)

fileName = sys.argv[1]

audiofile = eyed3.load(fileName)
descriptions = []
langs = []
for comment in audiofile.tag.comments:
    descriptions.append(comment.description)
    langs.append(comment.lang)
index = 0
for description in descriptions:
    audiofile.tag.comments.remove(description, langs[index])
    index = index + 1
audiofile.tag.artist_url = ""
audiofile.tag.save()
