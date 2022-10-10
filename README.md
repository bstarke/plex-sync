# plex-sync

## Todo:

- Convert to Cobra CLI
- Main Functions:
  - daemon : run as background job (requires cron schedule in config)
  - sync : one time sync (like daemon but doesn't sleep, quits after one run)
  - csv : create csv of library
  - print : print options to screen 
    - options: 
      - movies (-t, --type movies)
      - shows (-t, --type shows)
      - only resolution not = config (-q, --quality)
      - only format not equal to config (-f, --format)

## Sync Plan (daemon mode)

- Load DB for remote
- Load DB for local
- Fetch missing movies
- Fetch movies with better resolution
- Fetch movies with preferred format
- Fetch missing shows
- Fetch shows with better resolution (if enabled)

## Options

- Preferred Format
- Preferred Resolution
- If new shows/movies are of higher resolution than settings
  - Download (true/false)
  - If downloaded : email, write log, ignore?