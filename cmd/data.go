// cmd/data
//
// Data gathering functions.

package cmd

import (
	"fmt"

	"github.com/shurcooL/githubv4"
)

var (
	topEnterprise   Enterprise
	topOrganization Organization
)

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

func retrieveEnterpriseOrgs(enterprise string) {
	variables := map[string]interface{}{
		"enterpriseName": githubv4.String(enterprise),
		"cursor1":        (*githubv4.String)(nil),
	}

	if topEnterprise.Organizations == nil {
		topEnterprise.Organizations = map[string]Organization{}
	}
	for {
		err := Client.Query(Ctx, &EnterpriseOrganizationsQuery, variables)
		panicOnError(err)

		for _, edge := range EnterpriseOrganizationsQuery.Enterprise.Organizations.Edges {
			topEnterprise.Organizations[string(edge.Node.Login)] = Organization{
				Name:         string(edge.Node.Login),
				Repositories: map[string]Repository{},
				Users:        map[string]User{},
			}
		}

		if !EnterpriseOrganizationsQuery.Enterprise.Organizations.PageInfo.HasNextPage {
			break
		}
		variables["cursor1"] = githubv4.String(EnterpriseOrganizationsQuery.Enterprise.Organizations.PageInfo.EndCursor)
	}
}

func retrieveEnterpriseRepos(enterprise string, org string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(org),
		"cursor":           (*githubv4.String)(nil),
	}

	if topEnterprise.Organizations[org].Repositories == nil {
		orgtemp := topEnterprise.Organizations[org]
		orgtemp.Repositories = map[string]Repository{}
	}
	for {
		err := Client.Query(Ctx, &OrganizationRepositoriesQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationRepositoriesQuery.Organization.Repositories.Edges {
			topEnterprise.Organizations[org].Repositories[string(edge.Node.Name)] = Repository{
				Name: string(edge.Node.Name),
			}
		}

		if !OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.EndCursor)
	}
}

func retrieveOrganizationRepos(organization string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(organization),
		"cursor":           (*githubv4.String)(nil),
	}

	if topOrganization.Repositories == nil {
		topOrganization.Repositories = map[string]Repository{}
	}
	for {
		err := Client.Query(Ctx, &OrganizationRepositoriesQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationRepositoriesQuery.Organization.Repositories.Edges {
			topOrganization.Repositories[string(edge.Node.Name)] = Repository{
				Name: string(edge.Node.Name),
			}
		}

		if !OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.EndCursor)
	}
}

func retrieveEnterpriseMembers(enterprise string) {
	variables := map[string]interface{}{
		"enterpriseName": githubv4.String(enterprise),
		"cursor":         (*githubv4.String)(nil),
	}

	if topEnterprise.Users == nil {
		topEnterprise.Users = map[string]User{}
	}

	for {
		err := Client.Query(Ctx, &EnterpriseMembersQuery, variables)
		panicOnError(err)

		for _, edge := range EnterpriseMembersQuery.Enterprise.Members.Edges {
			if _, ok := topEnterprise.Users[string(edge.Node.EnterpriseUserAccount.Login)]; !ok {
				topEnterprise.Users[string(edge.Node.EnterpriseUserAccount.Login)] = User{
					Name: string(edge.Node.EnterpriseUserAccount.Name),
				}
			} else {
				tmpUser := topEnterprise.Users[string(edge.Node.EnterpriseUserAccount.Login)]
				tmpUser.Name = string(edge.Node.EnterpriseUserAccount.Name)
				topEnterprise.Users[string(edge.Node.EnterpriseUserAccount.Login)] = tmpUser
			}
		}
		if !EnterpriseMembersQuery.Enterprise.Members.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(EnterpriseMembersQuery.Enterprise.Members.PageInfo.EndCursor)
	}
}

func retrieveEnterpriseSamlIds(enterprise string) {
	variables := map[string]interface{}{
		"enterpriseName": githubv4.String(enterprise),
		"cursor":         (*githubv4.String)(nil),
	}

	if topEnterprise.Users == nil {
		topEnterprise.Users = map[string]User{}
	}
	for {
		err := Client.Query(Ctx, &EnterpriseSamlIdpUsersQuery, variables)
		panicOnError(err)

		for _, edge := range EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.Edges {
			if _, ok := topEnterprise.Users[string(edge.Node.User.Login)]; !ok {
				fmt.Printf("not ok, new %s\n", edge.Node.User.Login)
				topEnterprise.Users[string(edge.Node.User.Login)] = User{
					Name:       string(edge.Node.User.Name),
					SamlNameId: string(edge.Node.SamlIdentity.NameId),
				}
			} else {
				tmpUser := topEnterprise.Users[string(edge.Node.User.Login)]
				tmpUser.SamlNameId = string(edge.Node.SamlIdentity.NameId)
				topEnterprise.Users[string(edge.Node.User.Login)] = tmpUser
			}
		}
		if !EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}
}

func retrieveOrganizationUsers(organization string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(organization),
		"cursor":           (*githubv4.String)(nil),
	}

	if topOrganization.Users == nil {
		topOrganization.Users = map[string]User{}
	}
	for {
		err := Client.Query(Ctx, &OrganizationSamlIdpUsersQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.Edges {
			topOrganization.Users[string(edge.Node.User.Login)] = User{
				Name:       string(edge.Node.User.Name),
				SamlNameId: string(edge.Node.SamlIdentity.NameId),
			}
		}
		if !OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}
}
