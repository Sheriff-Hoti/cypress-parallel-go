package core

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Command struct {
	tool string
	args []string
}

func Run(tool string, dir string, cyArgs string) error {

	//getting the folders which match the 'dir' pattern
	matches, mError := filepath.Glob(dir)

	if mError != nil {
		return mError
	}

	log.Println("Files found:")
	for _, value := range matches {
		log.Println(value)
	}

	outputChan := make(chan []byte, 10)

	var wg sync.WaitGroup

	wg.Add(len(matches))

	for _, spec := range matches {
		go func() {
			defer wg.Done()
			command, err := buildCommand(tool, spec)

			if err != nil {
				log.Fatal(err)
			}

			cmd := exec.Command(command.tool, command.args...)

			log.Printf("Running command %s ...", cmd.String())

			var cmdOut bytes.Buffer

			var cmdErr bytes.Buffer

			cmd.Stdout = &cmdOut
			cmd.Stderr = &cmdErr

			runErr := cmd.Run()

			if runErr != nil {
				log.Printf("Error running test %s: %v\n", spec, err)
			} else {
				log.Printf("Finished running command %s", cmd.String())
			}

			outputChan <- append(cmdErr.Bytes(), cmdOut.Bytes()...)

		}()
	}

	wg.Wait()
	close(outputChan)

	var aggregatedOutput []string
	for output := range outputChan {

		smth := make([]byte, 0, 100)

		for _, val := range output {
			if val != 10 {
				smth = append(smth, val)
			}
		}
		aggregatedOutput = append(aggregatedOutput, string(smth))
	}

	fmt.Println("Aggregated Output:")
	for _, out := range aggregatedOutput {
		printJSON(out)
	}
	return nil
}

func buildCommand(tool string, specFile string) (*Command, error) {

	switch tool {
	case "yarn":
		return nil, errors.New("yarn not implemented")
	case "docker":
		dirs := strings.Split(filepath.Clean(specFile), "/")

		baseDir := dirs[0]

		rest := filepath.Join(dirs[1:]...)

		return &Command{
			tool: tool,
			args: []string{"run", "-i", "-v", fmt.Sprintf("./%s:/e2e", baseDir), "-w", "/e2e", "cypress/included:13.15.0", "-s", rest},
		}, nil
	}
	return nil, errors.New("the tool must be docker or yarn")

}

func printJSON(content string) {

	// Compile the regex
	re := regexp.MustCompile(`\{.*\}`)
	match := re.FindString(content)

	if match != "" {
		fmt.Println("Extracted JSON:")
		fmt.Println(match)
	} else {
		fmt.Println("No JSON found.")
	}
}
