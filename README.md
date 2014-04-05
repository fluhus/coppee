Coppee (0.4.0)
==============
An automatic file copier.
Programatically define file name templates and targets, and Coppee will copy
them automatically.

Please feel free to report bugs, comments and requests to
**doctor_troll at walla dot com**

How to Compile
--------------
* Make sure you have a **go** compiler.
* Set **GOPATH** to be the base folder of the project.
* Run `go install coppee`. Don't worry, it doesn't actually install anything on your system.
  It only compiles the code.
* Executable will be in folder **bin**.

Usage
-----
Executing `coppee` will print an informative usage explanation.  
For basic use execute `coppee <directory>`.
The target directory must contain a file named **.coppee**. In this file are the
copying instructions (see dedicated section).

Instruction File Format
-----------------------
### General structure
The instruction file contains the directives for Coppee's actions. It should be
encoded in UTF-8 (Notepad++ makes it easy), and supports all languages. A
sample instruction file can be
found with the main code files. The instruction file should be formatted as follows:  
```
template1
target1
template2
target2
...
```
Each template is a regular expression, followed by its specific target.
File names that match the template, will be copied and named according to
the target. The regular expressions should match the syntax specified
[here](http://code.google.com/p/re2/wiki/Syntax).

### Capturing groups
You can refer to
the source's name in the target, using capturing groups:
```
lecture_(\d+)\.ppt
lesson_${1}.ppt
```
The expression `${1}` will be replaced by the expression matched by the first
parenthesized group. In the same way, `${2}` will be replaced by the second
parenthesized expression, etc. In the above example, a file named **lecture_13.ppt** will
be copied to a file named **lesson_13.ppt**.

### Comments and empty lines
Comments are lines that start with `//`. Comments and empty lines are ignored by
the parser.
```
// This is a comment.
  // Enveloping whitespaces are trimmed for comments, templates and targets.

// Match .doc and .docx files
(.+)\.(docx?)
// Target is the same file with .coppee added before the suffix
${1}.coppee.${2}
```

### Nagated templates
Coppee can also copy all files that *don't* match a certain template.
To do that, add `!` before the template. Spaces are allowed between the `!`
and the template itself, for readability.

#### Capturing groups
Ok, so now you're probably thinking "how am I gonna use capturing groups on a
mismatch?" Lucky for you, we got that covered:  
${0} will be replaced by the entire file name, including directories.  
${1} will be replaced by the directories, including the last separator.  
${2} will be replaced by the file's prefix, up to the last period.  
${3} will be replaced by the file's suffix, including the period.

#### Examples
```
// Copy all except .log files, preserve paths (relative to input directory)
! .*\.log
C:\backup\${0}

// Add '_copy' to the prefix of files that don't begin with a 'd'
// a.txt will be copied to a_copy.txt
!d.*
${2}_copy{3}

// Match files that begin with '!' (with regular capturing)
\!.*
important\${0}
```

Future Features
---------------
* Option to choose a different instruction file.
* Target for files that didn't match any regex.
* Choose whether or not to ignore i/o errors (right now it exits on error).
* *Features suggested by users.*

Tasks
-----
* Write tests.
* Handle subdirectory recursion.

Version History
---------------
### 0.4.0
* Added negated templates.

### 0.3.2
* Split code into packages.
* Added line numbers to errors.

### 0.3.1
* Changed order of command line arguments.
* Fixed a bug that folder were attempted to be copied like files.

### 0.2.1
* Added BOM tolerance.

### 0.2
* Added overwrite mode.
* Added quiet mode.
* Added pretend mode.

### 0.1.1
* Fixed a bug that caused mismatches when char-13 was present in the instruction file.

### 0.1
* First functional version. Should be stable, though not fully functional.


