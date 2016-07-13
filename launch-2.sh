#!/bin/bash

clear

DIR=`echo ${PWD##*/}`
echo "Stopping ${DIR} With Supervisor"
sudo supervisorctl stop $DIR

read -p "Are you sure you want to delete the old files?
Be sure you're SSH'd onto the server before running this script!
Press y[Y] to continue..." -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]] then
    echo "Removing Old Files"
    sudo rm -rf ${DIR} assets/ templates/ $DIR.tar
else
    echo "Aborted Launch"
    exit 1
fi

echo "Moving New Tar to Current Location"
sudo mv ~/$DIR.tar .

if [ ! -f "${DIR}.tar" ]; then
    echo "Moving ${DIR}.tar Failed"
    exit 1
fi

echo "Extraction Contents of New Tar"
sudo tar xf $DIR.tar

if [ ! -f $DIR ]; then
    echo "No Binary Found"
    exit 1
fi

echo "Starting ${DIR} With Supervisor"
sudo supervisorctl start $DIR
