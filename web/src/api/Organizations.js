import { Request } from "./Base";

const GetOrganizations = async () => {
  const res = await Request("GET", "/organizations");
  if (res.error) {
    throw res.error;
  }

  return res.data;
};

export { GetOrganizations };
