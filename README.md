# Project YogaTube

Simple site that use youtube playlist and change to next video everyy day. Used for my daily morning yoga routine.
Demo: https://yoga.fumlig.com/

## Configuration

Set `YOUTUBE_API_KEY` in the server environment or in the local `.env` file before importing a playlist:

```dotenv
YOUTUBE_API_KEY=your-youtube-api-key
```

The `.env` file is ignored by Git. In production, set this variable in your deployment configuration and restart the application.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
