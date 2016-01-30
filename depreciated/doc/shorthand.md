
# shorthand

Shorthand is a simple text substitution program. It is based on a simple key value substitution.  It supports several types of assignments.

+ assigning a string to a LABEL
+ assigning the contents of a file to a LABEL
+ assigning the output of a shell command to a LABEL
+ Evaluating an short hand text string and assigning the results to a LABEL

The *shorthand* program replaces the LABEL with the value assigned to
it where ever it is encountered in the text being passed into the pre-processor.
The assignment statement is not output by the program. Shorthand works on standard input and ouput.

+ text substitutions defined with LABEL := STRING
+ file inclusion defined with LABEL :< PATH TO FILE TO INCLUDE
    + support middle of file extraction negative index refers to lines from end of file
    + middle 6,-10 would mean the buffer size would be ten lines and when you hit eof the buf will be discarded.
    + LABEL :< #,# PATH TO FILE FRAGMENT TO INCLUDE
+ assign the output of a shell command with LABEL :! SHELL_COMMAND
+ assign the output of evaluating a shorthand phrase with LABEL :{ SOME_SHORTHAND_REFERENCED_CONTENT
    + concatinating two labels 
        + @string1 := This is string one
        + @string2 := and this is string two
        + @COMBINED :{ @string1 @string2
    + evaludating a shorthand file and assigning it to @COMBINE
        + @shorthand_file_contents :< myfile.shorthand
        + @COMBINED :{ @shorthand_file_contents


Note the spaces surrounding " := ", " :< ", " :! ", and " :{ " are required.

## Example


In this example a file containing the text of pre-amble is assigned to the
label @PREAMBLE, the time is assigned to the label @NOW by using the shell command *date*.

```text
    @PREAMBLE :< preamble.txt
    @NOW :! date

    At @NOW I read the @PREAMBLE until everyone falls asleep.
```

If the time was "Wed Aug 26 10:33:13 PDT 2015" and the file preamble.txt contained the 
phrase "Hello World" (including the quotes but without any carriage return or line feed) 
the output after processing the shorthand would look like -

```text

    At Wed Aug 26 10:33:13 PDT 2015
    I read the "Hello World" until everyone falls asleep.
```

Notice the lines containing the assignments are not included in the output
and that no carriage returns or line feeds are added the the substituted labels but the assignments do respect any they
may contains (e.g. date includes a line feed so it forces a newline.).

