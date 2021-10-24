.PHONY: default clean

default: report.pdf

clean:
	git clean -fX .
	rm report.pdf

%.pdf: %.tex
	lualatex $*.tex
	biber $*
	lualatex $*.tex
	lualatex $*.tex
