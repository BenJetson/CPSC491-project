import { Request } from "./Base";
import { GetCurrentUser } from "./Auth";

const GetMyUser = GetCurrentUser;

const UpdateUserName = async (firstName, lastName) =>
  await Request("POST", `/my/profile/name`, {
    first_name: firstName,
    last_name: lastName,
  });

const UpdateUserEmail = async (email) =>
  await Request("POST", `/my/profile/email`, {
    email: email,
  });

const UpdateUserPassword = async (currentPassword, newPassword) =>
  await Request("POST", `/my/profile/password`, {
    new_password: newPassword,
    current_password: currentPassword,
  });

const DeactivateUser = async () =>
  await Request("POST", `/my/profile/deactivate`);

export {
  GetMyUser,
  UpdateUserName,
  UpdateUserEmail,
  UpdateUserPassword,
  DeactivateUser,
};
