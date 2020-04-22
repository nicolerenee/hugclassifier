# HUG Classifier

The HUG Classifier is a script that pulls data from the Meetup API and attempts to classify it. The data that is pulled
is the events for a given time range for all of the Meetup groups that belong to the
[HashiCorp User Group Pro Group](https://www.meetup.com/pro/hugs).

We get all of the events and the description for the events and then run them through some basic word matching to
categorize the meetup to determine what was discussed. The keywords are [defined here](pkg/classifier/classify.go#L14).

## System Requirements

This script pulls any required credentials from 1Password using the [op](https://support.1password.com/command-line/)
command line tool. Before you are able to run the script you will need to signin to your 1password vault. The default
account for the app is `hashicorp` but this can be changed with a setting. If you are not signed in you will see an
error similar to what is listed below.

```shell
$ ./hugclassifier generate
2020/04/22 17:14:32 info: using config file: /Users/nicole/.hugclassifier.yaml
2020/04/22 17:14:32 error: 1Password failed, ensure you are logged in by running: eval $(op signin hashicorp)
2020/04/22 17:14:32 error: failed calling 1Password: exit status 1
[LOG] 2020/04/22 17:14:32 (ERROR)  You are not currently signed in. Please run `op signin --help` for instructions
```

To resolve this you just need to login using the following command:

```shell
$ eval $(op signin hashicorp)
Enter the password for nhubbard@hashicorp.com at hashicorp.1password.com:
```

## Configuration

You may create a configuration file named `.hugclassifier.yaml` stored in your home directory. An example of this config
is below:

```yaml
onepassword:
  account: hashicorp
  vault: aaaaaaaaaaaaaaaaaaaaaaaaaa
  login: 11111111111111111111111111
```

The available configuration options are:
| Config Value        | Description |
|---------------------|-------------|
| duration.days       | The number of days to generate a report for. Command line flag: `-d` `--days` Default: `30` |
| duration.end        | The last day to pull data for. Formatted as `2020-02-01`. Command line flag: `--end`        |
| duration.start      | The first day to pull data for. Formatted as `2020-01-01`. Command line flag: `--start`     |
| onepassword.account | The account name of your 1Password account, the subdomain. Default: `hashicorp`             |
| onepassword.login   | The UUID of the login item storing the meetup API credentials. **REQUIRED**                 |
| onepassword.vault   | The UUID of the vault storing the login item. **REQUIRED**                                  |
| output              | The location of the CSV file that the results will be saved to. Default `./results.csv`     |

## Generating the report

To generate the report you will use the `hugclassifier generate` command. This command accepts a couple of arguments to
customize the data that is pulled. (All commands assume the 1Password config options are specified in your .hugclassifier.yaml file).

```shell
# Get the last 30 days worth of events and store in results.csv
$ hugclassifier generate
```

```shell
# Get the last 90 days worth of events and store in results.csv
$ hugclassifier generate --days 90
```

```shell
# Get the events since Jan 2020 and store in results.csv
$ hugclassifier generate --start 2020-01-01
```

```shell
# Get the events in Jan 2020 and store in hug2020-01.csv
$ hugclassifier generate --start 2020-01-01 --end 2020-02-01 --output hug2020-01.csv
```
