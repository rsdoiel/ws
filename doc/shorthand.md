
# shorthand

Shorthand is a simple text substitution pre-processor. It is based on a simple
key value substitution.  It only supports three operations

+ assigning a string to a LABEL
+ assigning the contents of a file to a LABEL
+ when the pre-processor is run the LABEL is replaced with the value assigned.

The *shorthand* pre-processor replaces the LABEL with the value assigned to
it where ever it is encountered in the text being passed into the pre-processor.
The assignment statement is not output by the preprocessor.

+ text substitutions defined with LABEL := STRING
+ file inclusion defined with LABEL :< PATH TO FILE TO INCLUDE
    + support middle of file extraction negative index refers to lines from end of file
    + middle 6,-10 would mean the buffer size would be ten lines and when you hit eof the buf will be discarded.
    + LABEL :< #,# PATH TO FILE FRAGMENT TO INCLUDE


The spaces surrounding ":=" and ":<" are required.

## Example


In this example a file containing the text of pre-amble is assigned to the
label @PREAMBLE, the time 3:30 is assigned to the label @NOW.

```text
    @PREAMBLE :< /home/me/preamble.text
    @NOW := 3:30

    At @NOW I will be reading the @PREAMBLE until everyone falls asleep.
```

If the file preamble.txt contained the phrase "Hello World" (including
the quotes but without any carriage return or line feed) the output after
processing the shorthand would look like -

```text

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.
```

Notice the lines containing the assignments are not included in the output
and that no carriage returns or line feeds are added the the substituted labels.
