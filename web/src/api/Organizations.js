import { Request } from "./Base";

const GetOrganizations = async () => {
  return [
    {
      id: 1,
      name: "Walmart",
      PointValue: 6,
    },
  ];

  const res = await Request("GET", "/organizations");
  if (res.error) {
    throw res.error;
  }

  return res.data;
};

const UpdateUserName = async (userID, name) =>
  //FIX change path here. unsure what would be appropriate
  await Request("POST", `/admin/organizations/${userID}/name`, {
    orgName: name,
  });

const UpdateUserEmail = async (userID, email) =>
  await Request("POST", `/admin/organizations/${userID}/email`, {
    email: email,
  });

const UpdateUserPassword = async (userID, password) =>
  await Request("POST", `/admin/organizations/${userID}/password`, {
    new_password: password,
  });

const DeactivateUser = async (userID) =>
  await Request("POST", `/admin/organizations/${userID}/deactivate`);

export {
  GetOrganizations,
  UpdateUserName,
  UpdateUserEmail,
  UpdateUserPassword,
  DeactivateUser,
};
