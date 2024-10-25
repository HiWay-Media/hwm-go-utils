package file_utils

import (
	"os"
    "log"
)

func WriteToFile(s []byte, file string) {
	//
	f, err := os.Create(file)
	if err != nil {
		//slack.Bug("WriteToFile", err.Error())
		log.Fatal(err)
	}
	f.Write(s)
	defer f.Close()
}
