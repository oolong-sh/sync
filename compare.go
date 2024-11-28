package sync

import (
	"log"
	"os"
)

func (c *SyncClient) compareFileState(localPath string, serverPath string) CompareState {
	localInfo, err := os.Stat(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return NotExistLocal
		}

		log.Println("Failed to get local file info: ", err)
		return Unknown
	}

	serverInfo, err := c.SftpClient.Stat(serverPath)
	if err != nil {
		if os.IsNotExist(err) {
			return NotExistServer
		}

		log.Println("Failed to get server file info: ", err)
		return Unknown
	}

	// compare by mod times
	// var newer CompareState
	if localInfo.ModTime().Unix() > serverInfo.ModTime().Unix() {
		// newer = Local
		return Local
	} else if localInfo.ModTime().Unix() < serverInfo.ModTime().Unix() {
		// newer = Server
		return Server
	}

	// TODO: compare file contents (for conflict resolution)

	return Same
}
