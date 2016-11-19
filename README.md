# Gobak - tool for rotating database backup archives.

Gobak is a simple tool written in go for rotating database archive files. 
You can rotate multipl slots with different backup counts. 

## Get started

1. Clone the repo and run script or download binaries and run.
#### Run the script
```
./run
```

#### Run the binary
```
./gobak.bin
```

2. Copy `config-sample.json` to `config.json`.
3. Run the app periodically such as using cron.

## Sample config file

```
{
    // Enable/Disable logging
    "log": true,
    // Global backup path
    "backup_path": "/home/db/backups",
    // Backup file format. 
    // Eg: if you backup file name is like `backup-2016-10-10.psql.gz`
    // backup file format is `backup-%s.psql.gz` and
    // backup date format is 'YYYY-MM-DD'
    // backup date format is based on excel date formats - https://goo.gl/oCTDIs
    "backup_file_format": "backup-%s.psql.gz",
    "backup_file_date_format": "YYYY-MM-DD",

    // Define your backup slots here in array
    "backup_slots": [
        {
            // slot name
            "name": "weekly",
            // slot relative path. Lets say if your absolute path is `/home/db/backups/weekly`
            // then slot path name is `weekly` 
            "path": "weekly",
            // Number of backup files to keep. If bacup file count exceeds this count then its deleted.
            "count": 4
        },
        {
            "name": "daily",
            "path": "daily",
            "count": 7
        }
    ]
}
```

Here is a sample structure of the backup folder

```
└── /home/db/backups
    ├── daily
    │   ├── backup-2016-10-05-psql.gz
    │   ├── backup-2016-10-06-psql.gz
    │   ├── backup-2016-10-07-psql.gz
    │   ├── backup-2016-10-08-psql.gz
    │   ├── backup-2016-10-09-psql.gz
    │   ├── backup-2016-10-10-psql.gz
    │   └── backup-2016-10-11-psql.gz
    └── weekly
        ├── backup-2016-10-15-psql.gz
        ├── backup-2016-10-22-psql.gz
        ├── backup-2016-10-29-psql.gz
        └── backup-2016-11-14-psql.gz
```

## What it doesn't do

1. It doesn't run as a dameon. Its needs to be run periodically based on your backup frequencies.]
2. It doesn't support rotation based on period. It just deletes files more than the backup count set.
For example it doesn't move this weeks last day daily backup to monthly backup and rotate likewise. 
You can expect this feature in future versions.
3. It doesn't take backup of your database. It just rotates the backup files which are periodically dumped to a folder.
For example I setup daily, weekly crons of `pg_dump` crons to dump database in a folder and run this program in cron.