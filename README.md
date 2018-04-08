# go-business
REST APIs for businesses

Hacking around with Go lang with MySQL which can easily be deployed to Heroku.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/amarjain07/go-business
$ cd $GOPATH/src/github.com/amarjain07/go-business
$ heroku local
```

Your app should now be running on [localhost:8000](http://localhost:8000/).

You should also install [godeps] if you are going to add any dependencies to the app.

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```
