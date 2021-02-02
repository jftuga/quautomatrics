# quautomatrics
Command-line automation of Qualtrics surveys

### Usage

```
Quautomatrics can perform basic operations on Qualtrics contacts and mailing lists

Usage:
  quautomatrics [command]

Available Commands:
  createContacts  Add contacts to a mailing list
  deleteContacts  Remove all contacts from a mailing list
  help            Help about any command
  replaceContacts Replace all mailing list entries with CSV entries

Flags:
      --config string   config file (default is quautomatrics_config.json)
  -h, --help            help for quautomatrics
  -v, --version         version for quautomatrics

Use "quautomatrics [command] --help" for more information about a command.
```

### Examples

* Rename `quautomatrics_config-dist.json` to `quautomatrics_config.json`
* * Edit these fields: `X-API-TOKEN`, `DATACENETR`
* A `CSV` file must use this format, with no header line:
* * `first name,last name,email address`
    
A mailing must already be created through the Qualtrics web interface.
In the examples below it is called `My_Fancy_Survey`.

*Adding Contacts*
```shell
quautomatrics createContacts -m My_Fancy_Survey -c people.csv
```

*Removing All Contacts*
```shell
quautomatrics deleteContacts -m My_Fancy_Survey
```

*Replace Contacts*
* This will first remove all contacts and then add new contacts from a `CSV` file.

```shell
quautomatrics replaceContacts -m My_Fancy_Survey -c newPeople.csv
```
### License
* [MIT License](https://github.com/jftuga/quautomatrics/blob/main/LICENSE)

### Acknowledgements
* [cobra](https://github.com/spf13/cobra)
* [viper](https://github.com/spf13/viper)
* [jsonparser](https://github.com/buger/jsonparser)

