# Book Chapter Fixer

I had a book that had all the chapters split into
indidividual MP3 files with incorrect chapter numbering.

What was listed as chapter 10 was really chapter 6.
Each one was ahead by 4.

This script fixes both the filename and the IDV3 'Title'
field to apply an offset of -4.

Most of the logic specific to this naming convention is in the excellently named function `getChapter`.

Usage:
```
go build
./book-chapter-fixer "*.mp3"
# OR 
./book-chapter-fixer "Book Name (Chapter 50).mp3" "Book Name (Chapter 51).mp3"
# accepts either a glob or individual file names as args
```
