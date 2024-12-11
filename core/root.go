package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

type Command struct {
	tool string
	args []string
}

type Root struct {
	Stats    Stats  `json:"stats"`
	Tests    []Test `json:"tests"`
	Pending  []Test `json:"pending"`
	Failures []Test `json:"failures"`
	Passes   []Test `json:"passes"`
}

// Stats represents the stats structure.
type Stats struct {
	Suites   int       `json:"suites"`
	Tests    int       `json:"tests"`
	Passes   int       `json:"passes"`
	Pending  int       `json:"pending"`
	Failures int       `json:"failures"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration int       `json:"duration"`
}

// Test represents the structure of individual test entries.
type Test struct {
	Title        string                 `json:"title"`
	FullTitle    string                 `json:"fullTitle"`
	Duration     int                    `json:"duration"`
	CurrentRetry int                    `json:"currentRetry"`
	Err          map[string]interface{} `json:"err"` // Using a map to represent the generic "err" object.
}

func writeSmth(input string) {
	filePath := "output"

	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(input)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File written successfully.")
}

func parallelCmdExec(wg *sync.WaitGroup, outputChan chan []byte, tool string, script string, spec string, cyArgs string) {
	defer wg.Done()
	command, err := buildCommand(tool, script, spec, []string{cyArgs})

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
		log.Printf("Error running test %s: %v\n", spec, runErr)
	} else {
		log.Printf("Finished running command %s", cmd.String())
	}

	// ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	// re := regexp.MustCompile(`\x1b\[[0-9;]*[mGKHF]`)

	// log.Print("success:", re.ReplaceAllString(cmdOut.String(), ""))
	// log.Print("error:", cmdErr.String())

	outputChan <- append(cmdErr.Bytes(), cmdOut.Bytes()...)
}

func Run(tool string, dir string, script string, cyArgs string) error {

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
		go parallelCmdExec(&wg, outputChan, tool, script, spec, cyArgs)
		// go func() {
		// 	defer wg.Done()
		// 	command, err := buildCommand(tool, script, spec, []string{cyArgs})

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	cmd := exec.Command(command.tool, command.args...)

		// 	log.Printf("Running command %s ...", cmd.String())

		// 	var cmdOut bytes.Buffer

		// 	var cmdErr bytes.Buffer

		// 	cmd.Stdout = &cmdOut
		// 	cmd.Stderr = &cmdErr

		// 	runErr := cmd.Run()

		// 	if runErr != nil {
		// 		log.Printf("Error running test %s: %v\n", spec, runErr)
		// 	} else {
		// 		log.Printf("Finished running command %s", cmd.String())
		// 	}

		// 	// ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		// 	re := regexp.MustCompile(`\x1b\[[0-9;]*[mGKHF]`)

		// 	log.Print("success:", re.ReplaceAllString(cmdOut.String(), ""))
		// 	log.Print("error:", cmdErr.String())

		// 	outputChan <- append(cmdErr.Bytes(), cmdOut.Bytes()...)

		// }()
	}

	wg.Wait()
	close(outputChan)

	// var aggregatedOutput []string
	// var aggregatedOutput []Root
	// for output := range outputChan {

	// 	smth := make([]byte, 0, 100)

	// 	for _, val := range output {
	// 		if val != 10 {
	// 			smth = append(smth, val)
	// 		}
	// 	}
	// 	// content, _ := extractJSON(string(smth))
	// 	content, err := parseJSON(string(smth))

	// 	if err != nil {
	// 		log.Println(err)
	// 	} else {
	// 		aggregatedOutput = append(aggregatedOutput, *content)
	// 	}

	// 	// displayTable(content)
	// 	// aggregatedOutput = append(aggregatedOutput, content)

	// }

	// var finalOutput Root = Root{}

	// for _, val := range aggregatedOutput {
	// 	finalOutput.Stats.Passes += val.Stats.Passes
	// 	finalOutput.Stats.Suites += val.Stats.Suites
	// 	finalOutput.Stats.Failures += val.Stats.Failures
	// 	finalOutput.Tests = append(finalOutput.Tests, val.Tests...)
	// 	finalOutput.Passes = append(finalOutput.Passes, val.Passes...)
	// 	finalOutput.Failures = append(finalOutput.Failures, val.Failures...)
	// }

	// displayTable(&finalOutput)

	// fmt.Println("Aggregated Output:")
	// for _, out := range aggregatedOutput {
	// 	// fmt.Println(out)
	// 	// printJSON(out)
	// }
	return nil
}

func buildCommand(tool string, script string, specFile string, cyArgs []string) (*Command, error) {

	switch tool {
	case "yarn":
		//TODO add the spec file
		return &Command{
			tool: tool,
			args: append([]string{"run", script, "--spec", specFile}, cyArgs...),
		}, nil
	case "npx yarn":
		return &Command{
			tool: "npx",
			args: append([]string{"yarn", "run", script, "--spec", specFile}, cyArgs...),
		}, nil
	case "docker":
		dirs := strings.Split(filepath.Clean(specFile), "/")

		baseDir := dirs[0]

		rest := filepath.Join(dirs[1:]...)

		return &Command{
			tool: tool,
			args: []string{"run", "-i", "-v", fmt.Sprintf("./%s:/e2e", baseDir), "-w", "/e2e", "cypress/included:13.15.0", "-s", rest},
		}, nil
		//TODO test the npm and docker -t flag
	case "npm":
		return &Command{
			tool: tool,
			args: append([]string{"run", script, "--"}, cyArgs...),
		}, nil
	}

	return nil, errors.New("the tool must be docker, yarn or npx yarn, npm, or npx")

}

func parseJSON(content string) (*Root, error) {
	re := regexp.MustCompile(`\{.*\}`)
	match := re.FindString(content)

	if match == "" {
		return nil, errors.New("the regex {.*} did not find anything from the cypress response")
	}
	var result Root
	err := json.Unmarshal([]byte(match), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func displayTable(result *Root) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	// fmt.Printf("%+v", result)
	fmt.Fprintln(writer, "Title\tFullTitle\tDuration(ms)\tCurrentRetry")

	// Print test data
	for _, test := range result.Tests {
		fmt.Fprintf(writer, "%s\t%s\t%d\t%d\n", test.Title, test.FullTitle, test.Duration, test.CurrentRetry)
	}

	// Flush the writer to output
	writer.Flush()
}

func extractJSON(content string) (string, bool) {

	// Compile the regex
	re := regexp.MustCompile(`\{.*\}`)
	match := re.FindString(content)

	if match != "" {
		fmt.Println("Extracted JSON:")
		fmt.Println(match)
		return match, true
	} else {
		fmt.Println("No JSON found.")
		return "", false
	}
}
