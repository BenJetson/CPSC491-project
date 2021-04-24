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

export { GetOrganizations };
