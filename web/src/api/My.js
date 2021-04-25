import { Request } from "./Base";
import { GetCurrentUser } from "./Auth";

const GetMyUser = GetCurrentUser;

const UpdateMyName = async (firstName, lastName) => {
  const res = await Request("POST", `/my/profile/name`, {
    first_name: firstName,
    last_name: lastName,
  });

  if (res.error) {
    throw res.error;
  }
};

const UpdateMyEmail = async (email) => {
  const res = await Request("POST", `/my/profile/email`, {
    email: email,
  });

  if (res.error) {
    throw res.error;
  }
};

const UpdateMyPassword = async (password) => {
  const res = await Request("POST", `/my/profile/password`, {
    new_password: password,
  });

  if (res.error) {
    throw res.error;
  }
};

export { GetMyUser, UpdateMyName, UpdateMyEmail, UpdateMyPassword };
