// cmd/list
//
// List command for gh-sso.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listFlags ListFlags

	listCmd = &cobra.Command{
		Use:   "list ( (-e | --enterprise) | (-o | --organization) ) <ent_or_org_name>...",
		Short: "List SSO Users",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runList(args)
		},
	}
)

// Initialization.
func init() {
	// Add the list command to the root command.
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&listFlags.Enterprises, "enterprise", "e", []string{}, "enterprise name(s)")
	listCmd.Flags().StringSliceVarP(&listFlags.Organizations, "organization", "o", []string{}, "organization name(s)")
}

// listEnterpriseUsers() - lists all users in the enterprise.
func listEnterpriseUsers(enterprise Enterprise) {
	fmt.Printf("Enterprise: %s (%s; %s)\n\n", enterprise.Name, enterprise.Login, enterprise.ID)
	fmt.Printf("%-40s %-40s\n", "Login", "SAML NameID")
	fmt.Printf("%-40s %-40s\n", "-----", "-----------")
	for name, user := range enterprise.Users {
		fmt.Printf("%-40s %-40s\n", name, user.SamlNameId)
	}
}

// listOrganizationUsers() - lists all users in the organization.
func listOrganizationUsers(organization Organization) {
	fmt.Printf("Organization: %s (%s; %s)\n\n", organization.Name, organization.Login, organization.ID)
	fmt.Printf("%-40s %-40s\n", "Login", "SAML NameID")
	fmt.Printf("%-40s %-40s\n", "-----", "-----------")
	for name, user := range organization.Users {
		fmt.Printf("%-40s %-40s\n", name, user.SamlNameId)
	}
}

// runList() - runs the list command.
func runList(args []string) {
	Cfg.Enterprises = append(Cfg.Enterprises, listFlags.Enterprises...)
	Cfg.Organizations = append(Cfg.Organizations, listFlags.Organizations...)

	fmt.Printf("enterprises:\n")
	for _, v := range Cfg.Enterprises {
		fmt.Printf("  - %s:\n", v)
		panicOnError(validateEnterprise(v))
		topEnterprise.Name = v
		retrieveEnterpriseMembers(topEnterprise.Name)
		retrieveEnterpriseSamlIds(topEnterprise.Name)
		listEnterpriseUsers(topEnterprise)
	}

	fmt.Printf("organizations:\n")
	for _, v := range Cfg.Organizations {
		fmt.Printf("  - %s\n", v)
	}

	// if listFlags.Enterprise != "" && listFlags.Organization != "" {
	// 	panicOnError(fmt.Errorf("only one of '-e' or '-o' is allowed"))
	// } else if listFlags.Enterprise != "" {
	// 	panicOnError(validateEnterprise(listFlags.Enterprise))
	// 	topEnterprise.Name = listFlags.Enterprise
	// 	retrieveEnterpriseMembers(topEnterprise.Name)
	// 	retrieveEnterpriseSamlIds(topEnterprise.Name)
	// 	listEnterpriseUsers(topEnterprise)
	// } else if listFlags.Organization != "" {
	// 	panicOnError(validateOrganization(listFlags.Organization))
	// 	topOrganization.Name = listFlags.Organization
	// 	retrieveOrganizationUsers(topOrganization.Name)
	// 	listOrganizationUsers(topOrganization)
	// } else {

	// 	panicOnError(fmt.Errorf("either '-e' or '-o' flag must be used, or enterprises/organizations used in config file"))
	// }
}
