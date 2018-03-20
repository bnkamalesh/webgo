# Release v2.0.0

1. Log levels
	1. Error logs are now printed to `os.Stderr` with a prefix `Error`
	2. Info logs are now printed to `os.Stdout` with a prefix `Info`
	3. Warning logs are now printed to `os.Stdour` with a prefix `Warning`

2. Removed per route access log control
3. Renamed option to toggle access log at router, to `AccessLog`.
	1. Access log is false/off by default
4. `Globals` in route is now renamed to `Globals` instead of `G`.