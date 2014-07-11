/**
 * WatchFfile.go - a simplistic file watcher to hold the fort until Go 1.4 when
 * fsnotify is expected to become core.
 */

package wsnotify

func WatchFile(filePath string, checks_every_n_seconds int) error {
	// Record the initial stat of the file.
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Sleep before we start polling
	time.Sleep(check_every_n_seconds * time.Second)

	// polling forever waiting for Guffman
	for {
		// Check the current stat of the file.
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		// Check if mod time or size changes
		if stat.ModTime() != initialStat.ModTime() {
			break
		}

		// Sleep for a time
		time.Sleep(check_every_n_seconds * time.Second)
	}
	return nil
}

/* USAGE example:
doneChan := make(chan bool)

go func(doneChan chan bool) {
    defer func() {
        doneChan <- true
    }()

    err := watchFile("/path/to/file")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("File has been changed")
}(doneChan)

<-doneChan
*/
