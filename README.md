# goyaml

## Introduction

Utility for performing simple actions on YAML files.  Primarily intended to be used inside shell scripts.

## Getting the source

To get the source of this project, issue the following command:
```
$ GO111MODULE=on GOSUMDB=off go get github.com/theochva/goyaml
```

## Installing the CLI

For Mac OS, you can install the command via homebrew. Since this is hosted on the IBM Github server, you will need to specify a github token in the `HOMEBREW_GITHUB_API_TOKEN` environment variable.

```
$ brew tap theochva/grizano
$ export HOMEBREW_GITHUB_API_TOKEN=YOUR_GITHUB_TOKEN
$ brew install goyaml
```

Once already tapped and installed you can update using:
```
$ HOMEBREW_GITHUB_API_TOKEN=YOUR_GITHUB_TOKEN brew upgrade goyaml
```

Alternatively, you can download the latest release from https://github.com/theochva/goyaml/releases or install the latest using:

```
$ GO111MODULE=on GOSUMDB=off go install github.com/theochva/goyaml
```

## Running the CLI

The main syntax for the CLI is:
```
Utility to perform simple operations on YAML files:
  - get/set/delete/check properties to/from YAML content/file
  - Validate YAML content/file
  - Convert to/from YAML/JSON content/file
  - Expand Go templates using YAML as the values file

All actions can be performed using either files or stdin/stdout.

Primarily intended to be used in scripts or command line.

RC is always 0 unless there was an error while processing.

Usage:
  goyaml [command] [<flags>]

Available Commands:
  contains    Check if a value is contained in the yaml
  delete      Delete a value from the yaml
  expand      Expand Go templates using the YAML as the values data. The templates are expanded to stdout
  from-json   Convert JSON to YAML
  get         Read a value from the yaml
  help        Help about any command
  set         Set a value in a YAML document
  to-json     Convert YAML to JSON
  validate    Validate the yaml syntax

Flags:
  -f, --file string   The yaml file to read/write. If not specified it reads from stdin
  -h, --help          help for goyaml
  -v, --version       version for goyaml

Use "goyaml [command] --help" or "goyaml help [command]" for more information about a command.

Examples:
  goyaml [-f <yaml_file>] <command> [options]
  goyaml --file <yaml_file> <command> [options]
  goyaml -f <yaml_file> <command> [options]
  cat foo.yaml | goyaml <command> [options]
```

All commands can either operate on YAML read from stdin or from a file using the `--file` or `-f` options.

All commands require that the YAML file is specified using the `--file` or `-f` options, since all commands either read from or write to the YAML file. 

In addition, all commands expecting a `key` parameter accept keys with a "dot" `.` notation for nested properties.  Array support is not available.  For example, given a simple YAML file:

```yaml
parent:
  child: someValue
prop2: value2
```

The "key" `parent.child` can be used to get/set the value `someValue`

### `goyaml` Commands

The available commands are invoked in the form `goyaml -f FILE <command> [options]` and are:

#### `get`: read values from the YAML file

  - Base syntax:
    ```
    goyaml get <key> [-o|--output json|yaml]
    ```
  - Can retrieve simple or container elements
  - Can select the output format (JSON, YAML, text)
  - For more examples, see `goyaml help get` or `goyaml get --help`

#### `set`: write values to the YAML file

  - Base syntax:
    ```
    goyaml set <key> <value> [-t|--type string|int|bool|json|yaml]
    goyaml -f|--file <yaml-file> set <key> --stdin [-t|--type string|int|bool|json|yaml]
    goyaml [-f|--file <yaml-file>] set <key> -i|--input <value-file> [-t|--type string|int|bool|json|yaml]
    ```

  - The intermediate keys are created if not present in the file, e.g. using 
    ```
    goyaml -f my.yml set firstLevel.secondLevel.thirdLevel hello
    ```

    will create the path `firstLevel.secondLevel` if it does not exist before setting the value for `thirdLevel`

  - Can set values directly from CLI and specify the data type to be written to the YAML file, e.g. `int`, `bool`, `string`, `yaml` or `json`:
    ```
    goyaml -f /tmp/foo.yaml set first.second.strProp "This is a string"
    goyaml -f /tmp/foo.yaml set first.second.intProp 10 -t int
    goyaml -f /tmp/foo.yaml set first.second.boolProp true -t bool
    ```

  - Can set values from a file:
    ```
    goyaml -f /tmp/sample.yaml set first.second.third -i ~/.ssh/id_rsa_priv
    goyaml -f /tmp/sample.yaml set first.second.third -i /tmp/foo.json -t json
    goyaml -f /tmp/sample.yaml set first.second.third -i /tmp/foo.yaml -t yaml
    ```

  - Can use in a "pipe" command to read the value to set from STDIN, e.g.:
    ```
    cat  ~/.ssh/id_rsa_priv | goyaml -f /tmp/sample.yaml set first.second.third --stdin
    cat  /tmp/foo.json | goyaml -f /tmp/sample.yaml set first.second.third --stdin -t json
    cat  /tmp/foo.yaml | goyaml -f /tmp/sample.yaml set first.second.third --stdin -t yaml
    curl [options] | goyaml -f /tmp/my.yml set prop1.prop2 --stdin --type json
    ```

  - For more examples, see `goyaml help set` or `goyaml set --help`

#### `delete`: delete a value from the YAML file

  - Base syntax:
    ```
    goyaml -f|--file FILE delete <key>
    ```
  - Deletes the requested key (and subkeys) from the YAML file
  - When processing a YAML file specified with the `-f` or `--file` options, it simply outputs `true` or `false` to indicate whether the value was deleted or not.
    - Examples:
      ```
      goyaml -f /tmp/foo.yaml delete first.second.third
      goyaml -f /tmp/foo.yaml del first.second.third
      ```

  - When processing YAML read from stdin, the result (updated) YAML is printed to stdout.
    - Examples:
      ```
      cat /tmp/foo.yaml | goyaml delete first.second.third
      cat /tmp/foo.yaml | goyaml del first.second.third
      ```
  - For more exmples, see `goyaml help delete` or `goyaml delete --help`

#### `contains`: check if a value is contained in the YAML file

  - Base syntax:
    ```
    goyaml -f|--file FILE contains <key>
    ```
  - Checks if the specified key is in the YAML file and outputs either `true` or `false`
  - For more examples, see `goyaml help contains` or `goyaml contains --help`

#### `validate`: check if the specified YAML file is syntactically correct

  - Base syntax:
    ```
    goyaml -f|--file FILE validate <key> [--details,-d]
    ```
  - Default behavior is to output `true` or `false`
  - Instead of `true` or `false`, you can get the any validation error using the `--details` or `-d` flags. In this case, when YAML is valid, nothing is outputed
  - Examples:
    ```
    goyaml -f /tmp/sample.yaml validate
    cat /tmp/sample.yaml | goyaml validate
    ```
  - For more examples, see `goyaml help validate` or `goyaml validate --help`

#### `expand`: expand Go-Lang templates using the YAML file as input

  - Base syntax:
    ```
    goyaml -f|--file FILE expand -t|--template <template-files> [-e|--ext <extensions>] [-o|--output text|html]
    goyaml -f|--file FILE expand --text <template-text> [-o|--output text|html]
    ```
  - The `--output` option simply defines which Go-Lang template package to use, either `html/template` or `text/template`. For text files, the output is identical.
  - Can expand using one or more template files, in one or more directories.
  - Can expand using an inline text template
  - For more examples, see `goyaml help expand` or `goyaml expand --help`

#### `from-json`: create a YAML file from a JSON file

  - Base syntax:
    ```
    goyaml f from-json [-i|--input <input-json-file>]
    ```

  - Can be used to "convert" a JSON file to a YAML file

    - **NOTE**: some ordering might be lost in maps and arrays due to the different way maps/arrays are implemented in Go. However, the data should all be intact.

  - Can convert a JSON file into a YAML file:
    ```
    goyaml -f /tmp/foo.yaml from-json -i /tmp/foo.json
    ```
  - Can use in a "pipe" command to read the JSON value to "convert" from STDIN and print to STDOUT, e.g.:
    ```
    cat /tmp/foo.json | goyaml from-json
    ```

  - For more examples, see `goyaml help from-json` or `goyaml from-json --help`

#### `to-json`: convert the YAML file to JSON format

  - Base syntax:
    ```
    goyaml -f|--file FILE to-json [-o|--output <output-json-file>] [-p|--pretty]
    ```
  - Can be used to "convert" a YAML file to JSON
    - **NOTE**: some ordering might be lost in maps and arrays due to the different way maps/arrays are implemented in Go. However, the data should all be intact.
  - Can convert a YAML file into a JSON file:
    ```
    goyaml -f /tmp/ff.yaml to-json -p > /tmp/ff.json
    goyaml -f /tmp/ff.yaml to-json | jq .
    goyaml -f /tmp/ff.yaml to-json -o /tmp/ff.json -p
    cat /tmp/ff.yaml | goyaml to-json -p > /tmp/ff.json
    cat /tmp/ff.yaml | goyaml to-json | jq .
    cat /tmp/ff.yaml | goyaml to-json -o /tmp/ff.json -p
    ```

  - For more examples, see `goyaml help to-json` or `goyaml to-json --help`

