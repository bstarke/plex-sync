# plex-sync

## Todo:

- ~~Convert to Cobra CLI~~
- Main Functions:
  - ~~csv : create csv of library~~
  - sync : one time sync (like daemon but doesn't sleep, quits after one run)
  - daemon : run as background job (requires cron schedule in config)
  - print : print options to screen
    - options:
      - movies (-t, --type movies)
      - shows (-s, --type shows)
      - only resolution not equal to config (-q, --quality)
      - only format not equal to config (-f, --format)

## Sync Plan

- Load DB for local
- Load DB for remote
- Fetch missing movies
- Push missing movies
- Fetch movies with better resolution
- Push movies with better resolution
- Fetch movies with preferred format
- Push movies with preferred format
- Fetch missing shows
- Push missing shows
- Fetch shows with better resolution
- Push shows with better resolution

## Options

- Preferred Format
- Preferred Resolution
- If new shows/movies are of higher resolution than settings
  - Download (true/false)
  - If downloaded : email, write log, ignore?