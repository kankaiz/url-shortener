# url-shortener
this is a mini server for URL shortener

run `go get github.com/zkkzero/url-shortener` to get the code 

The server has been deployed on AWS.

POST http://52.62.156.51:3008/ with "url" and the value to get a unique short url

For example:

`curl --data "url=https://www.broadsheet.com.au/" 52.62.156.51:3008`
then it will return 34nIa0

we can now get http://52.62.156.51:3008/34nIa0 to redirect to our original url
