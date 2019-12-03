# mp3copy

Music players supporting USB drives and SD cards often only support the FAT file system and play media in the order files were written to the drive.

This utility copies music files (MP3 currently, FLAC and/or AAC can be added if there is demand) and can sort by:

| option | description        |
| ------ | -----------        |
| track  | track id           |
| artist | artist name        |
| album  | album name         |
| song   | song name          |
| file   | file name          |
| date   | last modified date |
| random | random order       |

Sort order is determined by (in order of precedence):

* .mp3copy file in directory
* .mp3copy file in parent directory with *children=true*
* command line arguments
* built-in default: artist name, album name, track id; all using ascending order

## .mp3copy

Place this file in any directories to control the sort order.

```bash
# Comma separated list of sort criteria in order of precedence.
# Order can be specified by suffixing `:a` (ascending, default)
# or `:d` (descending)
sort = artist:a, album:a, track:a

# true (default) if these settings should apply to child directories
# that do not contain their own .mp3copy files
children = true
```
