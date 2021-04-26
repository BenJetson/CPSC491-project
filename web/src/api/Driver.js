import { Request } from "./Base";

const GetBalances = async () => await Request("GET", "/driver/balances");

const GetApplications = async () =>
  await Request("GET", "/driver/applications");

const SubmitApplication = async (organizationID, comment) =>
  await Request("POST", "/driver/applications/submit", {
    organization_id: organizationID,
    comment: comment,
  });

const GetMyOrganizations = async () =>
  await Request("GET", "/driver/organizations");

const GetAllOrganizations = async () =>
  await Request("GET", "/driver/organizations/all");

const SearchOrganizationCatalog = async (organizationID, keywords) => {
  const params = new URLSearchParams();
  params.set("q", keywords);

  const query = params.toString();

  return await Request(
    "GET",
    `/driver/catalog/${organizationID}/search?${query}`
  );
};

export {
  GetBalances,
  GetApplications,
  SubmitApplication,
  GetAllOrganizations,
  GetMyOrganizations,
  SearchOrganizationCatalog,
};
