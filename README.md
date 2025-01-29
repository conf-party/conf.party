<h1 align="center">
<img src="./out/img/logo.jpg" height="100px" style="border-radius:200px" ><br/>
conf.party
</h1>

Tracking all those conference after parties

## Contributing

If you know of a conference party happening and have all the details we'd LOVE to hear about it.

You can do so in a few ways:

1. Fork this repo and open a PR with the change needed under `./conferences`
2. Open a [new issue](https://github.com/conf-party/conf.party/issues/new) with all the details of the party
3. Let us know on [Mastodon](https://mastodon.social/@confparty) or [Bluesky](https://bsky.app/profile/conf.party)

### Running locally

```shell
cd ./src
go run . --action serve --rootdir ../ --out ../gh-pages
```

Then navigate to http://localhost:8080
