Navarchiver archives your [Navidrome](https://www.navidrome.org/) audio library and metadata using GCS (Google Cloud Storage).

It runs in 3 modes:

* [Scheduled](#scheduled)
* [Ledger](#ledger)
* [Batch](#batch)

Tested on versions:

* [v0.54.5](https://github.com/navidrome/navidrome/releases/tag/v0.54.5)
* [v0.58.0](https://github.com/navidrome/navidrome/releases/tag/v0.58.0)

## Scheduled
Scheduled mode is designed to be used for archiving the library "moving forward" on a nightly basis. For backfilling, use [batch mode](#batch).

Scheduled mode should run in same directory as where you keep your navidrome DB.

### Environment variables

* GCS_PROJECT_ID - Project ID for GCS uploading
* GCS_BUCKET_NAME - Bucket name for GCS uploading
* GCS_TIMEOUT_SECONDS - maximum time to allow for a GCS upload before timing out
* GCS_CREDS_FILE - path to GCS JSON credentials file

For alerting, [discord-alert](https://github.com/apkatsikas/discord-alert) is used. Please see the documentation for this tool and for more info on [creating a bot](https://github.com/apkatsikas/discord-alert?tab=readme-ov-file#creating-a-bot).

* BOT_TOKEN - Discord bot token for alerting on failure of the nightly backup
* CHANNEL_ID - Discord channel ID for alerting

### Running scheduled mode

One nice way to run scheduled mode is a Linux cron job.

An example cron command might look like:

`cd /var/lib/navidrome && /opt/navarchiver/navarchiver >> /var/lib/navidrome/archiver.log 2>&1`

## Batch

Batch mode takes a JSON file produced by the [ledger mode](#ledger). You might choose to split the resultant ledger into smaller JSON files and run multiple batches. This mode is suitable for backfilling an archive for an existing library. It only backfills the audio library, not the SQLite metadata.

You will need to set:

* GCS_PROJECT_ID
* GCS_BUCKET_NAME
* GCS_TIMEOUT_SECONDS
* GCS_CREDS_FILE

from the [environment variables](#environment-variables) section.

This mode is invoked using the `-runMode=batch` flag, and takes a single positional argument for the ledger path.

## Ledger

The ledger mode produces a JSON file suitable for backfilling a batch in [batch mode](#batch) of the user's Navidrome library.

This mode is invoked using the `-runMode=ledger` flag, and takes 2 positional arguments for:

* Location of the Navidrome SQLite DB file
* Destination for the ledger JSON file
