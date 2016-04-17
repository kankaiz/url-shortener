# url-shortener
this is a mini server for URL shortener

The server has been deployed on AWS.

POST http://52.62.156.51:3008/ with "url" and the value to get a unique short url

For example:

curl --data "url=https://www.broadsheet.com.au/" 52.62.156.51:3008
then return 55B3yB

we can now get http://52.62.156.51:3008/55B3yB to redirect to our original url
