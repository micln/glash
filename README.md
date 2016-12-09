# glash
能指定文件目录的下载器


## Install

```sh
git clone https://github.com/micln/glash.git
cd glash
go build
```

##Usage



```sh
$ ./glash

Usage of ./glash:
  -d	start http server.
	listen on http://0.0.0.0:10087?url=xx&filename=xx
  -f string
    	filename. the file will be download to "Downloads/20161209/filename"
  -json
    	parse the file as json
  -t	parse the file like 'url(\t)filename'
  -tool string
    	choose the download tools. support [curl, aria2]

```

## Example

```
# use httpd daemon
./glash -d

# use json file
./glash -f ~/test.json -json -tool aria2

# use aria2c
./glash -f ~/test.json -json -tool aria2

```