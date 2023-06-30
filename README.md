# GHPR
Small cli utility (still under construction) to get a list of all PR's that are assigned to you.

## Install
```
go install github.com/jm96441n/ghpr
```

## Usage

The cli assumes you have your github token stored in the env at `GITHUB_ACCESS_TOKEN` you can override that via the cli
args (eventually this will be all set during an initial configuration command). You are also required to pass in your
github handle on each call, this also will eventually be done once via a configuration command.

- using default env
```
ghpr -handle [YOUR HANDLE]
```

- override env
```
ghpr -handle [YOUR HANDLE] -tokenEnv [NAME OF ENV VAR WHERE YOUR GITHUB TOKEN LIVES]
```
