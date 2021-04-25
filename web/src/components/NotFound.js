import { Link } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { Alert, AlertTitle } from "@material-ui/lab";
import React from "react";

const NotFound = () => (
  <Alert severity="error">
    <AlertTitle>Page Not Found</AlertTitle>
    The page you have requested was not found within this application. Consider
    nagivating from the{" "}
    <Link component={RouterLink} to="/">
      homepage
    </Link>
    , perhaps?
  </Alert>
);

export default NotFound;
