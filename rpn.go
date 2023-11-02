package main

var manpage = `
NAME
    rpn - Reverse Polish Notation Calculator for the commandline

SYNOPSIS
        Usage: rpn [-bdvh] [<operator>]
    
        Options:
          -b, --batchmode   enable batch mode
          -d, --debug       enable debug mode
          -v, --version     show version
          -h, --help        show help
    
        When <operator>  is given, batch  mode ist automatically  enabled. Use
        this only when working with stdin. E.g.: echo "2 3 4 5" | rpn +

DESCRIPTION
    rpn is a command line calculator using reverse polish notation.

  Working principle
    Reverse Polish Notation (short: RPN) requires to have a stack where
    numbers and results are being put. So, you put numbers onto the stack
    and each math operation uses these for calculation, removes them and
    puts the result back.

    To visualize it, let's look at a calculation:

        ((80 + 20) / 2) * 4

    This is how you enter the formula int an RPN calculator and how the
    stack evolves during the operation:

        | rpn commands | stack contents | calculation   |
        |--------------|----------------|---------------|
        |           80 |             80 |               |
        |           20 |          80 20 |               |
        |            + |            100 | 80 + 20 = 100 |
        |            2 |          100 2 |               |
        |            / |             50 | 100 / 2 = 50  |
        |            4 |           50 4 |               |
        |            x |            200 | 50 * 4 = 200  |

    The last stack element 200 is the calculation result.

  USAGE
    The default mode of operation is the interactive mode. You'll get a
    prompt which shows you the current size of the stack. At the prompt you
    enter numbers followed by operators or mathematical functions. You can
    use completion for the functions. You can either enter each number or
    operator on its own line or separated by whitespace, that doesn't
    matter. After a calculation the result will be immediately displayed
    (and added to the stack). You can quit interactive mode using the
    commands quit or exit or hit one of the "ctrl-d" or "ctrl-c" key
    combinations.

    If you feed data to standard input (STDIN), rpn just does the
    calculation denoted in the contet fed in via stdin, prints the result
    and exits. You can also specify a calculation on the commandline.

    Here are the three variants ($ is the shell prompt):

        $ rpn
        rpn> 2
        rpn> 2
        rpn> +
        = 4
    
        $ rpn
        rpn> 2 2 +
        = 4
    
        $ echo 2 2 + | rpn
        4
    
        $ rpn 2 2 +
        4

    The rpn calculator provides a batch mode which you can use to do math
    operations on many numbers. Batch mode can be enabled using the
    commandline option "-b" or toggled using the interactive command batch.
    Not all math operations and functions work in batch mode though.

    Example of batch mode usage:

        $ rpn -b
        rpn->batch > 2 2 2 2 +
        = 8
    
        $ rpn
        rpn> batch
        rpn->batch> 2 2 2 2 +
        8
    
        $ echo 2 2 2 2 + | rpn -b
        8
    
        $ echo 2 2 2 2 | rpn +
        8

    If the first parameter to rpn is a math operator or function, batch mode
    is enabled automatically, see last example.

  STACK MANIPULATION
    There are lots of stack manipulation commands provided. The most
    important one is undo which goes back to the stack before the last math
    operation.

    You can use dump to display the stack. If debugging is enabled ("-d"
    switch or debug toggle command), then the backup stack is also being
    displayed.

    The stack can be reversed using the reverse command.

    You can use the shift command to remove the last number from the stack.

  BUILTIN OPERATORS AND FUNCTIONS
    Basic operators: + - x /

    Math functions:

        sqrt                 square root
        mod                  remainder of division (alias: remainder)
        max                  batch mode only: max of all values
        min                  batch mode only: min of all values
        mean                 batch mode only: mean of all values (alias: avg)
        median               batch mode only: median of all values
        %                    percent
        %-                   substract percent
        %+                   add percent
        ^                    power

EXTENDING RPN USING LUA
    You can use a lua script with lua functions to extend the calculator. By
    default the tool looks for "~/.rpn.lua". You can also specify a script
    using the <kbd>-c</kbd> flag.

    Here's an example of such a script:

        function add(a,b)
          return a + b
        end
    
        function init()
          register("add", 2, "addition")
        end

    Here we created a function "add()" which adds two parameters. All
    parameters are "FLOAT64" numbers. You don't have to worry about stack
    management, this is taken care of automatically.

    The function "init()" MUST be defined, it will be called on startup. You
    can do anything you like in there, but you need to call the "register()"
    function to register your functions to the calculator. This function
    takes these parameters:

    *   function name

    *   number of arguments expected (1,2 or -1 allowed), -1 means batch
        mode.

    *   help text

    Please refer to the lua language reference:
    <https://www.lua.org/manual/5.4/> for more details about LUA.

    Please note, that io, networking and system stuff is not allowed though.
    So you can't open files, execute other programs or open a connection to
    the outside!

GETTING HELP
    In interactive mode you can enter the help command (or ?) to get a short
    help along with a list of all supported operators and functions.

    To read the manual you can use the manual command in interactive mode.
    The commandline option "-m" does the same thing.

    If you have installed rpn as a package or using the distributed tarball,
    there will also be a manual page you can read using "man rpn".

BUGS
    In order to report a bug, unexpected behavior, feature requests or to
    submit a patch, please open an issue on github:
    <https://github.com/TLINDEN/rpnc/issues>.

LICENSE
    This software is licensed under the GNU GENERAL PUBLIC LICENSE version
    3.

    Copyright (c) 2023 by Thomas von Dein

    This software uses the following GO modules:

    readline (github.com/chzyer/readline)
        Released under the MIT License, Copyright (c) 2016-2023 ChenYe

    pflag (https://github.com/spf13/pflag)
        Released under the BSD 3 license, Copyright 2013-2023 Steve Francia

    gopher-lua (github.com/yuin/gopher-lua)
        Released under the MIT License, Copyright (c) 2015-2023 Yusuke
        Inuzuka

AUTHORS
    Thomas von Dein tom AT vondein DOT org

`
var usage = `

Usage: rpn [-bdvh] [<operator>]

Options:
  -b, --batchmode   enable batch mode
  -d, --debug       enable debug mode
  -v, --version     show version
  -h, --help        show help

When <operator>  is given, batch  mode ist automatically  enabled. Use
this only when working with stdin. E.g.: echo "2 3 4 5" | rpn +

`
