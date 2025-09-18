
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

Now please read https://codingchallenges.fyi/challenges/challenge-wc. 

For each step outlined in that page, write in docs/STEPS.md a prompt for achieving that step. 

Our philosophy is to start at Step 0 with a "walking skeleton" where all the functions are implemented as stubs and the program just returns "0 0 0".  With every additional requirement, we flesh each function more, just enough to meet the current requirements, with no anticipation.

After every stage, it should be possible to build and test manually

The prompt should include:

 - write an appropriate AT, or looking for an existing one and unskipping it
 - verify that the AT fails for the right reason
 - write appropriate unit test for each of the functions in the pipeline
 - verify that they fail for the right reason (compilation errors are not a good reason)
 - write the simplest possible implementation that makes all the tests pass 
 - verify that all the tests pass, including the AT
 - wait for feedback from the user
 
What do you think of this plan?
 

