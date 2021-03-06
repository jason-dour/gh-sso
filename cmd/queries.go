// cmd/queries
//
// GraphQL queries for Github API v4.

package cmd

import "github.com/shurcooL/githubv4"

var (
	//
	// ENTERPRISE GRAPHQL QUERIES

	// EnterpriseQuery is a GraphQL query to validate a Github Enterprise Account.
	EnterpriseQuery struct {
		Enterprise struct {
			ID   githubv4.String `graphql:"id"`
			Name githubv4.String `graphql:"name"`
		} `graphql:"enterprise(slug: $enterpriseName)"`
	}

	// EnterpriseUsersQuery is a GraphQL query to retrieve all users in an Enterprise.
	EnterpriseMembersQuery struct {
		Enterprise struct {
			Members struct {
				Edges []struct {
					Node struct {
						EnterpriseUserAccount struct {
							Login githubv4.String `graphql:"login"`
							Name  githubv4.String `graphql:"name"`
						} `graphql:"... on EnterpriseUserAccount"`
					} `graphql:"node"`
				} `graphql:"edges"`
				PageInfo struct {
					EndCursor   githubv4.String  `graphql:"endCursor"`
					HasNextPage githubv4.Boolean `graphql:"hasNextPage"`
				} `graphql:"pageInfo"`
			} `graphql:"members(first: 100, after: $cursor)"`
		} `graphql:"enterprise(slug: $enterpriseName)"`
	}

	// EnterpriseSamlIdpUsersQuery is a GraphQL query for a Github Enterprise Account's SAML IDP User details.
	EnterpriseSamlIdpUsersQuery struct {
		Enterprise struct {
			OwnerInfo struct {
				SamlIdentityProvider struct {
					ExternalIdentities struct {
						Edges []struct {
							Node struct {
								Guid         githubv4.String `graphql:"guid"`
								SamlIdentity struct {
									NameId githubv4.String `graphql:"nameId"`
								} `graphql:"samlIdentity"`
								User struct {
									Login githubv4.String `graphql:"login"`
									Name  githubv4.String `graphql:"name"`
								} `graphql:"user"`
							} `graphql:"node"`
						} `graphql:"edges"`
						PageInfo struct {
							EndCursor   githubv4.String  `graphql:"endCursor"`
							HasNextPage githubv4.Boolean `graphql:"hasNextPage"`
						} `graphql:"pageInfo"`
					} `graphql:"externalIdentities(first: 100, after: $cursor)"`
				} `graphql:"samlIdentityProvider"`
			} `graphql:"ownerInfo"`
		} `graphql:"enterprise(slug: $enterpriseName)"`
	}

	// EnterpriseRepositoriesQuery lists the organizations in the enterprise.
	EnterpriseOrganizationsQuery struct {
		Enterprise struct {
			Organizations struct {
				Edges []struct {
					Node struct {
						Login githubv4.String `graphql:"login"`
					} `graphql:"node"`
				} `graphql:"edges"`
				PageInfo struct {
					EndCursor   githubv4.String  `graphql:"endCursor"`
					HasNextPage githubv4.Boolean `graphql:"hasNextPage"`
				} `graphql:"pageInfo"`
			} `graphql:"organizations(first: 100, after: $cursor1)"`
		} `graphql:"enterprise(slug: $enterpriseName)"`
	}

	//
	// ORGANIZATION GRAPHQL QUERIES

	// OrganizationQuery is a GraphQL query for a Github Organization.
	OrganizationQuery struct {
		Organization struct {
			ID   githubv4.String `graphql:"id"`
			Name githubv4.String `graphql:"name"`
		} `graphql:"organization(login: $organizationName)"`
	}

	// OrganizationSamlIdpUsersQuery is a GraphQL query for a Github Organization SAML IDP Users.
	OrganizationSamlIdpUsersQuery struct {
		Organization struct {
			SamlIdentityProvider struct {
				ExternalIdentities struct {
					Edges []struct {
						Node struct {
							Guid         githubv4.String `graphql:"guid"`
							SamlIdentity struct {
								NameId   githubv4.String `graphql:"nameId"`
								Username githubv4.String `graphql:"username"`
							} `graphql:"samlIdentity"`
							User struct {
								Login githubv4.String `graphql:"login"`
								Name  githubv4.String `graphql:"name"`
							} `graphql:"user"`
						} `graphql:"node"`
					} `graphql:"edges"`
					PageInfo struct {
						EndCursor   githubv4.String  `graphql:"endCursor"`
						HasNextPage githubv4.Boolean `graphql:"hasNextPage"`
					} `graphql:"pageInfo"`
				} `graphql:"externalIdentities(first: 100, after: $cursor)"`
			} `graphql:"samlIdentityProvider"`
		} `graphql:"organization(login: $organizationName)"`
	}

	// OrganizationRepositoriesQuery lists the repositories in the organization.
	OrganizationRepositoriesQuery struct {
		Organization struct {
			Name         githubv4.String `graphql:"name"`
			Repositories struct {
				Edges []struct {
					Node struct {
						Name githubv4.String `graphql:"name"`
					} `graphql:"node"`
				} `graphql:"edges"`
				PageInfo struct {
					EndCursor   githubv4.String  `graphql:"endCursor"`
					HasNextPage githubv4.Boolean `graphql:"hasNextPage"`
				} `graphql:"pageInfo"`
			} `graphql:"repositories(first: 100, after: $cursor)"`
		} `graphql:"organization(login: $organizationName)"`
	}
)
