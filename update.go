package sync

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CompareState byte

const (
	Same CompareState = iota
	Local
	Server
	OutOfSync
	NotExistLocal
	NotExistServer
	Unknown
)

func (c *SyncClient) Update(files ...string) error {
	for _, path := range files {
		// TODO: set up handling for distinct local and server paths
		state := c.compareFileState(path, path)
		switch state {
		case Same:
			log.Println("Files are the same.")
			return nil
		case Local, NotExistServer:
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			f, err := c.SftpClient.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.Write(data); err != nil {
				return err
			}
			log.Printf("Wrote %s to %s (wrote to server)\n", truncateString(string(data)), path)
		case Server:
			f, err := c.SftpClient.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, f)
			if err != nil {
				log.Fatalf("failed to read file: %v", err)
			}

			if err := os.WriteFile(path, buffer.Bytes(), os.ModePerm); err != nil {
				return err
			}
			log.Printf("Wrote %s to %s (wrote to client)\n", truncateString(buffer.String()), path)
		case NotExistLocal:
			// TODO: delete from local (or not? -- multiple local instances -- keep log on server?)
			// NOTE: probably use separate signal from watcher to handle deletions
			return fmt.Errorf("Not yet implemented")
		case OutOfSync:
			// TODO: probably give up :(
			return fmt.Errorf("Local and server states are out of sync. (automatic conflict resolution not yet implemented)")
		case Unknown:
			return fmt.Errorf("Compare state is Unknown. How did you get here?")
		}
	}

	return nil
}

func truncateString(input string) string {
	const maxLength = 15
	newlineIndex := strings.Index(input, "\n")

	if newlineIndex != -1 && newlineIndex < maxLength {
		return input[:newlineIndex]
	}

	if len(input) > maxLength {
		lines := strings.Count(input, "\n") + 1
		return fmt.Sprintf("%s... (%d more lines)", input[:maxLength], lines)
	}

	return input
}
