package types

type NewBackupDTO struct {
	ID               uint
	Name             string
	Frequency        string
	Timezone         string
	CustomFrequency  string
	Type             string
	Bucket           string
	Region           string
	DayOfWeek        *int
	DayOfMonth       int
	Month            int
	Time             string
	DbName           string
	DbUser           string
	DbPassword       string
	DbHost           string
	DbPort           string
	PodLabel         string
	PodName          string
	Container        string
	FilesPath        string
	S3AccessKey      string
	S3SecretKey      string
	StorageDirectory string
	Retention        uint64
	Emails           string
}
