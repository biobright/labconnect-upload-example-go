package main

import (
	"fmt"
	"github.com/google/uuid"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	testFile, err := os.Open("./testformat2file.csv")
	check(err)
	defer testFile.Close()
	hash, err := calculate(testFile, 1024*1024*10) // always 10MB for LLC
	check(err)
	//println(hash)
	agentId, err := getLLCAgentId("upload_test_scott")
	// generating a uuid is just for example - we can decide to use whatever fot the folder name
	folderName := uuid.Must(uuid.NewV7())
	err = createVirtualFolder(fmt.Sprintf("/%s", folderName.String()), agentId)
	check(err)
	testFileRead2, err := os.Open("./testformat2file.csv")
	check(err)
	defer testFileRead2.Close()
	err = uploadFileToVirtualFolder(hash, agentId, fmt.Sprintf("/%s/testformat2file.csv", folderName.String()), "testformat2file.csv", testFileRead2)
	check(err)
}
