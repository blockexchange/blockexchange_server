
%.png: %.dot
	cat $< | docker run --rm -i vladgolubev/dot2png > $@

all: $(patsubst %.dot,%.png,$(wildcard *.dot))
