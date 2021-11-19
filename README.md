```
$ go install github.com/vzvu3k6k/twitter_avatar_downloader@latest
$ twitter_avatar_downloader cnn nytimes nasa
$ ls
cnn.jpg  nasa.jpg  nytimes.png
$ twitter_avatar_downloader thisuserdoesnotexist
2021/11/20 08:00:00 get thisuserdoesnotexist
2021/11/20 08:00:30 context deadline exceeded
```