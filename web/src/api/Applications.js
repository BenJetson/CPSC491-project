import { Request } from "./Base";

const SubmitApplication = (application) =>
  Request("POST", "/applications/submit", application);

export { SubmitApplication };
