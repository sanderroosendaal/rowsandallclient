
# Table of Contents

1.  [Rowsandall Client](#org1008df2)
    1.  [Introduction](#org0704c72)
    2.  [Run options and configuration](#orgc690bb2)
    3.  [Using the client](#orge7cd2d0)



<a id="org1008df2"></a>

# Rowsandall Client


<a id="org0704c72"></a>

## Introduction

This is an example client demonstrating the Rowsandall API, see <https://rowsandall.com/rowers/developers/>.

The client demonstrates the v2 and v3 versions of the API.


<a id="orgc690bb2"></a>

## Run options and configuration

You need to ask the developer of [https://rowsandall.com](https://rowsandall.com) for API keys (client ID, client secret), and
put these into a config.yaml file together with your redirect URL.

By default, the program reads the config.yaml file.

Options:

-   v (verbose) Verbose mode, logging information to stdout
-   c (config file) Set a different YaML file as configuration file
-   a (authorized) Use the stored tokens from a tokens.yaml


<a id="orge7cd2d0"></a>

## Using the client

The client runs a web server on localhost:9094 with the following endpoints:

-   localhost:9094/ Homepage (redirects to authorization if not authorized)
-   localhost:9094/oauth2 Redirects to the Rowsandall site where you login and exchange the tokens
-   localhost:9094/workouts Gets your workouts
-   localhost:9094/workout Creates a new workout on Rowsandall, without stroke data (API v2)
-   localhost:9094/strokedata/{id} Adds stroke data from teststrokes.json file to the created workout with ID {id}
-   localhost:9094/strokedatav3 Create Workout and add strokedata from teststrokes2.json (API v3)
-   localhost:9094/form Allows you to upload a CSV file according to [https://github.com/sanderroosendaal/rowingdata#csv-file-standard](https://github.com/sanderroosendaal/rowingdata#csv-file-standard)
    and specify workout name and other fields, then uploads it to Rowsandall

