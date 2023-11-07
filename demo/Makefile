#


all:
	asciinema rec --cols 80 --row 25 -c "env - PS1='> ' PATH=..:$PATH /bin/bash --norc --noprofile" --overwrite demo.cast
	agg demo.cast demo.gif

