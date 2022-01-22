# JRCTL

Jrctl is a set of tools for developers, we can create different commands to abstract manual tasks.

This tool is written in [Go](https://go.dev/) with the [cobra](https://github.com/spf13/cobra) framework.

As result of the project, we obtain a binary which can be compiled for the different platforms and architectures.

---

## Build project

To build the project from your side, you need `Go` installed in your computer, and execute:

Remember to check your `OS` and your architecture.

```
mkdir -p ./bin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/jrctl ./
```

---

## Sonar

Jrctl has the command `sonar` which allows you to setup and configure a **SonarQube** in your local dev env.

This command is an abstraction to setup and configure the project in **SonarQube**.

### What it does?

The CLI with the command `sonar` uses [docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/) to create the resources needed.

Some [features](#features) are:

- Use three docker containers:`SonarQube`, `PostgreSQL` and `sonar-scanner`.
  - `SonarQube` -> server
  - `PostgreSQL` -> database engine
  - `sonar-scanner` -> tool from Sonar to analyse the code
- The tool `jrctl` communicates to the `SonarQube` API to create the projects and the tokens automatically, the tokens are store in the path `~/.jrctl/sonar/token`
- `jrctl` has the flag `-i` which install the needed requirements for you, the requirements are:
  - docker
  - docker-compose
- The flag `-i` also add your user to the `Docker` group.
- It asks you to restart your computer for changes to take effect.

### Examples

- Install requirements, this step install all the requirements to run `Docker` and the `sonar-scanner` in your computer. [features](#features)
- After the installation, you will be prompt to restart your computer, this is because is needed after add your user to the `Docker` group [source](https://docs.docker.com/engine/install/linux-postinstall/)
```bash
jrctl sonar -i
```

- Start the service, creating the projects and scan them
```bash
jrctl sonar -s -c --scan -p "someProject" -o "someOrganization"
```

- Start the service
```bash
jrctl sonar -s -p "someProject" -o "someOrganization"
```

- Start the service creating the projects
```bash
jrctl sonar -s -c -p "someProject" -o "someOrganization"
```

- With SonarQube running, create the projects and scan them
```bash
jrctl sonar -c --scan -p "piktostory" -o "Piktochart"
```

- Check the status of the service
```bash
jrctl sonar --status 
```

- Start the SonarQube service
```bash
jrctl sonar -s
```

---

### Sonar-scanner Docker


You can execute the following command in order to run the analysis directly.

```bash
docker run \
      --rm \
      --network=host \
      -e SONAR_HOST_URL="http://sonarqube:9000" \
      -v $PWD:/usr/src sonarsource/sonar-scanner-cli \
      -Dsonar.projectKey=jrctl \
      -Dsonar.sonar.projectName=jrctl \
      -Dsonar.sonar.projectVersion=1.0 \
      -Dsonar.scm.disabled=true \
      -Dsonar.sources=./ \
      -Dsonar.sonar.host.url=http://sonarqube:9000 \
      -Dsonar.login=`+token
```

---

## Add to path

You can easily execute `jrctl` anywhere adding the binary to some path that you have configured in your `$PATH`.

You can use for instance the path:
`/usr/bin/`

Execute the command:
```bash
sudo cp ./jrctl /usr/bin/
```

---

## TODO
- [x] Install dependencies
- [x] Start all the containers
- [x] Check containers status
- [x] Generate project in Sonar
- [x] Generate token for project
- [x] Create config path
  - [x] Create tokens inside the config path
  - [x] List tokens
  - [ ] Delete tokens
  - [ ] Delete all resources created
- [x] Setup debug argument
- [x] Update release from the CLI
- [x] Remove installation packages for `sonar-scanner`
- [x] Refactor
- [ ] Flag to specify the code coverage file

---

Jose Ramon Mañes

---