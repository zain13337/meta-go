package socket

/*
	type 3 = task
*/

var TaskLabels = map[string]string{
	"ThreadMarkRead":          "21",
	"AddParticipantsTask":     "23",
	"UpdateAdminTask":         "25",
	"SendReactionTask":        "29",
	"SearchUserTask":          "30",
	"SearchUserSecondaryTask": "31",
	"RenameThreadTask":        "32",
	"DeleteMessageTask":       "33",
	"SetThreadImageTask":      "37",
	"SendMessageTask":         "46",
	"ReportAppStateTask":      "123",
	"CreateGroupTask":         "130",
	"RemoveParticipantTask":   "140",
	"MuteThreadTask":          "144",
	"FetchThreadsTask":        "145",
	"DeleteThreadTask":        "146",
	"DeleteMessageMeOnlyTask": "155",
	"CreatePollTask":          "163",
	"UpdatePollTask":          "164",
	"GetContactsFullTask":     "207",
	"CreateThreadTask":        "209",
	"FetchMessagesTask":       "228",
	"GetContactsTask":         "452",
	"EditMessageTask":         "742",
}

type Task interface {
	GetLabel() string
	Create() (interface{}, interface{}, bool) // payload, queue_name, marshal_queuename
}

type TaskData struct {
	FailureCount interface{} `json:"failure_count"`
	Label        string      `json:"label,omitempty"`
	Payload      interface{} `json:"payload,omitempty"`
	QueueName    interface{} `json:"queue_name,omitempty"`
	TaskId       int64       `json:"task_id"`
}
