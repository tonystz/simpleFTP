SimpleFTP client wrote by Golang
Usage:
```
ftp.New("10.64.70.73")
ftp.GetFile("/pub/atop-1.27-3.tar.gz")
ftp.Close()
```
or
```
ftp.New("10.64.70.73","username","passwd")
ftp.GetFile("/pub/atop-1.27-3.tar.gz")
ftp.Close()
```
