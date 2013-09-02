go-redirect-checker
===================

A 301 redirect checker written in go

Usage
=====

Create a CSV file called 301s.csv, with the fields originalUrl, expectedUrl. Then,

./redirectChecker

Output
======

Prints out the results of each URL, whether the redirect occured as expected or if there was a loop
