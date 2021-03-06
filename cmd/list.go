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
		Use:   "list [-e | --enterprise] <enterprise> | list [-o | --organization] <organization>",
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

	listCmd.Flags().StringVarP(&listFlags.Enterprise, "enterprise", "e", "", "enterprise name")
	listCmd.Flags().StringVarP(&listFlags.Organization, "organization", "o", "", "organization name")
}

// listEnterpriseUsers() - lists all users in the enterprise.
func listEnterpriseUsers(enterprise Enterprise) {
	fmt.Printf("Enterprise: %s (%s; %s)\n\n", EnterpriseQuery.Enterprise.Name, enterprise, EnterpriseQuery.Enterprise.ID)
	fmt.Printf("%-40s %-40s\n", "Login", "SAML NameID")
	fmt.Printf("%-40s %-40s\n", "-----", "-----------")
	for name, user := range enterprise.Users {
		fmt.Printf("%-40s %-40s\n", name, user.SamlNameId)
	}
}

// listOrganizationUsers() - lists all users in the organization.
func listOrganizationUsers(organization Organization) {
	fmt.Printf("Organization: %s (%s; %s)\n\n", OrganizationQuery.Organization.Name, organization, OrganizationQuery.Organization.ID)
	fmt.Printf("%-40s %-40s\n", "Login", "SAML NameID")
	fmt.Printf("%-40s %-40s\n", "-----", "-----------")
	for name, user := range organization.Users {
		fmt.Printf("%-40s %-40s\n", name, user.SamlNameId)
	}
}

// runList() - runs the list command.
func runList(args []string) {
	if listFlags.Enterprise != "" && listFlags.Organization != "" {
		panicOnError(fmt.Errorf("only one of '-e' or '-o' is allowed"))
	} else if listFlags.Enterprise != "" {
		panicOnError(validateEnterprise(listFlags.Enterprise))
		topEnterprise.Name = listFlags.Enterprise
		retrieveEnterpriseMembers(topEnterprise.Name)
		retrieveEnterpriseSamlIds(topEnterprise.Name)
		listEnterpriseUsers(topEnterprise)
	} else if listFlags.Organization != "" {
		panicOnError(validateOrganization(listFlags.Organization))
		topOrganization.Name = listFlags.Organization
		retrieveOrganizationUsers(topOrganization.Name)
		listOrganizationUsers(topOrganization)
	}
}
