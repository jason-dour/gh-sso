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

	orgs := map[string]Organization{}
	for {
		err := Client.Query(Ctx, &EnterpriseOrganizationsQuery, variables)
		panicOnError(err)

		for _, edge := range EnterpriseOrganizationsQuery.Enterprise.Organizations.Edges {
			orgs[string(edge.Node.Login)] = Organization{
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
	topEnterprise.Organizations = orgs
}

func retrieveEnterpriseRepos(enterprise string, org string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(org),
		"cursor":           (*githubv4.String)(nil),
	}

	repos := map[string]Repository{}
	for {
		err := Client.Query(Ctx, &OrganizationRepositoriesQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationRepositoriesQuery.Organization.Repositories.Edges {
			repos[string(edge.Node.Name)] = Repository{
				Name: string(edge.Node.Name),
			}
		}

		if !OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.EndCursor)
	}
	org2 := topEnterprise.Organizations[org]
	org2.Repositories = repos
}

func retrieveOrganizationRepos(organization string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(organization),
		"cursor":           (*githubv4.String)(nil),
	}

	repos := map[string]Repository{}
	for {
		err := Client.Query(Ctx, &OrganizationRepositoriesQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationRepositoriesQuery.Organization.Repositories.Edges {
			repos[string(edge.Node.Name)] = Repository{
				Name: string(edge.Node.Name),
			}
		}

		if !OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationRepositoriesQuery.Organization.Repositories.PageInfo.EndCursor)
	}
	topOrganization.Repositories = repos
}

func retrieveEnterpriseUsers(enterprise string) {
	variables := map[string]interface{}{
		"enterpriseName": githubv4.String(enterprise),
		"cursor":         (*githubv4.String)(nil),
	}

	users := map[string]User{}
	for {
		err := Client.Query(Ctx, &EnterpriseSamlIdpUsersQuery, variables)
		panicOnError(err)

		for _, edge := range EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.Edges {
			users[string(edge.Node.User.Login)] = User{
				Login:      string(edge.Node.User.Login),
				SamlNameId: string(edge.Node.SamlIdentity.NameId),
			}
		}
		if !EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(EnterpriseSamlIdpUsersQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}
	topEnterprise.Users = users
}

func retrieveOrganizationUsers(organization string) {
	variables := map[string]interface{}{
		"organizationName": githubv4.String(organization),
		"cursor":           (*githubv4.String)(nil),
	}

	users := map[string]User{}
	for {
		err := Client.Query(Ctx, &OrganizationSamlIdpUsersQuery, variables)
		panicOnError(err)

		for _, edge := range OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.Edges {
			users[string(edge.Node.User.Login)] = User{
				Login:      string(edge.Node.User.Login),
				SamlNameId: string(edge.Node.SamlIdentity.NameId),
			}
		}
		if !OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(OrganizationSamlIdpUsersQuery.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}
	topOrganization.Users = users
}
