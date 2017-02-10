#!/bin/bash

mongoimport -h mongo -d students -c students --file /dbfixtures/students_fixtures.json
mongoimport -h mongo -d subjects -c subjects --file /dbfixtures/subjects_fixtures.json

