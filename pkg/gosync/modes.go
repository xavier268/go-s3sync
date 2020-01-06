package gosync

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

	ModeCleanEmptyDirs
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
	case ModeCleanEmptyDirs:
		return "Cleaning empty dirs"
	default:
		panic(m)
	}
}
