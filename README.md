# PhotoManager

Little CLI based photo organiser that was made because I'm too lazy to organize my photos.

## Note: only works on Windows at the moment as it uses the metadata of the image to determine the folder to move it to. 

How to run:
```shell
go build .
photoManager.exe -i ./inputFolder -o ./outputFolder -r -t jpeg,jpg,raw,dng,raf
```

flags
```
-i: input folder
-o: output folder
-r: recursive (will search subfolders)
-t: file types to search for (comma separated)
```