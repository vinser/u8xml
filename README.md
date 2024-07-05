## u8xml Go package

### _Reinvent the wheel!_

The __u8xml__ package implements `NewDecoder` which can be used to parse XML files with IANA character encodings such as Windows-1252, ISO-8859-1, etc. It can be used to decode XML files/strings with Go Standard Library xml package Decoder type methods like Decode(), Token(), etc.  

XML files must contain a BOM at the beginning in the case of unicode characters or an XML declaration with an encoding attribute otherwise.  

XML files with UTF-8 content may be detected either by BOM or XML declaration. XML files with no BOM or XML declaration will be treated as UTF-8.  

The package also implements functions `NewReader,` which creates io.Reader that converts content to UTF-8, and `DetectEncoding`, which can detect the IANA encoding.

### u8hex CLI utility
The `cmd` folder contains the source code of the `u8hex` command-line interface utility, which may be used to get the hex representation of a string with a given character set. It may be useful for debugging.

### Credits
__u8hex__ is inspired by [__cpd__](https://github.com/softlandia/cpd)


