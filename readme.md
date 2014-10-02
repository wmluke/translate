# Translate

A command line utility to translate [Java ResourceBundle Properties Files](http://docs.oracle.com/javase/tutorial/i18n/resbundle/propfile.html) with Google Translate.

## Install

```bash
$ go get github.com/wmluke/translate
```

## Usage

```
NAME:
   translate - translate a Java ResourceBundle Properties file with Google Translate

USAGE:
   translate [options] source_file target_file

VERSION:
   0.1.0

AUTHOR:
  Luke Bunselmeyer

OPTIONS:
   --source, -s 		source langauge code
   --target, -t 		target langauge code
   --key, -k 			Google translate API key [$GOOGLE_API_KEY]
   --help, -h			show help
   --generate-bash-completion
   --version, -v		print the version
```

## License
MIT
