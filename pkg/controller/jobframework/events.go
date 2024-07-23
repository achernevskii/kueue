package jobframework

// JobReconciler event action list
const (
	Started            = "Started"
	Suspended          = "Suspended"
	Stopped            = "Stopped"
	CreatedWorkload    = "CreatedWorkload"
	FinalizingWorkload = "FinalizingWorkload"
	DeletedWorkload    = "DeletedWorkload"
	UpdatedWorkload    = "UpdatedWorkload"
	WorkloadCompose    = "WorkloadCompose"
)

// JobReconciler event reason list
const (
	FinishedWorkload   = "FinishedWorkload"
	ErrWorkloadCompose = "ErrWorkloadCompose"
)
