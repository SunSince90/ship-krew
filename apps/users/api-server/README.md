# Ship Krew Users

This is the readme for the users section of the application.

It is supposed to be the core and most important part of it.

TODO: better description.

## TODO

### Flags

- [ ] `--database-url`: for where to find the database. (default: users.ship-krew-databases)
- [ ] `--cache-url`: for where to find the cache. (default: users.ship-krew-caches)

### Others

- [ ] Prevent users from registering as `healthz`, as this is used for liveness probe
- [ ] Return appropriate status codes on probes, i.e. is not ready if it does not detect a database, cache etc...
- [ ] Rename `per-page` as `perPage`? And same for other query parameteres
