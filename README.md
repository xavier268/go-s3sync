# go-s3sync

## What is it ?

This package provides an efficient (concurrent processing) and simple synchronisation tool between a directory and an s3 bucket.
* parallel processing of files and S3 buckets
* backup or restore modes, possibly in a mock (do-nothing) format

Command line tools are provided, using the package :
* backup (or backupmock, to simulate a backup)
* restore (or restoremock, to simulate a restore)

Bucket name and directory are set with cli options. Use the -h flag more more details.

AWS authentication is done via credentials files or IAM setting. 
There are obviously no secret in the code !

## Design principles and notes :

The S3 bucket content can be manually restored/edited/examined if needed. No meta data are being added (beyond name, size, lastupdated, all automatically managed by S3).

Empty directories are ignored. A cleaup utility is provided to remove them locally - it is voluntary not done automatically while synchronising.

UTC is use as the sole time reference.

Soft links are copied as is as files, hard links are copied as separate objects.

Upload/Download use the s3manager version of the API, allowing for up to 5 TB (!!) per file/object.

The max object key length (see AWS documentation) is enforced at 1000 bytes. A longer file name or key will panic and stop processing.

Synchronizations decisions are based solely upon file or s3 object  name, size, and last updated time. ETAGS are not used.

Except for the mock versions, restore and backup may and **will overwite or delete existing information**, if needed. 
**USE WITH CARE** on real world data !

Local files are never accessed locally outside of the file system "prefix" set at configuration time.
However, the entire S3 buscket specified will be accessed, and possibly modified upon backup.

Special attention was given to the concurrency design to maximize the throughput while taking into account that S3 does not provide any transactionnal support. For instance, I decided not to let the fileprocessing and the s3 processing run in parallel ...

Public API surface was reduced to the minimum.
