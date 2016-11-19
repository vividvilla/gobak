// Gobak - Simple tool for database backup archive rotation.
// Vivek R, http://vivekr.net
// Nov 2016
// MIT License

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"

	"path"

	"io/ioutil"

	"github.com/knadh/jsonconfig"
	"github.com/metakeule/fmtdate"
)

var (
	cfg                   *Config
	logger                *log.Logger
	backupFileFormatRegex *regexp.Regexp
	slotFrequencyRegex    *regexp.Regexp
)

// Sort implementation for slices of time
type timeSlice []time.Time

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	return p[i].Before(p[j])
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	logger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)

	// Enable of disable logger
	// Load json config from `config.json``
	err := jsonconfig.Load("config.json", &cfg)
	if err != nil {
		logger.Println("Unable to load config:", err)
	}

	if !cfg.Log {
		logger = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime)
	}

	// Regex for matching date string in backup file names.
	backupFileFormatRegex = regexp.MustCompile(fmt.Sprintf(cfg.BackupFileFormat, "(.+?)"))

	for _, slot := range cfg.BackupSlots {
		slotBackupPath := path.Join(cfg.BackupPath, slot.Path)
		logger.Println(fmt.Sprintf("Rotating slot '%s' in path '%s'", slot.Name, slotBackupPath))
		rotateSlot(slot, slotBackupPath, cfg.BackupFileFormat, cfg.BackupFileDateFormat)
	}

	logger.Println("Voila..it's done.")
}

// Rotate files in a slot
func rotateSlot(slot BackupSlot, slotBackupPath string, backupFileFormat string, backupDateFormat string) {
	// Read all files in slot directory
	files, err := ioutil.ReadDir(slotBackupPath)

	if err != nil {
		logger.Println(fmt.Sprintf("Couldn't slot backup directory '%s', moving on.", slotBackupPath), err)
		return
	}

	if len(files) == 0 {
		logger.Println("No files found, moving on.")
		return
	}

	// Get sorted array of dates representing the backup files in the slot
	backupDates := getBackupDates(files, backupDateFormat)

	// Remove old backup files if it exceeds backup count.
	if len(backupDates) > slot.Count {
		deleteFiles(backupDates[slot.Count:len(backupDates)], slotBackupPath, backupFileFormat, backupDateFormat)
	} else {
		logger.Println("Nothing to rotate..cheers.")
	}
}

// Return sorted array of dates representing the backup files in the slot
func getBackupDates(backupFiles []os.FileInfo, backupDateFormat string) timeSlice {
	var backupDates timeSlice
	for _, f := range backupFiles {
		match := backupFileFormatRegex.FindStringSubmatch(f.Name())

		if len(match) != 2 {
			continue
		}

		fileDate, err := fmtdate.Parse(backupDateFormat, match[1])

		if err != nil {
			continue
		}

		backupDates = append(backupDates, fileDate)
	}

	// sort all the files matching the file format
	sort.Sort(sort.Reverse(backupDates))
	return backupDates
}

// Delete files based on given array of time representing the backup file
func deleteFiles(backupDates timeSlice, slotBackupPath string, backupFileFormat string, backupDateFormat string) {
	for _, date := range backupDates {
		formattedDate := fmtdate.Format(backupDateFormat, date)
		formattedFilename := fmt.Sprintf(backupFileFormat, formattedDate)
		filePath := path.Join(slotBackupPath, formattedFilename)

		var err = os.Remove(filePath)
		if err != nil {
			logger.Println(fmt.Sprintf("Error while deleteing the file '%s'", filePath), err)
		} else {
			logger.Println(fmt.Sprintf("Delete a file '%s'", filePath))
		}
	}
}
