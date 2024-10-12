#!bin/bash 

message="$1"

if [ -z "$message" ]; then 
    echo "The command usage: bash autopush.sh <commit message here>  👽"
    exit 1
fi

git add .
git commit -m "$message"
git push origin "$(git branch --show-current)"