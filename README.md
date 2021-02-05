# quautomatrics
Command-line automation of Qualtrics™ surveys

### Usage

```
Quautomatrics can perform basic operations on Qualtrics™ contacts, mailing lists, and distributions

Usage:
  quautomatrics [command]

Available Commands:
  createContacts     Add contacts to a mailing list
  createDistribution Create a distribution file in JSON format
  deleteContacts     Remove all contacts from a mailing list
  help               Help about any command
  listLibraries      List all libraries. A library is needed in order to create a Distribution.
  listMailingLists   Get a mailing-list ID
  listSurveys        Get a survey ID.
  replaceContacts    Replace all mailing list entries with CSV entries
  uploadDistribution Upload a distribution file

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

**Adding Contacts**
```shell
quautomatrics createContacts -m My_Fancy_Survey -c people.csv
```

**Removing All Contacts**
```shell
quautomatrics deleteContacts -m My_Fancy_Survey
```

**Replace Contacts**
* This will first remove all contacts and then add new contacts from a `CSV` file.

```shell
quautomatrics replaceContacts -m My_Fancy_Survey -c newPeople.csv
```

**Contents of quautomatrics_config.json**
```json
{
  "X-API-TOKEN": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
  "DATACENTER": "zz1",
  "fromName": "My Company Name",
  "replyToEmail": "noreply@qualtrics.com",
  "fromEmail": "noreply@qualtrics.com",
  "subject": "Please take our survey"
}
```

**Creating a Distribution**
```shell
quautomatrics.exe createDistribution -c quautomatrics_config.json -o distribution.json -l "Inquiry Survey" 
-m "Invitation Email" -n "My_Contacts"  -s "My_Fancy_Survey" -d "_NOW_" -e "_DAYS:5_T23:59:59Z"
```

**Uploading a Distribution**
```shell
quautomatrics.exe uploadDistribution -d distribution.json
```

### Date-Time Macros

These macros can be used in the `createDistribution` command with the `-d` and `-e` options:

| Macro         | Description 
|---------------|------------- 
| `_NOW_`       | replaced with current date/time such as `2006-01-02T15:04:05Z` |
| `_TODAY_`     | replaced with current date such as `2006-01-02` |
| `_YMD_`       | same as `_TODAY_` |
| `_HMS_`       | replaced with current time such as `15:04:05` | 
| `_TOMORROW_`  | replaced with tomorrow's date such as `2006-01-03` | 
| `_DAYS:n_`    | replaced with *n* days into the future; when n=3 then `2006-01-05` |

### License
* [MIT License](https://github.com/jftuga/quautomatrics/blob/main/LICENSE)

### API

This program **does not** use the XM Directory API. Instead, it uses the `Research Core Contacts` API.

APIs used:

* [Create Distribution](https://api.qualtrics.com/instructions/reference/distributions.json/paths/~1distributions/post)
* [List Libraries](https://api.qualtrics.com/instructions/reference/libraries.json/paths/~1libraries/get)
* [List Library Messages](https://api.qualtrics.com/instructions/reference/libraries.json/paths/~1libraries~1%7BlibraryId%7D~1messages/get)
* [List Mailing Lists](https://api.qualtrics.com/instructions/reference/researchCore.json/paths/~1mailinglists/get)
* [Update Mailing List](https://api.qualtrics.com/instructions/reference/researchCore.json/paths/~1mailinglists~1%7BmailingListId%7D/put)
* [List Contacts](https://api.qualtrics.com/instructions/reference/researchCore.json/paths/~1mailinglists~1%7BmailingListId%7D~1contacts/get)
* [Delete Contacts](https://api.qualtrics.com/instructions/reference/researchCore.json/paths/~1mailinglists~1%7BmailingListId%7D~1contacts~1%7BcontactId%7D/delete)
* [Create Contacts](https://api.qualtrics.com/instructions/reference/researchCore.json/paths/~1mailinglists~1%7BmailingListId%7D~1contacts/post)

### Acknowledgements
* [cobra](https://github.com/spf13/cobra)
* [viper](https://github.com/spf13/viper)
* [jsonparser](https://github.com/buger/jsonparser)
* **Please note that there is no official affiliation between Qualtrics™ and `quautomatrics`**
