package gosync

// Action defines what needs to be done to the file and s3 object.
type Action int

// Allowed action as constants.
const (
	ActionNone Action = iota
	ActionDeleteObject
	ActionUploadFile
	ActionDownloadObject
	ActionDeleteFile
)

func (a Action) String() string {

	switch a {
	case ActionNone:
		return "NO ACTION NEEDED"
	case ActionDeleteObject:
		return "DELETE OBJECT"
	case ActionUploadFile:
		return "UPLOAD FILE"
	case ActionDeleteFile:
		return "DELETE FILE"
	case ActionDownloadObject:
		return "DOWNLOAD OBJECT"
	default:
		return "UNKNOWN ACTION"
	}

}

// SyncDirection indicates in what direction we are synchronizing.
type SyncDirection int

// Default is neither backup nor restore
const (
	DirectionNone    SyncDirection = iota // No effect on either side
	DirectionBackup                       // File => S3
	DirectionRestore                      // S3 => File
)
