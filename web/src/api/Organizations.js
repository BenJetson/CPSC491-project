import { Request } from "./Base";
import HTTPStatus from "./HTTPStatus";

const GetOrganizations = async () => {
  const res = await Request("GET", "/organizations");
  if (res.error) {
    throw res.error;
  }

  return res.data;
};

const GetOrgByID = async (userID) => {
  const res = await Request("GET", `/admin/organizations/${userID}`);
  if (res.error) {
    if (res.status === HTTPStatus.NOT_FOUND) {
      return false;
    }
    throw res.error;
  }

  return res.data;
};

export { GetOrganizations, GetOrgByID };
