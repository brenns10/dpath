package main

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"os"
)

/*
A command-line driver for evaluating DPath expressions.
*/
func main() {
	var parseTreeBuf bytes.Buffer
	var r bool

	if len(os.Args) < 2 {
		log.Fatal("Must provide a DPath expression.")
	}

	// Parse the DPath expression.
	tree, err := ParseString(os.Args[1])
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Syntax error.")
	}

	// Log the parse tree.
	parseTreeBuf.WriteString("Parse Tree:\n")
	tree.Print(&parseTreeBuf, 0)
	log.WithFields(log.Fields{
		"tree": parseTreeBuf.String(),
	}).Debug("Created parse tree.")

	// Evaluate the expression and print the results.
	ctx := DefaultContext()
	seq, err := tree.Evaluate(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error while evaluating expression.")
	}

	for r, err = seq.Next(ctx); r && err == nil; r, err = seq.Next(ctx) {
		err = seq.Value().Print(os.Stdout)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error while iterating.")
	}
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}
