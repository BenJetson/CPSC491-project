import { Link } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { Alert, AlertTitle } from "@material-ui/lab";
import React from "react";

const LoginRequired = () => (
  <Alert severity="warning">
    <AlertTitle>Access Denied</AlertTitle>
    To access this application, first you must authenticate. Visit the{" "}
    <Link component={RouterLink} to="/login">
      login page
    </Link>{" "}
    to get started.
  </Alert>
);

export default LoginRequired;
