# adb-cli

adb-cli is yet another cli to manage Autonomous Databases in your OCI tenancy

## Quickstart Guide

The cli needs a configuration file, you can create a sample one by using the following command:

```sh
adb-cli create config
```

The command will create a file named _adb-cli.yaml.template_, please modify it to suit your OCI configuration.

The syntax is straigthforward and the template is fully documented.

The cli has the following syntax:

```sh
adb-cli [verb] [resource] [flags]
```

where

* [verb] is:
  * _start_
  * _stop_
  * _query_
  * _create_
  * _delete_
* [resource] is:
  * _database_
  * _config_
  * _wallet_
* [flags] are specific for each and every verb/resource, the most important one is __--name__ that specifies the name of the Autonomous Database to manage.

### Examples

```sh
adb-cli create database --name ATP1

adb-cli stop database --name ATP1

adb-cli delete database --name ATP1
```

Please see the full documentation [here](docs/adb-cli.md)
