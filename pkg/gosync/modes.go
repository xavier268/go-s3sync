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

// Mode indicates in what direction we are synchronizing,
// and if modifications are actually made.
type Mode int

// Dining the Mode constants.
// xxxMock means no operation is actually performed.
const (
	ModeBackupMock Mode = iota // File => S3
	ModeBackup

	ModeRestoreMock // S3 => File
	ModeRestore
)

func (m *Mode) String() string {
	switch *m {
	case ModeBackup:
		return "Backup : File --> S3"
	case ModeBackupMock:
		return "Backup (mock) : File --> S3"
	case ModeRestore:
		return "Restore : S3 --> File"
	case ModeRestoreMock:
		return "Restore (mock): S3 --> File"

	default:
		panic(m)
	}
}
