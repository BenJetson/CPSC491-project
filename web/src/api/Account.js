import { Request } from "./Base";

const DoAccountRegistration = async ({
  firstName,
  lastName,
  email,
  password,
  shouldNotify,
}) => {
  return await Request("POST", "/account/register", {
    first_name: firstName,
    last_name: lastName,
    email: email,
    password: password,
    should_notify: shouldNotify,
  });
};

export { DoAccountRegistration };
