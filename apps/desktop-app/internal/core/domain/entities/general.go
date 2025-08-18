package entities

type AsyncOperationResult struct {
	Success    bool
	Message    string
	BackupData map[string]interface{}
	Error      error
}