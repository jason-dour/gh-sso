// cmd/list
//
// List command for gh-sso.

package cmd

import (
	"fmt"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
)

type Flags struct {
	Enterprise   string
	Organization string
}

var (
	cmdFlags Flags

	// List command.
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

	listCmd.Flags().StringVarP(&cmdFlags.Enterprise, "enterprise", "e", "", "enterprise name")
	listCmd.Flags().StringVarP(&cmdFlags.Organization, "organization", "o", "", "organization name")
}

// listEnterpriseUsers() - lists all users in the enterprise.
func listEnterpriseUsers(enterprise string) {
	variables := map[string]interface{}{
		"enterpriseName": githubv4.String(enterprise),
		"cursor":         (*githubv4.String)(nil),
	}

	fmt.Printf("%-32s %-48s\n", "Login", "SAML NameID")
	fmt.Printf("%-32s %-48s\n", "-----", "-----------")
	for {
		err := Client.Query(Ctx, &EnterpriseSamlIdpUsersQuery, variables)
		panicOnError(err)

		for _, edge := range EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.Edges {
			fmt.Printf("%-32s %-48s\n", edge.Node.User.Login, edge.Node.SamlIdentity.NameId)
		}
		if !EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}
}

// listOrganizationUsers() - lists all users in the organization.
func listOrganizationUsers(organization string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(organization),
		"cursor":           (*githubv4.String)(nil),
	}

	fmt.Printf("%-32s %-48s\n", "Login", "SAML NameID")
	fmt.Printf("%-32s %-48s\n", "-----", "-----------")
	for {
		err := Client.Query(Ctx, &OrganizationSamlIdpUsersQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.Edges {
			fmt.Printf("%-32s %-48s\n", edge.Node.User.Login, edge.Node.SamlIdentity.NameId)
		}
		if !OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}
}

// validateEnterprise() - validates the enterprise exists.
func validateEnterprise(enterprise string) error {
	variables := map[string]interface{}{
		"enterpriseName": githubv4.String(enterprise),
	}

	err := Client.Query(Ctx, &EnterpriseQuery, variables)
	panicOnError(err)

	if EnterpriseQuery.Enterprise.ID == "" {
		return fmt.Errorf("enterprise '%s' not found", enterprise)
	} else {
		fmt.Printf("Enterprise: %s (%s; %s)\n\n", EnterpriseQuery.Enterprise.Name, enterprise, EnterpriseQuery.Enterprise.ID)
	}

	return nil
}

// validateOrganization() - validates the organization exists.
func validateOrganization(organization string) error {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(organization),
	}

	err := Client.Query(Ctx, &OrganizationQuery, variables)
	panicOnError(err)

	if OrganizationQuery.Organization.ID == "" {
		return fmt.Errorf("organization '%s' not found", organization)
	} else {
		fmt.Printf("Organization: %s (%s; %s)\n\n", OrganizationQuery.Organization.Name, organization, OrganizationQuery.Organization.ID)
	}

	return nil
}

// runList() - runs the list command.
func runList(args []string) {
	if cmdFlags.Enterprise != "" && cmdFlags.Organization != "" {
		panicOnError(fmt.Errorf("only one of '-e' or '-o' is allowed"))
	} else if cmdFlags.Enterprise != "" {
		panicOnError(validateEnterprise(cmdFlags.Enterprise))
		listEnterpriseUsers(cmdFlags.Enterprise)
	} else if cmdFlags.Organization != "" {
		panicOnError(validateOrganization(cmdFlags.Organization))
		listOrganizationUsers(cmdFlags.Organization)
	}
}
