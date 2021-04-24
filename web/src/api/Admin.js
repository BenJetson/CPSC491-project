import { Request } from "./Base";
import HTTPStatus from "./HTTPStatus";

const GetAllUsers = async () => {
  const res = await Request("GET", "/admin/users");
  if (res.error) {
    throw res.error;
  }

  return res.data;
};

const GetUserByID = async (userID) => {
  const res = await Request("GET", `/admin/users/${userID}`);
  if (res.error) {
    if (res.status === HTTPStatus.NOT_FOUND) {
      return false;
    }
    throw res.error;
  }

  return res.data;
};

const UpdateUserName = async (userID, firstName, lastName) =>
  await Request("POST", `/admin/users/${userID}/name`, {
    first_name: firstName,
    last_name: lastName,
  });

const UpdateUserEmail = async (userID, email) =>
  await Request("POST", `/admin/users/${userID}/email`, {
    email: email,
  });

const UpdateUserPassword = async (userID, password) =>
  await Request("POST", `/admin/users/${userID}/password`, {
    new_password: password,
  });

const UpdateUserAffiliations = async (userID, organization_ids) =>
  await Request("POST", `/admin/users/${userID}/affiliations`, {
    organization_ids: organization_ids,
  });

const ActivateUser = async (userID) => {
  const res = await Request("POST", `/admin/users/${userID}/activate`);

  if (res.error) {
    throw res.error;
  }
};

const DeactivateUser = async (userID) => {
  const res = await Request("POST", `/admin/users/${userID}/deactivate`);

  if (res.error) {
    throw res.error;
  }
};

export {
  GetAllUsers,
  GetUserByID,
  UpdateUserName,
  UpdateUserEmail,
  UpdateUserPassword,
  UpdateUserAffiliations,
  ActivateUser,
  DeactivateUser,
};
