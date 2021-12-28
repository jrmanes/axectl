/*
Copyright © 2021 Jose Ramon Mañes jr.mb47@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Command struct which contains an info message, command to execute and an array of arguments
type Command struct {
	message string
	command string
	args    []string
}

// Commands list of commands
type Commands []Command

// sonarCmd represents the sonar command
var sonarCmd = &cobra.Command{
	Use:   "sonar",
	Short: "SonarQube command options",
	Long: `Using the command 'sonar' 
You will be able to configure a SonarQube with docker for local development.
Start the container.
Scan projects.`,
	// Validate if there is any flag added, if not, we send the user to Usage func
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flag.Parse()
		tail := flag.Args()
		if len(tail) <= 1 {
			cmd.Usage()
			os.Exit(0)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		organization, _ = cmd.Flags().GetString("organization")
		project, _ = cmd.Flags().GetString("project")
		debug := cmd.Flags().Changed("debug")

		// check if the user and password were provided
		if cmd.Flags().Changed("user") {
			u, _ = cmd.Flags().GetString("user")
			userData := strings.Split(u, ":")
			fmt.Println("[INFO] New user config:", u)
			fmt.Println("[INFO] New user data:", userData)

			sonarUser = userData[0]
			sonarPass = userData[1]

			fmt.Println("[INFO] New user config:", sonarUser, sonarPass)
		}
		if cmd.Flags().Changed("install") {
			install(debug)
		}
		if cmd.Flags().Changed("run") {
			run()
		}
		if cmd.Flags().Changed("organization") || cmd.Flags().Changed("project") {
			if organization == "" || project == "" {
				fmt.Println("[ERROR] Organization needs to be set, use parameter: -o ")
			}
		}
		if cmd.Flags().Changed("create") {
			createProject()
			createProjectToken()
		}
		if cmd.Flags().Changed("scan") {
			scan()
		}
		if cmd.Flags().Changed("status") {
			status()
		}
		if cmd.Flags().Changed("stop") {
			stop()
		}
	},
}

var (
	filePath                 = "/tmp/"
	fileName                 = "docker-compose.piktochart-sonarqube"
	sonarUser                = "admin"
	sonarPass                = "admin123."
	project, organization, u string
	tokensFolder             = "/.piktoctl/sonar/tokens/"
)

func init() {
	rootCmd.AddCommand(sonarCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sonarCmd.PersistentFlags().String("foo", "", "A help for foo")
	// sonarCmd.Flags().String("", "", "A help for foo")

	sonarCmd.PersistentFlags().BoolP("install", "i", true, "[*] TODO: Install all requirements needed")
	sonarCmd.PersistentFlags().BoolP("scan", "", true, "Scan a project")

	sonarCmd.PersistentFlags().BoolP("create", "c", true, "Create a project and tokens")
	rootCmd.PersistentFlags().StringP("organization", "o", "", "Organization in SonarQube")
	rootCmd.PersistentFlags().StringP("project", "p", "", "You can add one project name or multiple separated by comas.")
	sonarCmd.PersistentFlags().BoolP("run", "r", true, "Start running the SonarQube container")
	sonarCmd.PersistentFlags().BoolP("stop", "", true, "Stop the SonarQube container")
	sonarCmd.PersistentFlags().BoolP("status", "", true, "Check the docker container status")
	sonarCmd.PersistentFlags().StringP("user", "u", "admin:admin123.", "Use your user:password  -> Example: admin:admin123.")

	sonarCmd.PersistentFlags().BoolP("debug", "", false, "Set debug option")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sonarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// install the needed software
func install(debug bool) {
	// TODO: allow installation for MacOS & Windows
	switch os := detectOS(); os {
	case "darwin":
		fmt.Println("TODO: Development pending...")
	case "linux":
		LinuxPkg(debug)
	default:
		LinuxPkg(debug)
	}
	// Install Linux Requirements
}

// LinuxPkg Install needed SonarQube packages for Linux environments
func LinuxPkg(debug bool) {
	// create a list with all the packages needed
	packages := []string{
		"docker",
		"docker-compose",
		"wget",
		"unzip",
		"openjdk-11-jre-headless",
		"default-jre",
		"default-jdk",
	}

	// Update the package list
	fmt.Println("[INFO] 📦 Update package list... ")
	cmd := exec.Command("sudo", "apt", "update")
	if debug {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Loop inside all packages and install them one by one
	for _, p := range packages {
		fmt.Println("[INFO] 📦 Installing package: ", p)

		cmd := exec.Command("sudo", "apt", "install", "-y", p)
		if debug {
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout

		}
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Configure system
	linuxSystem(debug)
}

// linuxSystem Configure system to execute SonarQube in Linux
func linuxSystem(debug bool) {
	// get current user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := user.HomeDir

	// lco list of commands to execute
	commands := Commands{Command{
		message: "[INFO] 📦 Download package Sonar Scanner... ",
		command: "wget",
		args:    []string{"https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.6.2.2472-linux.zip", "-O", "/tmp/sonar-scanner.zip"},
	}, Command{
		message: "[INFO] 📦 Unzip package Sonar Scanner... ",
		command: "unzip",
		args:    []string{"-o", "/tmp/sonar-scanner.zip", "-d", "/tmp/"},
	}, Command{
		message: "[INFO] 📦 Copy package Sonar Scanner to ~/",
		command: "cp",
		args:    []string{"-R", "/tmp/sonar-scanner-4.6.2.2472-linux/", home + "/"},
	}, Command{
		message: "[INFO] 📦 Clean folder if exists in: " + home,
		command: "rm",
		args:    []string{"-rf", home + "/.sonar-scanner-4.6.2.2472-linux/"},
	}, Command{
		message: "[INFO] 📦 Hide folder package Sonar Scanner in: " + home,
		command: "mv",
		args:    []string{home + "/sonar-scanner-4.6.2.2472-linux/", home + "/.sonar-scanner-4.6.2.2472-linux"},
	}, Command{
		message: "[INFO] 📦 Copy package Sonar Scanner to /usr/local/bin",
		command: "sudo",
		args:    []string{"cp", home + "/.sonar-scanner-4.6.2.2472-linux/bin/sonar-scanner", "/usr/local/bin/"},
	}, Command{
		message: "[INFO] 📦 Copy library Sonar Scanner to /usr/local/lib",
		command: "sudo",
		args:    []string{"cp", home + "/.sonar-scanner-4.6.2.2472-linux/lib/sonar-scanner-cli-4.6.2.2472.jar", "/usr/local/lib/"},
	}, Command{
		message: "[INFO] 📦 Remove java from sonar-scanner",
		command: "rm",
		args:    []string{home + "/.sonar-scanner-4.6.2.2472-linux/jre/bin/java"},
	}, Command{
		message: "[INFO] 📦 Copy java from system",
		command: "ln",
		args:    []string{"-s", "/usr/bin/java", home + "/.sonar-scanner-4.6.2.2472-linux/jre/bin/java"},
	}, Command{
		message: "[INFO] 📦 Add docker group to the user: " + user.Username + "",
		command: "sudo",
		args:    []string{"usermod", "-aG", "docker", user.Username},
	},
	}

	// loop inside all the commands to execute
	for _, c := range commands {
		fmt.Println(c.message)
		c := exec.Command(c.command, c.args...)
		if debug {
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
		}
		err = c.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Setenv("JAVA_HOME", "/usr/lib/jvm/java-11-openjdk")

	fmt.Println("[INFO] 📦 Cleaning temporary path...")
	cmd := exec.Command("rm", "-fr", "/tmp/sonar-scanner-4.6.2.2472-linux/", "/tmp/sonar-scanner.zip")
	if debug {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[INFO] ✅ All packages have been installed successfully!")
	fmt.Println("[INFO] 🔄 Please restart your computer to execute: sonar-scanner")
}

// status check the status of the containers
func status() {
	cmd := exec.Command("docker", "ps", "-a")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

// scan check every project and scan it on sonar
func scan() {
	now := time.Now()

	fmt.Println("[INFO] 🔭 Scanning projects...")

	projects := strings.Split(project, ",")

	// crate the project in SQ
	for _, p := range projects {
		fmt.Println("[INFO] 🔭 Scanning project...", p)
		out, err := exec.Command("pwd").Output()
		if err != nil {
			log.Fatal(err)
		}
		path := string(out)

		strings.Replace(path, `\n`, "\n", -1)

		// get token value if exists
		token, err := GetTokenInFile(p)
		if err != nil {
			log.Fatal(err)
		}

		// removed last character - new line
		path = path[:len(path)-1]

		err = SonarScanner(p, token)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("[INFO] ---------------------------- ")
	fmt.Println("[INFO] Elapse: ", time.Since(now))
	fmt.Println("[INFO] ---------------------------- ")
}

// SonarScanner executes the scanner of code
func SonarScanner(p, token string) error {
	// TODO: check if docker execution worth it
	//command := `docker run \
	//--rm \
	//--network=tmp_sonar \
	//-e SONAR_HOST_URL="http://sonarqube:9000" \
	//-v /home/joseramon/Tools/piktostory/:/root/src sonarsource/sonar-scanner-cli \
	//-Dsonar.projectKey=`+p+` \
	//-Dsonar.sonar.projectName=`+p+` \
	//-Dsonar.sonar.projectVersion=1.0 \
	//-Dsonar.scm.disabled=true \
	//-Dsonar.sources=./ \
	//-Dsonar.sonar.host.url=http://sonarqube:9000 \
	//-Dsonar.login=`+token

	// chekk if command: sonar-scanner exists in path
	exists := commandExists("sonar-scanner")
	if !exists {
		log.Fatal(exists)
	}

	// get current user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	home := string(user.HomeDir)

	command := home + "/.sonar-scanner-4.6.2.2472-linux/bin/" + `sonar-scanner  \
-Dsonar.host.url=http://localhost:9000 \
-Dsonar.projectKey=` + p + ` \
-Dsonar.sonar.projectName=` + p + ` \
-Dsonar.sources=./` + p + ` \
-Dsonar.exclusions="**/node_modules/**" \
-Dsonar.inclusions="**,**/*.js" \
-Dsonar.tests.inclusions="src/**/*.spec.js,src/**/*.spec.jsx,src/**/*.test.js,src/**/*.test.jsx,**/__tests__/**,e2e/**" \
-Dsonar.login=` + token

	fmt.Println("[INFO] ---------------")
	fmt.Println("command: ", command)
	fmt.Println("[INFO] ---------------")
	//
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// run configure and initialize the containers
func run() {
	fmt.Println("[INFO] 🚢 We are starting the setup process... this can take some seconds...")
	configureSystem()

	dockerComposeFile := dockerComposeFile()
	fileName := CreateFileWithContent(filePath+fileName, dockerComposeFile)

	cmd := exec.Command("docker-compose", "-f", filePath+fileName, "up", "-d")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal("Please, check that your current user has permissions to execute docker \n"+
			"you can add it to the docker group executing: sudo usermod -aG docker $USER \n"+
			"After that, please, reboot your system \n", err)
	}

	// Grant enough time to allow the service start
	fmt.Println("[INFO] 🚢 SonarQube is starting, wait some seconds...")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕐")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕑")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕓")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕔")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕕")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕖")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕗")
	time.Sleep(1 * time.Second)
	fmt.Println("[INFO] 🕘")
	time.Sleep(1 * time.Second)

	fmt.Println("[INFO] 🚧 Please, open the following link and change the password when the service will be up")
	fmt.Println("[INFO] 👤 Default user [" + sonarUser + ":admin]")
	fmt.Println("[INFO]  http://localhost:9000/")
	fmt.Println("[INFO] 🚨 RECOMMENDATION: Change the password to: " + "[" + sonarPass + "] ")
	fmt.Println("[INFO] Press enter once you have change the password. ")

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	_, err = buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("[INFO] 🙉 SonarQube is up an running!")
	time.Sleep(3 * time.Second)
}

// stop the docker-compose containers
func stop() {
	fmt.Println("[INFO] Stopping SonarQube...")
	cmd := exec.Command("docker-compose", "-f", filePath+fileName, "stop")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println("[INFO] 👋 SonarQube is stopped!")
}

// Util functions

// configureSystem set the needed path to the sysctl
func configureSystem() {
	o := detectOS()
	if o == "linux" {

		fmt.Println("[INFO] Starting the containers, it can take a while...")
		fmt.Println("[INFO] 🔧 We are going to configuring system, we need to increase the configuration of: sysctl to -> vm.max_map_count=262144")
		fmt.Println("[INFO] 🔓 We need to run as ROOT...")

		cmd := exec.Command("sudo", "sysctl", "-w", "vm.max_map_count=262144")

		cmd.Stdin = os.Stdin
		//cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

}

// dockerComposeFile returns all the data inside the docker-compose.yml file
func dockerComposeFile() string {
	dockerFile := `
version: "3"
services:
  sonarqube:
    image: sonarqube:9.2-community
    expose:
      - 9000
    ports:
      - "9000:9000"
    networks:
      - sonar
    environment:
      - sonar.jdbc.username=sonar
      - sonar.jdbc.password=sonar
      - sonar.jdbc.url=jdbc:postgresql://psql:5432/sonar
  psql:
    image: postgres:9.5
    networks:
      - sonar
    ports:
      - 5432:5432
    environment:
      - POSTGRES_USER=sonar
      - POSTGRES_PASSWORD=sonar
      - POSTGRES_DATABASE=sonar
    volumes:
      - postgresql:/var/lib/postgresql
      - postgresql_data:/var/lib/postgresql/data
networks:
  sonar:
volumes:
  postgresql_data:
  postgresql:
`
	return dockerFile
}

// CreateFileWithContent generates the docker file in the path specified
func CreateFileWithContent(path, content string) string {
	// create file
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write content
	_, err2 := f.WriteString(content)
	if err2 != nil {
		log.Fatal(err2)
	}

	return fileName
}

// createProject generates the project in SonarQube
func createProject() {
	fmt.Println("[INFO] --------------------------------------------------------------")
	fmt.Println("[INFO] 💡 The organization to create the project is: ", organization)
	fmt.Println("[INFO] --------------------------------------------------------------")

	projects := strings.Split(project, ",")
	// crate the project in Sonar
	for _, p := range projects {
		fmt.Println("[INFO] 📚 Project to create: ", p)

		params := url.Values{}
		params.Add("project", p)
		params.Add("organization", organization)
		params.Add("name", p)
		body := strings.NewReader(params.Encode())

		req, err := http.NewRequest("POST", "http://localhost:9000/api/projects/create", body)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(sonarUser, sonarPass)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// TokenResponse is the struct that we use for our Sonar responses
type TokenResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Token     string `json:"token"`
	CreatedAt string `json:"createdAt"`
}

// createProjectToken generates the token for the project in SonarQube
func createProjectToken() {
	fmt.Println("[INFO] --------------------------------------------------------")
	fmt.Println("[INFO] 💡 The organization to create the token is: ", organization)
	fmt.Println("[INFO] --------------------------------------------------------")

	projects := strings.Split(project, ",")
	// crate the project in SQ
	for _, p := range projects {
		fmt.Println("[INFO] 💡 Project to create the token: ", p)
		// Get info from the actual tokens configuration
		token, err := GetTokenInFile(p)
		if err != nil && token == "" {
			fmt.Println("[INFO] ✔️ Creating new token for project: ", p)

			params := url.Values{}
			params.Add("name", p)
			body := strings.NewReader(params.Encode())

			req, err := http.NewRequest("POST", "http://localhost:9000/api/user_tokens/generate", body)
			if err != nil {
				log.Fatal(err)
			}
			req.SetBasicAuth(sonarUser, sonarPass)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			// Check response, if it's ok, store the token into the FS
			if resp.StatusCode == 200 {
				token := TokenResponse{}
				err2 := json.NewDecoder(resp.Body).Decode(&token)
				if err2 != nil {
					log.Fatal(err)
				}
				fmt.Println("[INFO]: ", token.Name, "=", token.Token)

				home, err := os.UserHomeDir()
				cobra.CheckErr(err)

				// Store the content insisde ~/.piktoctl/sonar/tokens/
				configHome := filepath.Join(home, tokensFolder)
				fileInPath := filepath.Join(configHome, token.Name)

				err = CreateFileInPath(configHome, fileInPath)
				if err != nil {
					log.Fatal(err)
				}

				tokenFile := filepath.Join(configHome, token.Name)
				CreateFileWithContent(tokenFile, token.Token)
			} else {
				fmt.Println("[ERROR] Failed token creation, it's possible that the token already exists, for check it, got to:")
				fmt.Println("[ERROR] Try to check the token in your path: ~/.piktoctl/sonar/tokens/ - or check it in the panel:")
				fmt.Println("[ERROR] http://localhost:9000/account/security")
			}
		} else {
			fmt.Println("[INFO] 📜️ Using existing token for project: ", p)
		}
		fmt.Println("[INFO] --------------------------------------------------------")
	}
}

// GetTokenInFile check the content inside the file and return it
func GetTokenInFile(tokenName string) (string, error) {
	// Get current user directory
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Println("user home dir not found...")
		return "", err
	}
	tokenValue := dirname + "/.piktoctl/sonar/tokens/" + tokenName

	t, err2 := ioutil.ReadFile(tokenValue) // just pass the file name
	if err2 != nil {
		return "", err2
	}
	token := string(t)

	return token, nil
}

// commandExists verify if a command exists in path
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// detectOS check the current OS where the tool is being executed
func detectOS() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	default:
		return "linux"
	}
	return ""
}
