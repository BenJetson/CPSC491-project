import HTTPStatus from "./HTTPStatus";

// Tier constants.
const TIER_PROD = "prod";
const TIER_DEV = "dev";
const TIER_LOCAL = "local";

const GetTier = () => {
  const hostname = window.location.hostname;
  let tier;

  switch (hostname) {
    case "app.teamxiv.space":
      tier = TIER_PROD;
      break;
    case "dev.teamxiv.space":
      tier = TIER_DEV;
      break;
    default:
      tier = TIER_LOCAL;
      break;
  }

  return tier;
};

const GetBaseURL = () => {
  let protocol = "https";
  if (GetTier() === TIER_LOCAL) {
    protocol = "http";
  }

  const hostname = window.location.hostname;
  const port = window.location.port;

  return `${protocol}://${hostname}:${port}/api`;
};

const Request = async (method, endpoint, data = undefined, options = {}) => {
  // Craft the full request URL.
  const url = GetBaseURL() + endpoint;

  // Attach the body data, marshaling to JSON.
  if (data !== undefined) {
    options.body = JSON.stringify(data);
  }

  // Request the data from the backend.
  const res = await fetch(url, {
    method: method,
    ...options,
  });

  // Decode the response received from the server.
  let outData = null;
  let error = null;

  try {
    if (res.status !== HTTPStatus.NO_CONTENT) {
      // Attempt to read the response body as JSON.
      outData = await res.json();
    }

    // If the status is not OK, attempt to read the error message, guarantee
    // that error !== null.
    if (res.status >= HTTPStatus.BAD_REQUEST) {
      error = outData?.message ?? outData?.status ?? "unknown error";
    }
  } catch (e) {
    // Bad JSON data would produce a SyntaxError here. This is expected. All
    // other errors are a problem and will be re-thrown.
    if (!(e instanceof SyntaxError)) {
      throw e;
    }

    // Guarantee that error !== null for bad JSON.
    error = "received malformed data from server";
  }

  // If there is an error condition, write a log as well.
  if (error !== null) {
    console.error("Request error.", { status: res.status, error: error });
  }

  return {
    status: res.status,
    data: outData,
    error: error,
  };
};

export { GetTier, GetBaseURL, Request };
