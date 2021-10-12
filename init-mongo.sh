#!/usr/bin/env bash

echo 'Creating application user and db';

mongo godb  --username godb  --password 123456  --authenticationDatabase admin  --host localhost  --port 27017  --eval "db.createUser({user: 'godb', pwd: '123456', roles:[{role:'dbOwner', db: 'godb'}]});"

echo 'User: godb create to database godb';

mongo testDB  --username godb  --password 123456  --authenticationDatabase admin  --host localhost  --port 27017  --eval "db.createUser({user: 'godb', pwd: '123456', roles:[{role:'dbOwner', db: 'testDB'}]});"

echo 'User: godb create to database testDB';
