## Reverse Polish Notation Calculator for the commandline

This is a small commandline calculator which takes its input in
[reverse polish notation](https://en.wikipedia.org/wiki/Reverse_Polish_notation)
form.

Features:

- unlimited stack
- undo
- various stack manipulation commands
- basic math operators
- advanced math functions (not yet complete)
- provides interactive repl
- can be used on the commandline
- can calculate data in batch mode (also from STDIN)
- extensible with custom LUA functions

## Usage

Basically   you  enter   numbers  followed   by  an   operator  or   a
function. Each number you enter will be put into the stack. Say you
entered two numbers, 2 and 4. If you now enter the `+` operator, those
two numbers will be removed from the stack, added and the result will
be put back onto the stack.

Here's a  comprehensive example:  calculate the summary  resistance of
parallel  resistors with  220, 330  and  440 Ohm  using the  following
formula:

    1 / (1/R1 + 1/R2 + 1/R3)

Here's the sample session:

```
rpn [0]» 1
rpn [1]» 1 220 /
= 0.004545454545454545
rpn [2]» 1 330 /
= 0.0030303030303030303
rpn [3]» 1 440 /
= 0.0022727272727272726
rpn [4]» +
= 0.0053030303030303025
rpn [3]» +
= 0.009848484848484848
rpn [2]» /
= 101.53846153846155
```

It doesn't matter  wether you enter numbers  and operators/function on
the same line or separated by whitespace:

```
rpn [0]» 1 1 220 / 1 330 / 1 440 / + + /
= 0.004545454545454545
= 0.0030303030303030303
= 0.0022727272727272726
= 0.0053030303030303025
= 0.009848484848484848
= 101.53846153846155
```

Works on the commandline as well:

```
rpn 1 1 220 / 1 330 / 1 440 / + + /
0.004545454545454545
0.0030303030303030303
0.0022727272727272726
0.0053030303030303025
0.009848484848484848
101.53846153846155
```

And via STDIN:
```
echo "1 1 220 / 1 330 / 1 440 / + + /" | rpn
0.004545454545454545
0.0030303030303030303
0.0022727272727272726
0.0053030303030303025
0.009848484848484848
101.53846153846155
```

What we basically entered was:

    1 1 220 / 1 330 / 1 440 / + + /
    
Which translates to:

    1 ((1 / 220) + (1 / 330) + (1 / 440))
    
So,  you're entering  the numbers  and operators  as you  would do  on
paper. To learn more, refer to the Wikipedia page linked above.

## Batch mode

Beside  traditional RPN  you can  also  enter a  special mode,  called
*batch  mode* either  by entering  the  `batch` command  or using  the
commandline  switch <kbd>-b</kbd>.

Most  operators and  functions can  be used  with batch  mode but  not
all. In this mode the calculation works on all numbers on the stack so
far.

So, let's compare. If you had in normal RPN mode the following stack:

    3
    5
    6
    
and then enter  the <kbd>+</kbd> operator, the calculator  would pop 5
and 6  from the stack,  add them  and push the  result 11 back  to the
stack.

However, if  you are in  batch mode, then  all the items  would be
added, the sub stack would be cleared and the result 14 would be added
to the stack.

To leave batch mode just enter the `batch` command again (this is a
toggle).

Here's an example using a math function:

    echo 1 2 3 4 5 6 7 | rpn -b median
    4

Really simple.

## Undo

Every operation which  modifies the stack can be  reversed by entering
the `undo` command. There's only one level of undo and no redo.

## Extend the calculator with LUA functions

You can use a lua script with lua functions to extend the
calculator. By default the tool looks for `~/.rpn.lua`. You can also
specify a script using the <kbd>-c</kbd> flag.

Here's an example of such a script:

```lua
function add(a,b)
  return a + b
end

function init()
  register("add", 2, "addition")
end
```

Here  we created  a function  `add()` which  adds two  parameters. All
parameters are `FLOAT64` numbers. You  don't have to worry about stack
management, this is taken care of automatically.

The function `init()` **MUST** be defined, it will be called on
startup. You can do anything you like in there, but you need to call
the `register()` function to register your functions to the
calculator. This function takes these parameters:

- function name
- number of arguments expected (1,2 or -1 allowed), -1 means batch
  mode
- help text

Please [refer to the lua language
reference](https://www.lua.org/manual/5.4/) for more details about
LUA.

**Please note, that io, networking and system stuff is not allowed
though. So you can't open files, execute other programs or open a
connection to the outside!**

## Documentation

The  documentation  is  provided  as  a unix  man-page.   It  will  be
automatically installed if  you install from source.  However, you can
read the man-page online:

https://github.com/TLINDEN/rpnc/blob/main/rpn.pod

Or if you cloned  the repository you can read it  this way (perl needs
to be installed though): `perldoc rpn.pod`.

If you have the binary installed, you  can also read the man page with
this command:

    rpn --man

## Getting help

Although I'm happy to hear from rpn users in private email, that's the
best way for me to forget to do something.

In order to report a bug,  unexpected behavior, feature requests or to
submit    a    patch,    please    open   an    issue    on    github:
https://github.com/TLINDEN/rpnc/issues.

## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://github.com/TLINDEN/rpnc
