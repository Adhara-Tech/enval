[![CircleCI](https://circleci.com/gh/Adhara-Tech/enval.svg?style=shield)](https://circleci.com/gh/Adhara-Tech/enval) 
[![<Sonarcloud quality gate>](https://sonarcloud.io/api/project_badges/measure?project=citrus&metric=alert_status)](https://sonarcloud.io/dashboard?id=citrus)

# Enval

Enval is an environment validator that checks if your system is ready to work with a given project. It helps development teams to avoid different behaviors when using different versions of the same tools and a quick way to identify all the tools you need to install to work in a project. It also helps to ensure CI pipelines are using the same version of the tools that the development team is using. It is easy to integrate and has a nice output


# Demo

<p align="center">
  <img src="./docs/demo.svg" width="100%">
</p>


# Install

Download enval binary from [releases page](https://github.com/Adhara-Tech/enval/releases) and place in your path

check if it correctly working executing

```
enval version
```

We recommend to update enval frequently because new features and tools specifications are being added

# Quick start

To run enval execute inside your project directory:
```
enval
```

a valid enval manifest is needed to perform the validations



# Configuring enval manifest

To setup a project to use enval create a file called ```.enval.manifest.yaml``` in the root directory like this

```yaml
# Name of the tool
name: my-cool-app-enval-manifest
# Directory where custom tool specs are placed (see how to define a tool section)
custom_specs: my-specs/tools
# List of tools to validate
tools:
    # Name of the tool (must be aligned with the name used in the tool specification)
    # Once configured enval will check that the tool is in path and can be executed
    # See available tools section for tools and fields
  - name: node
    # Checks will be performed next checks with the version command
    checks:
      # checks are defined in the format fieldName: "constraint"
      # for semver field constraints are defined with next symbols
      #       "<" minor than 
      #       ">" greater than
      #       "="equal to
      # Version value can omit minor and or patch. Valid values are
      #       ">= 12" that is equivalent to ">= 12.0.0"
      #       ">= 12.*.1" that requires a version 12 with patch version greater or equal to 1
      #       ">= 12.3" that is equivalent to ">= 12.3.0"
      # In this example the version for node tool must be greater or equal to 12.13.1
      version: ">= 12.13.1"
```


There are some tools with different implementations and sometimes it is needed to ensure that one specific implementation is present in the system. Enval uses the concept ```flavor``` to represent this scenario.


Following manifest requires openjdk flavor for java tool

```yaml
name: example-application-manifest
custom_specs: my-specs/tools
tools:
  - name: java
    flavor: openjdk
    checks:
      version: ">= 11.0.0"
```

While this manifest requires java (without any implementation preference) so openjdk, oracle, ... are valid

```yaml
name: example-application-manifest
custom_specs: my-specs/tools
tools:
  - name: java
    checks:
      version: ">= 11.0.0"
```

When a flavored tool is configured without requesting an specific flavor all the flavors will be considered to validate the check following next process
* For every flavor
* Execute version checker parsers. If one of the match use that info to validate the checks
* if The checks are satisfied the tool will be marked as available and valid
* if The parsers don't match or match but the field checks are not satisfied enval continues with the next flavor
* If no more flavors the tool will be marked as invalid  

Read next section for a complete list of tools, flavors and fields

# Available tools

Tools that are shipped with enval. If you need a tool not available you can define custom tools (see section how to configure custom tools) and you can open a pr with your specification to allow the enval community benefit from it


### [Go language](https://golang.org/)

Tool Name: **go**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

### [GolangCI-Lint](https://github.com/golangci/golangci-lint)

Tool Name: **golanci-lint**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |


### [gotestsum](https://github.com/gotestyourself/gotestsum)

Tool Name: **gotestsum**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

### Java

Tool Name: **java**
Flavors: yes

#### [Java Openjdk](https://openjdk.java.net/)

Flavor Name: **openjdk**

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

#### [Java Oracle](https://www.oracle.com/java/)

Flavor Name: **oracle**

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |



### [Node.js](https://nodejs.org/)

Tool Name: **node**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

### [npm](https://www.npmjs.com/)

Tool Name: **npm**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |


### [Open api generator](https://github.com/OpenAPITools/openapi-generator)

Tool Name: **openapi-generator**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

### [swagger-cli](https://www.npmjs.com/package/swagger-cli)

Tool Name: **swagger-cli**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |


### [Terraform](https://www.terraform.io/)

Tool Name: **terraform**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

### [Truffle suite](https://www.trufflesuite.com/)

Tool Name: **truffle**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |
| core    | semver            |
| solidity    | semver            |

### [PHP](https://www.php.net/)

Tool Name: **PHP**
Flavors: No

| Field        | Kind  |
| ------------- |:-------------|
| version    | semver            |

# How to define a new tool or override an existing one

Tools are defined in a tool specifications yaml files. To add a custom validation or override an existing one just place a specification yaml inside the directory configured in custom_specs in the ```.enval.manifest.yaml```

To override a tool just use the same name as the enval defined tool.


## Basic specification

Most of the tools are easy to configure to be supported by enval by just deciding the name of the tool and having a regular expression to validate it. 

Following specification example demonstrates how to configure a tool of this type.

```yaml
# Name of the tool that will be used in manifests
name: go
# Description of the tool
description: go language compiler and tools
# Command that will be checked. It must be in system's path and with execution permissions
command: go
# Arguments needed to get the version from the command. It can be overridden by flavors
version_command_args:
  - version
# Version checkers defines how to parse the command version output and extract the needed data that will be checked based on the manifest specification. It can be overridden by a flavor.
version_checker:
  # parsers that will try to parse the output. They are executed in order until one of them matches the output
  parsers:
      # type of the parser. Now only regexp is supported
    - type: regexp
      # regexp following golang regexp syntax to parse the output. Regexp groups must match the fields section. In this case version group
      regexp: 'go version go(?P<version>\d+\.\d+\.?\d*)'
  # Fields that will be available to define checks in the manifest
  fields:
      #Name of the field
    - name: version
      #Type of the field. Now supported semver that allows semantic version checks and fixed that allow exact match validations
      type: semver
      # Configures how to deal with a field not found by the parser. If it is required and the parser does not return, enval will produce an error
      required: true
```


## Advanced specifications

In some cases the tools are not so simple and they:
* have different flavors
* The version commands returns different information depending on versions 
* There are different fields that are interesting to be checked
* ...

So more complex specifications are needed. 

This section covers those scenarios

**Specification with flavors**

Typically flavors need to define custom version_checkers. When flavors are configured the tool level version_checkes is still valid but will be overridden by the flavor so if a standard behavior can be defined it can be reused for different flavor only redefining those ones with different behavior.

Following example configures 2 flavors. First one with 2 parsers one for the GA versions and the second one for de EA versions

```yaml
name: java
description: java
command: java
version_command_args:
  - --version
# Flavors supported for this tool
flavors:
  - name: openjdk
    # Version checker for the openjdk flavor
    version_checker:
      # Two parsers are needed for this flavor one for ga versions and another for ea versions. 
      # This two parsers are executed in order and the first matching will stop the loop. 
      # Both parsers provides version field so it was configured as mandatory
      parsers:
        - type: regexp
          regexp: 'openjdk (?P<version>\d+\.\d+\.\d+)' #ga versions
        - type: regexp
          regexp: 'openjdk (?P<version>\d+)-ea' #ea versions
      fields:
        - name: version
          type: semver 
          required: true 
  - name: oracle
    version_checker:
      # This flavor 
      parsers:
        - type: regexp
          regexp: 'java (?P<version>\d+\.\d+\.\d+)'
      fields:
        - name: version
          type: semver
          required: true

```


**Specification with two or more fields**

Sometimes the version command output returns data that is needed to validate if the version is correct and that data can not be aggregated in one field. In those cases different fields are configured like the following example

```yaml
name: truffle
description: Truffle ethereum suite
command: truffle
version_command_args:
  - version
version_checker:
  parsers:
    - type: regexp
      regexp: 'Truffle v(?P<version>\d+\.\d+\.?\d*) \(core: (?P<core>\d+\.\d+\.?\d*)\)\sSolidity v(?P<solidity>\d+\.\d+\.?\d*)'
  fields:
    - name: version
      type: semver
      required: true
    - name: core
      type: semver
      required: true
    - name: solidity
      type: semver
      required: true
```
