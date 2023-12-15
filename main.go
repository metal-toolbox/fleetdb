// package main is the entry point
package main

//go:generate sqlboiler crdb --add-soft-deletes

import "go.hollow.sh/fleetdb/cmd"

func main() {
	cmd.Execute()
}
