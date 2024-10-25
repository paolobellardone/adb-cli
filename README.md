# adb-cli

adb-cli is a lightweight cli to manage Autonomous Databases in your OCI tenancy

<details>
<summary>Demo</summary>

<img src="docs/adb-cli_demo.gif"/>

</details>
<br/>

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
  * _inspect_
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

## License

Copyright © 2022, 2024 PaoloB

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
