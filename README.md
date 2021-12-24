# Piktoctl

Piktoctl is a set of tools for developers, we have different options for the tools.

---

## Sonar

Piktoctl has the command `sonar` which allows you to have a **SonarQube** in your local dev env.

There are different options for `sonar`:

- Install needed packages:
```bash
piktoctl sonar -i
```

- Start the service
```bash
piktoctl sonar -r -p "someProject" -o "someOrganization"
```

- Start the service creating the projects
```bash
piktoctl sonar -r -c -p "someProject" -o "someOrganization"
```

- Check the status of the service 
```bash
piktoctl sonar --status 
```

---

## TODO
- [x] Install dependencies
- [x] Execute docker-compose up
- [x] Check containers status
- [x] Generate project in Sonar
- [x] Generate token for project
- [x] Create config path
  - [x] Create tokens inside the config path
  - [x] List tokens
  - [ ] Delete tokens
- [x] Setup debug argument

---

