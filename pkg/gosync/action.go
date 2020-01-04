package gosync

// Action defines what needs to be done to the file and s3 object.
type Action int

// Allowed action as constants.
const (
	ActionNone Action = iota
	ActionDeleteObject
	ActionUploadFile
)

func (a Action) String() string {

	switch a {
	case ActionNone:
		return "NO ACTION NEEDED"
	case ActionDeleteObject:
		return "DELETE OBJECT"
	case ActionUploadFile:
		return "UPLOAD FILE"
	default:
		return "UNKNOWN ACTION"
	}

}
