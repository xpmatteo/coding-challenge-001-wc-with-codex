
# Design prompt

## Goal
We want to implement a clone of the Unix wc command, following the steps outlined in
https://codingchallenges.fyi/challenges/challenge-wc

## Design

We would like a sort of pipeline implementation.  Here is how it would look like:

1. ParseArgs(args []string) (wc.Config, error)
2. AnalyzeFiles(wc.Config) ([]wc.Stats, error)
3. AddTotal(wc.Config, []wc.File) ([]wc.Stats, error) // might not do anything if only one file is analyzed
4. Format(wc.Config, []wc.Stats) ([]string, error)

The input is the args passed to the command line
The output is a slice of strings that can be printed on stdout

The custom types is:

type Stats struct {
	// filename
	// count of lines, chars, words, etc
}

What do you think of this proposed design?

(Then clarified that each stage is a function)

## Steps