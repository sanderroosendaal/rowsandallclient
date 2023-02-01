
# Table of Contents

1.  [Rowsandall Client](#org4689731)
    1.  [Introduction](#org97b9cd3)
    2.  [Run options and configuration](#org5ae5bea)
    3.  [Using the client](#org18b2830)



<a id="org4689731"></a>

# Rowsandall Client


<a id="org97b9cd3"></a>

## Introduction

This is an example client demonstrating the Rowsandall API, see <https://rowsandall.com/rowers/developers/>.

The client demonstrates the v2 and v3 versions of the API.


<a id="org5ae5bea"></a>

## Run options and configuration

You need to ask the developer of [https://rowsandall.com](https://rowsandall.com) for API keys (client ID, client secret), and
put these into a config.yaml file together with your redirect URL.

By default, the program reads the config.yaml file.

Options:

-   v (verbose) Verbose mode, logging information to stdout
-   c (config file) Set a different YaML file as configuration file
-   a (authorized) Use the stored tokens from a tokens.yaml


<a id="org18b2830"></a>

## Using the client

The client runs a web server on localhost:9094 with the following endpoints:

-   [BROKEN LINK: localhost:9094/] Homepage (redirects to authorization if not authorized)
-   [BROKEN LINK: localhost:9094/oauth2] Redirects to the Rowsandall site where you login and exchange the tokens
-   [BROKEN LINK: localhost:9094/workouts] Gets your workouts
-   [BROKEN LINK: localhost:9094/workout] Creates a new workout on Rowsandall, without stroke data (API v2)
-   [BROKEN LINK: localhost:9094/strokedata/{id}] Adds stroke data from teststrokes.json file to the created workout with ID {id}
-   [BROKEN LINK: localhost:9094/strokedatav3] Create Workout and add strokedata from teststrokes2.json (API v3)
-   [BROKEN LINK: localhost:9094/form] Allows you to upload a CSV file according to [https://github.com/sanderroosendaal/rowingdata#csv-file-standard](https://github.com/sanderroosendaal/rowingdata#csv-file-standard)
    and specify workout name and other fields, then uploads it to Rowsandall

