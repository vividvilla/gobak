package main

type Config struct {
	Log                  bool         `json:"log"`
	BackupPath           string       `json:"backup_path"`
	BackupFileFormat     string       `json:"backup_file_format"`
	BackupFileDateFormat string       `json:"backup_file_date_format"`
	BackupSlots          []BackupSlot `json:"backup_slots"`
}

type BackupSlot struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Count int    `json:"count"`
}
