const RoleIDs = {
  ADMIN: 1,
  SPONSOR: 2,
  USER: 3,
  DRIVER: 4,
};

const RoleDescriptions = {
  [RoleIDs.ADMIN]: "Admin",
  [RoleIDs.SPONSOR]: "Sponsor",
  [RoleIDs.USER]: "User",
  [RoleIDs.DRIVER]: "Driver",
};

const Roles = {
  IDOf: RoleIDs,
  Describe: RoleDescriptions,
};

export default Roles;
