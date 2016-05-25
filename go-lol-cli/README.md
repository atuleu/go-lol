# go-lol-cli

`go-lol-cli` is a command line utility to save and display League Of
Legends Replays


## Installation

You must first need to install
[go](https://golang.org/doc/install). Please make sure to
[set up your environment](https://golang.org/doc/code.html)
accordingly in order for `go get` to work.

Then simply run :

```bash
go get -u github.com/charlywazzap/go-lol/go-lol-cli
```

You should be now able to run this :

```bash
go-lol-cli -h
```

## Setup

In order to work and fetch data from Riot Games server, one need an
API Key. Since this project is in alpha, there is no global, high
bandwith key shared on a server, but user are asked to use their own
key. This key is only used to get data from the server, and cannot
alter any data on Riot Games server (so the utility won't sniff all
your RP points from the game).


Go to https://developer.riotgames.com/ and login with a valid LoL
account to generate an api-key. Then you can use the following command
to finish the set-up of `go-lol-cli` (it will just stire the magic
number somewhere safe, where only you can see it).

```bash
go-lol-cli set-api-key <key-from-riot-games>
```

## Manual

### Record replay

To record in loop any current game of a player, simply use :
```bash
go-lol-cli -r <region> watch-summoner <SummonerName>
```

`-r` option can be used to select a region like `na`, `kr` or
`euw`. Default is `euw` because EUW > NA.

### List recorded replay

```bash
go-lol-cli [-r <region>] list-replays
````

Will display on stdout all recorded replays for the specified region

### Watch a replay

```bash
go-lol-cli [-r <region] replay [GameID]
```

If you do not pass something to this command, the newest replay will
be recorded. Otherwise you need to specify the long GameID (10 digits
at the moment) to the command

### Clean Up

```bash
go-lol-cli [-r <region>] garbage-collect --older-than 168h --limit 10
```

This will keep no more than 10 replays (older are deleted first), and
all replays that are older than 168 hours (aka 1 week).


## License

This software is Licensed under GPL version 3. A copy of the license
can be found in the file ../LICENSE

## Disclaimer

go-lol and go-lol-cli isn’t endorsed by Riot Games and doesn’t reflect
the views or opinions of Riot Games or anyone officially involved in
producing or managing League of Legends. League of Legends and Riot
Games are trademarks or registered trademarks of Riot Games,
Inc. League of Legends (c) Riot Games, Inc.

