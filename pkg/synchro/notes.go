// Package synchro provides synchronization services between part of a file system and a s3 bucket.
// It is not called sync, beacause sync is actually used !
//
// Here are some thoughts/notes/todos about the sync process :
// * remember S3 keys are utf8 and must be below 1024 bytes in length. That should be enforced by design.
// * It seems double shashs // will be processed in a specific way, double check that.
// * check for max object size limits ? use Multipart upload above 100Mb, enforce max object size of 5 TB (not really useful - enforce a lower limit to keep the process practical !) ?
// * how do we confirm the upload (potentially multipart) was sucessful ?
//
// Syncing is based upon :
// * potential actions are : detroy a (obsolete) remote file or upload a (new or updated) local file
// * one should be able to copy bucket cantent verbatim to restore
//   the local file system : this implies, no file name transformation, no meta tags, ...
// * when was the local file last updated, when was the s3 file written
// * use UTC for time comparison
// * verify if size comparison can be an effcient copmpare ?
// * do we need to maintain a checksum/md5 table somewhere of all the files modified since lost sync ?
// * how do we handle links (soft/hard) ?
//
// Efficiency concerns
// * allow for parallel uploading
// * single, central manafgement of what needs to be destroyed/uploaded
package synchro
