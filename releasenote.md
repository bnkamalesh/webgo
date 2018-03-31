# Release v2.0.0

1. Log levels
	1. Error logs are now printed to `os.Stderr` with a prefix `Error`
	2. Info logs are now printed to `os.Stdout` with a prefix `Info`
	3. Warning logs are now printed to `os.Stdour` with a prefix `Warning`

2. Removed per route access log control
3. Renamed option to toggle access log at router, to `AccessLog`
	1. Access log is false/off by default
4. `Globals` is now removed, and router holds the app context now
5. Removed templates, as it can just be added to Globals' app context 
(`App` which is a map[string]interface{})
6. Removed configuration `HTTPSOnly` as it can be started by calling `StartHTTPS`
7. Read and write timeout are now added in configuration instead of passing to Start