package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Credentials struct {
	BaseURL  string `json:"_BASEURL"`
	Username string `json:"_USERNAME"`
	Password string `json:"_PASSWORD"`
}

const credentialsFileName = "nexus-credentials.json"

const containerName = "conan2-pollinator-container"

func readCredentials() Credentials {
	// Read the JSON file
	file, err := os.Open(credentialsFileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Unmarshal the JSON content into a Credentials struct
	var creds Credentials
	if err := json.Unmarshal(byteValue, &creds); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}
	return creds
}

func buildConanImage() {
	// Define the docker build command
	cmd := exec.Command("docker", "build", "-t", "conan2-pollinator", ".")

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
		return
	}

	// Print the output
	log.Printf("Output: %s\n", output)
}

func configureConanContainer(credentials Credentials) {
	var configureCommands []string = []string{
		"conan remote remove conancenter",
		fmt.Sprintf("conan remote add conan-hosted-v2 %s/repository/conan-hosted2/ --insecure", credentials.BaseURL),
		fmt.Sprintf("conan remote login conan-hosted-v2 %s -p %s", credentials.Username, credentials.Password),
	}

	for _, command := range configureCommands {
		cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing command '%s': %s\n", command, output)
			return
		}
		log.Printf("Output of command '%s': %s\n", command, output)
	}
}

func runConanContainer() {
	cmd := exec.Command("docker", "run", "-d", "--name", containerName, "conan2-pollinator")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running container %s: %s\n", containerName, err)
		return
	}
	log.Printf("Output docker run: %s\n", output)
}

func stopConanContainer() {
	cmd := exec.Command("docker", "stop", containerName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error deleting container %s: %s\n", containerName, err)
		return
	}
	log.Printf("Output docker stop: %s\n", output)
}

func removeConanContainer() {
	cmd := exec.Command("docker", "rm", containerName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error deleting container %s: %s\n", containerName, err)
		return
	}
	log.Printf("Output docker remove: %s\n", output)
}

func main() {
	var credentials Credentials = readCredentials()
	buildConanImage()
	runConanContainer()
	configureConanContainer(credentials)
	stopConanContainer()
	removeConanContainer()
}
