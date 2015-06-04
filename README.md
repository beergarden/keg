# keg

Data uploader for beerserver.

## Usage

Compile and install:

```
go install
```

Upload temperature to a channel:

```
keg {channel-id}
```

Execute periodically with crontab:

```
*/5 * * * * /home/pi/go/bin/keg {channel-id}
```
