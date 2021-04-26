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

const UpdateUserName = async (userID, name) =>
  //need to replace these userIDs with currently logged in organization
  await Request("POST", `/admin/organizations/${userID}/name`, {
    orgName: name,
  });

const UpdateUserEmail = async (userID, email) =>
  await Request("POST", `/admin/organizations/${userID}/email`, {
    email: email,
  });

const UpdateUserPassword = async (userID, currentPassword, newPassword) =>
  await Request("POST", `/admin/organizations/${userID}/password`, {
    new_password: newPassword,
    current_password: currentPassword,
  });

const UpdateConversionRate = async (userID, rate) =>
  await Request("POST", `/admin/organizations/${userID}/rate`, {
    new_rate: rate,
  });

const DeactivateUser = async (userID) =>
  await Request("POST", `/admin/organizations/${userID}/deactivate`);

export {
  GetOrganizations,
  GetOrgByID,
  UpdateUserName,
  UpdateUserEmail,
  UpdateUserPassword,
  UpdateConversionRate,
  DeactivateUser,
};
