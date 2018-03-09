package useradm

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var RootCmd = &cobra.Command{
	Use: "useradm",
	Short: "User administration tool for the login scaffold and based upon projects",
	Long: "Manages (create, password reset, delete) Users for the login scaffold and based upon projects",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


//Setup Flags Here
func init(){

}