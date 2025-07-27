default:
	make -f makefile.wsl test

win:
	make -f makefile.win test

.PHONY: default win
