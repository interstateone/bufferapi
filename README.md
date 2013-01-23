bufferapi
============

A little wrapper for the Buffer API in Go.

Right now it will issue requests if you already have an auth token. You should be able to use Google's goauth2 package to let a user authenticate ([example](https://code.google.com/p/goauth2/source/browse/oauth/example/oauthreq.go)), and then pass in the transport and auth token.

There are some basic tests to run with `go test` that make sure requests respond correctly, given a `config.json` file that contains the following:

    {
      "ClientId": "your client id",
      "ClientSecret": "your client secret",
      "AuthToken": "your auth token (in your emails from buffer)"
    }
