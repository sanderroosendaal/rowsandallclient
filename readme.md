
# Table of Contents

1.  [Rowsandall Client](#org92125ec)
    1.  [Introduction](#org0099ae0)
    2.  [Run options and configuration](#org903a73a)
    3.  [Using the client](#org64a45b9)



<a id="org92125ec"></a>

# Rowsandall Client


<a id="org0099ae0"></a>

## Introduction

This is an example client demonstrating the Rowsandall API, see <https://rowsandall.com/rowers/developers/>.

The client demonstrates the v2 and v3 versions of the API.


<a id="org903a73a"></a>

## Run options and configuration

You need to ask the developer of [https://rowsandall.com](https://rowsandall.com) for API keys (client ID, client secret), and
put these into a config.yaml file together with your redirect URL.

By default, the program reads the config.yaml file.

Options:

-   v (verbose) Verbose mode, logging information to stdout
-   c (config file) Set a different YaML file as configuration file
-   a (authorized) Use the stored tokens from a tokens.yaml


<a id="org64a45b9"></a>

## Using the client

The client runs a web server on localhost:9094 with the following endpoints:

-   Homepage (redirects to authorization if not authorized)
-   Redirects to the Rowsandall site where you login and exchange the tokens
-   Gets your workouts
-   Creates a new workout on Rowsandall, without stroke data (API v2)
-   Adds stroke data from teststrokes.json file to the created workout with ID {id}
-   Create Workout and add strokedata from teststrokes2.json (API v3)
-   Allows you to upload a CSV file according to [https://github.com/sanderroosendaal/rowingdata#csv-file-standard](https://github.com/sanderroosendaal/rowingdata#csv-file-standard)
    and specify workout name and other fields, then uploads it to Rowsandall

