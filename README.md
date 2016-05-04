# url-shortener
The web app has been dockerised!

after install docker and docker-compose properly

simply run `docker-compose up`

POST http://localhost:3008/ with "url" and the value to get a unique short url

For example:

`curl --data "url=https://www.broadsheet.com.au/" http://localhost:3008`
then it will return 34nIa0

we can now get http://localhost:3008/34nIa0 to redirect to our original url
