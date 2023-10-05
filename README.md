# PhotoManager

Little CLI based photo organiser that was made because I'm too lazy to organize my photos into folders organised by date.

## Note: Should work cross-platform now due to exif support.

How to run:
```shell
go build .
photoManager.exe -i ./inputFolder -o ./outputFolder -r -t jpeg,jpg,raw,dng,raf -threads 10
```

flags
```
-i: input folder (required)
-o: output folder (required)
-r: recursive (will search subfolders, optional)
-t: file types to search for (comma separated, required)
-threads: number of threads to use (defaults to 5, optional)
```
