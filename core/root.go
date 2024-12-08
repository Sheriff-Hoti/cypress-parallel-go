package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type Command struct {
	tool string
	args []string
}

// func NewCypressPath(fullPath string) (*CypressPath, error) {

// }

func GetMatches(dir string) []string {
	matches, error := filepath.Glob(dir)

	if error != nil {
		log.Panic(error)
	}

	return matches
}

func writeSmth(input string) {
	filePath := "output.txt"

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

func Run(tool string, dir string) {

	matches := GetMatches(dir)

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
				log.Panic(err)
			}

			cmd := exec.Command(command.tool, command.args...)

			log.Printf("Running command %s ...", cmd.String())
			// cmd.Stdout = os.Stdout
			// cmd.Stderr = os.Stderr
			bytes, error := cmd.CombinedOutput()

			if error != nil {
				log.Printf("Error running test %s: %v\n", spec, error)
				// outputChan <- fmt.Sprintf("Error: %v", err)
			} else {
				// log.Println(string(bytes))
				log.Printf("Finnished running command %s ...", cmd.String())
				outputChan <- bytes
			}

			// if err := cmd.Run(); err != nil {
			// 	fmt.Printf("Error running test %s: %v\n", spec, err)
			// }
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

	// defer fmt.Println("All Cypress tests completed.")

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for i := 0; i < 100; i++ {
	// 		fmt.Println("IDX from first func:", i)
	// 		time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
	// 	}
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for i := 0; i < 100; i++ {
	// 		fmt.Println("IDX from second func:", i)
	// 		time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
	// 	}
	// }()

	// wg.Wait()
	// fmt.Println("Done")

	//TODO accept the cypress command
	// log.Printf("script %v \n", script)
	// // fmt.Printf("process %d", runtime.GOMAXPROCS(0))
	// // cmds := strings.Fields(script)
	// command := exec.Command("docker", "run", "-i", "-v", "/home/sherif/personal/cypress-parallel-go/test:/e2e", "-w", "/e2e", "cypress/included:13.15.0", "-s", "cypress/e2e/GoogleImage.cy.js")

	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr

	// log.Println(command.Args)

	// // output, err := command.CombinedOutput()
	// if err := command.Run(); err != nil {
	// 	fmt.Printf("Command failed: %v\n", err)
	// 	return
	// }

	// Print the output
	// log.Println(string(output))
	// matches, error := filepath.Glob("cypress/ParallelTest/*/*.cy.ts")
	//TODO replace the hardcoded path with the dir argument
	// cores := runtime.GOMAXPROCS(2)
	// start := time.Now()
	// MainStuff()
	// // ParallelExec()
	// elapsed := time.Since(start)
	// log.Printf("The time it took %s", elapsed)
	// log.Printf("Cores %v", cores)

}

func buildCommand(tool string, specFile string) (*Command, error) {

	if tool == "yarn" {
		log.Panic("yarn not implemented")
	}

	if tool == "docker" {

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

func runCypress() {

	// cmd := exec.Command("docker", "run", "-i", "-v", "/home/sherif/personal/cypress-parallel-go/test:/e2e", "-w", "/e2e", "cypress/included:13.15.0", "-s", testSpec)
	cmd := exec.Command("docker", "run", "-i", "-v", "/home/sherif/personal/cypress-parallel-go/test:/e2e", "-w", "/e2e", "cypress/included:13.15.0", "-s", "cypress/e2e/*.cy.js")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running test %s: %v\n", err, "2")
	}
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

func MainStuff() {
	// testSpecs := []string{"cypress/e2e/GoogleImage.cy.js", "cypress/e2e/GoogleSearch.cy.js", "cypress/e2e/GoogleSearch1.cy.js"}

	// for _, spec := range testSpecs {
	// 	go runCypress(spec, &wg)
	// }
	runCypress()

	log.Print(runtime.GOMAXPROCS(0))
	fmt.Println("All Cypress tests completed.")
}

func runCypressParallel(testSpec string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("docker", "run", "-i", "-v", "./test:/e2e", "-w", "/e2e", "cypress/included:13.15.0", "-s", testSpec)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running test %s: %v\n", testSpec, err)
	}
}

func MainStuffParallel() {
	testSpecs := []string{"cypress/e2e/GoogleImage.cy.js", "cypress/e2e/GoogleSearch.cy.js", "cypress/e2e/GoogleSearch1.cy.js", "cypress/e2e/FirstTest.cy.js"}

	var wg sync.WaitGroup
	wg.Add(len(testSpecs))

	for _, spec := range testSpecs {
		go runCypressParallel(spec, &wg)
	}

	defer wg.Wait()
	fmt.Println("All Cypress tests completed.")
}

func ParallelExec() {
	testSpecs := []string{"cypress/e2e/GoogleImage.cy.js", "cypress/e2e/GoogleSearch.cy.js", "cypress/e2e/GoogleSearch1.cy.js", "cypress/e2e/FirstTest.cy.js"}
	var commands []*exec.Cmd

	for _, spec := range testSpecs {
		cmd := exec.Command("docker", "run", "-i", "-v", "./test:/e2e", "-w", "/e2e", "cypress/included:13.15.0", "-s", spec)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting test %s: %v\n", spec, err)
			continue
		}
		commands = append(commands, cmd)
	}

	for _, cmd := range commands {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("Error waiting for command: %v\n", err)
		}
	}

	fmt.Println("All Cypress tests completed.")
}
