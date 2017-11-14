#!/bin/bash
# Check prerequisties

main() {
	tk_check "mysql"
	tk_check "unzip"
	import_do
}

tk_check() {
	command -v $1 >/dev/null 2>&1 || { echo >&2 "This script requires '$1', please install it to do import.  Aborting."; exit 1; }
}

import_do() {
	source "../src/.env"

	# Default DB name
	local DB_NAME="foodcourt"
	local DIR_PUBLIC="../src/public"
	local DIR_PHOTOS="$DIR_PUBLIC/photos"

	echo "* Importing dishes list ..."

	# Import data
	cat ./data.sql | mysql --user="root" --password="wh3r315myc00k135" --port="3366" --host="127.0.0.1" --database="foodcourt"

	if [ ! -d "$DIR_PUBLIC" ]; then
		echo "* Creating 'public' directory ..."
		mkdir "$DIR_PUBLIC"
	fi

	if [ -d "$DIR_PHOTOS" ]; then
		echo "* Deleting existing photos ..."
		rm -rf "$DIR_PHOTOS"
	fi


	echo "* Extracting photos ..."
	unzip photos.zip -d "$DIR_PUBLIC"

	echo "===== IMPORT FINISH ======"
}

main
