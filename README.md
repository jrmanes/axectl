# Piktoctl

Piktoctl is a set of tools for developers, we have different options for the tools.

---

## Installation Steps

- Copy to path 
```bash
sudo cp ./piktoctl /usr/bin/
```
- Install requirements
```bash
piktoctl sonar -i
```
- Reboot system
```bash
sudo reboot now
```
- Start the SonarQube service
```bash
piktoctl sonar -s
```
- Access in your browser and set the admin password (follow the instruction in the tool)
- Go to the parent folder of your project
- Create the projects and scan them
```bash
piktoctl sonar -c --scan -p "piktostory" -o "Piktochart"
```

---

## Sonar

Piktoctl has the command `sonar` which allows you to have a **SonarQube** in your local dev env.

There are different options for `sonar`:

- Install needed packages:
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

- Check the status of the service 
```bash
piktoctl sonar --status 
```

---

## Add to path

You can easy execute `piktoctl` adding it to some path that you have configured in your `$PATH`.

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
- [x] Setup debug argument
- [ ] Flag to specify the code coverage file
- [ ] Update from the CLI

---

Jose Ramon Mañes

---
