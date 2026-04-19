# Gator

## Requirements

- Go
- Postgres

## Installation

to install run the command "go install github.com/augustfrih/gator" from the project directory

## Set up

To set up create a .config.gator file in your home directory with the contents:

```
{
  "db_url": "connection_string_goes_here",
}
```


## Usage

- To register a new user "gator register _username_"
- To login as a user "gator login _username_"
- To add a new feed "gator addfeed _name url_"
- To run the aggregator "gator agg _time between reqs_"
- To browse posts "gator browse _(optional) limit_"
