# Piktoctl

Piktoctl is a set of tools for developers, we can create different commands to abstract manual tasks.

This tool is written in [Go](https://go.dev/) with the [cobra](https://github.com/spf13/cobra) framework.

As result of the project, we obtain a binary which can be compiled for the different platforms and architectures.

---

## Sonar

Piktoctl has the command `sonar` which allows you to setup and configure a **SonarQube** in your local dev env.

This command is an abstraction to setup and configure the project in **SonarQube**.

### What it does?

The CLI with the command `sonar` uses [docker](https://www.docker.com/), [docker-compose](https://docs.docker.com/compose/) and the [sonar-scanner](https://docs.sonarqube.org/latest/analysis/scan/sonarscanner/)

Some [features](#features) are:

- Use two docker containers:`SonarQube` and `PostgresSQL`
- Use the `sonar-scanner` to analyze the code
- The tool `piktoctl` communicates to the `SonarQube` API to create the projects and the tokens automatically, the tokens are store in the path `~/.piktoctl/sonar/token`
- It downloads the `sonar-scanner` automatically and place it in `~/.sonar-scanner-4.6.2.2472-linux/`, however, you can do it manually
- `piktoctl` has the flag `-i` which install you all the needed requirements for you, the requirements are needed to run `Docker` in your computer and execute `sonar-scanner`, the requirements are:
  - wget
  - unzip
  - docker
  - docker-compose
  - openjdk-11-jre-headless
  - default-jre
  - default-jdk
- The flag `-i` also add your user to the `Docker` group.
- It asks you to restart your computer for changes to take effect.

### Examples

- Install requirements, this step install all the requirements to run `Docker` and the `sonar-scanner` in your computer. [features](#features)
- After the installation, you will be prompt to restart your computer, this is because is needed after add your user to the `Docker` group [source](https://docs.docker.com/engine/install/linux-postinstall/)
```bash
piktoctl sonar -i
```

- Start the service, creating the projects and scan them
```bash
piktoctl sonar -s -c --scan -p "someProject" -o "someOrganization"
```

- Start the service
```bash
piktoctl sonar -s -p "someProject" -o "someOrganization"
```

- Start the service creating the projects
```bash
piktoctl sonar -s -c -p "someProject" -o "someOrganization"
```

- With SonarQube running, create the projects and scan them
```bash
piktoctl sonar -c --scan -p "piktostory" -o "Piktochart"
```

- Check the status of the service
```bash
piktoctl sonar --status 
```

- Start the SonarQube service
```bash
piktoctl sonar -s
```

---

### Sonar-scanner Docker


Base command

```bash
	
docker run \
      --rm \
      --network=tmp_sonar \
      -e SONAR_HOST_URL="http://sonarqube:9000" \
      -v /System/Volumes/Data/Volumes/Data/workspace/dev/go/github.com/jrmanes/piktoctl:/root/src sonarsource/sonar-scanner-cli \
      -Dsonar.projectKey=piktoctl \
      -Dsonar.sonar.projectName=piktoctl \
      -Dsonar.sonar.projectVersion=1.0 \
      -Dsonar.scm.disabled=true \
      -Dsonar.sources=./ \
      -Dsonar.sonar.host.url=http://sonarqube:9000 \
      -Dsonar.login=`+token
```

Command executed inside the tool

```bash
docker run \
      --rm \
      --network=tmp_sonar \
      -e SONAR_HOST_URL="http://sonarqube:9000" \
      -v /home/joseramon/Tools/piktostory/:/root/src sonarsource/sonar-scanner-cli \
      -Dsonar.projectKey=`+p+` \
      -Dsonar.sonar.projectName=`+p+` \
      -Dsonar.sonar.projectVersion=1.0 \
      -Dsonar.scm.disabled=true \
      -Dsonar.sources=./ \
      -Dsonar.sonar.host.url=http://sonarqube:9000 \
      -Dsonar.login=`+token
```

---

## Add to path

You can easy execute `piktoctl` anywhere adding the binary to some path that you have configured in your `$PATH`.

You can use for instance the path:
`/usr/bin/`

Execute the command:
```bash
sudo cp ./piktoctl /usr/bin/
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
- [ ] Flag to specify the code coverage file
- [ ] Update release from the CLI
- [ ] Refactor

---