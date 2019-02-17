## Reverse Polish Notation Calculator for the commandline

This is a small commandline calculator which takes its input in
[https://en.wikipedia.org/wiki/Reverse_Polish_notation](reverse polish notation)
form.

It  has  an  unlimited  stack,  supports  various  stack  manipulation
commands, can be used interactively or  via a pipe and has a collector
mode. It doesn't have any other dependencies than Perl.

## Usage

Calculate the summary  resistance of parallel resistors  with 220, 330
and 440 Ohm using the following formula:

    1 / (1/R1 + 1/R2 + 1/R3)

Here's the sample session:

     0 % 1
    stack    1: 1
    
     1 % 1
    stack    2: 1
    stack    1: 1
    
     2 % 220
    stack    3: 1
    stack    2: 1
    stack    1: 220
    
     3 % /
    stack    2: 1
    stack    1: 0.00454545454545455
    
    => 0.00454545454545455
    
     2 % 1
    stack    3: 1
    stack    2: 0.00454545454545455
    stack    1: 1
    
     3 % 330
    stack    4: 1
    stack    3: 0.00454545454545455
    stack    2: 1
    stack    1: 330
    
     4 % /
    stack    3: 1
    stack    2: 0.00454545454545455
    stack    1: 0.00303030303030303
    
    => 0.00303030303030303
    
     3 % 1
    stack    4: 1
    stack    3: 0.00454545454545455
    stack    2: 0.00303030303030303
    stack    1: 1
    
     4 % 440
    stack    5: 1
    stack    4: 0.00454545454545455
    stack    3: 0.00303030303030303
    stack    2: 1
    stack    1: 440
    
     5 % /
    stack    4: 1
    stack    3: 0.00454545454545455
    stack    2: 0.00303030303030303
    stack    1: 0.00227272727272727
    
    => 0.00227272727272727
    
     4 % +
    stack    3: 1
    stack    2: 0.00454545454545455
    stack    1: 0.0053030303030303
    
    => 0.0053030303030303
    
     3 % +
    stack    2: 1
    stack    1: 0.00984848484848485
    
    => 0.00984848484848485
    
     2 % /
    stack    1: 101.538461538462
    
    => 101.538461538462

The *%* character denotes the interactive prompt. What we basically entered was:

    1 1 220 / 1 330 / 1 440 / + + /
    
Which translates to:

    1 ((1 / 220) + (1 / 330) + (1 / 440))
    
So,  you're entering  the numbers  and operators  as you  would do  on
paper. To learn more, refer to the Wikipedia page linked above.

## Collector mode

Beside  traditional RPN  you can  also  enter a  special mode,  called
*collector mode* by entering  the <kbd>(</kbd> command.  The collector
mode  has its  own stack  (a  sub stack)  which is  independed of  the
primary stack.  Inside this  mode you can  use all  operators, however
they work on *ALL* items on the sub stack.

So, let's compare. If you had in normal RPN mode the following stack:

    3
    5
    6
    
and then entering the <kbd>+</kbd>  operator, the calculator would pop
5 and 6  from the stack, add them  and push the result 11  back to the
stack.

However, if  you are in collector  mode with this stack,  then all the
items would be added, the sub stack would be cleared and the result 14
would be added to the primary stack.

You  will  leave  the  collector  mode  after  an  operator  has  been
executed. But  you can  also just  leave the  collector mode  with the
command  <kbd>)</kbd> leaving  the sub  stack intact.   That is,  upon
re-entering collector mode at a  later time, you'll find the unaltered
sub stack of before.

## Undo

Every operation which  modifies the stack can be  reversed by entering
the <kbd>u</kbd> command. There's only one level of undo and no redo.

## Using STDIN via a PIPE

If the commandline  includes any operator, commands will  be read from
STDIN, the result will be printed  to STDOUT wihout any decoration and
the  program will  exit. Commands  can be  separated by  whitespace or
newline.

Examples:

    echo "2 2" | rpnc +
    (echo 2; echo 2) | rpnc +
    
Both commands will print 4 to STDOUT.

## Complete list of all supported commands:

* <kbd>c</kbd>    clear stack
* <kbd>s</kbd>    show the stack
* <kbd>d</kbd>    toggle debugging (current setting: 0)
* <kbd>r</kbd>    reverse the stack
* <kbd>R</kbd>    rotate the stack
* <kbd>(</kbd>    enter collect mode
* <kbd>)</kbd>    leave collect mode
* <kbd>u</kbd>    undo last operation
* <kbd>q</kbd>    finish (<kbd>C-d</kbd> works as well)
* <kbd>?</kbd>    print help

## Supported mathematical operators:

* <kbd>+</kbd>     add
* <kbd>-</kbd>     substract
* <kbd>/</kbd>     divide
* <kbd>*</kbd>     multiply
* <kbd>^</kbd>     expotentiate
* <kbd>%</kbd>     percent
* <kbd>&</kbd>     bitwise AND
* <kbd>|</kbd>     bitwise OR
* <kbd>x</kbd>     bitwise XOR
* <kbd>V</kbd>     pull root (2nd if stack==1)

## Copyleft

Copyleft (L) 2019 - Thomas von Dein.
Licensed under the terms of the GPL 3.0.
