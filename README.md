# Go Strava Client

This is a barebones and hardly complete Go client for the [Strava API v3](http://strava.github.io/api). It was written mainly as an exercise to play around with Go and therefore only supports a small set of operations and includes a subset of fields on model objects.

## Authentication

The Strava API requires authentication and uses OAuth. The OAuth flow is not implemented in this API, but access can still be gained using a developer access token which can be easily obtained by any Strava user. The Strava API page has details. *Don't share your access token*.

## Example CLI App

A sample command line app is included that can list activities or segments. Output is in CSV (though custom delimiters are supported with with `--delimiter`).

### Building the App

```
go install github.com/alecholmes/strava
```

### Get All Activities

```
STRAVA_ACCESS_TOKEN=your_private_token
$GOPATH/bin/strava --accessToken $STRAVA_ACCESS_TOKEN
```

Or, to get activities after a given id:

```
STRAVA_ACCESS_TOKEN=your_private_token
$GOPATH/bin/strava --accessToken $STRAVA_ACCESS_TOKEN --afterId 212147000
```

### Get Segment Details for Activities

All segments for activities can be printed instead of activity summaries. This is done by including the `--segment` flag.

```
STRAVA_ACCESS_TOKEN=your_private_token
$GOPATH/bin/strava --accessToken $STRAVA_ACCESS_TOKEN --afterId 212147000 --segments
```

All times for a specific segment, and the date of the effort:

```
STRAVA_ACCESS_TOKEN=your_private_token

# Dump everything into all_segments to allow reuse, since fetching is relatively slow
$GOPATH/bin/strava --accessToken $STRAVA_ACCESS_TOKEN --delimiter $'\t' --segments > all_segments

grep "\tHawk Hill\t" all_segments | cut -d $'\t' -f 8,10
```

## Bugs

* GetActivity blows up if a ride is private
