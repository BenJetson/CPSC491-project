import { Typography } from "@material-ui/core";
import { Alert, AlertTitle } from "@material-ui/lab";
import React from "react";

const AccessDenied = ({ reason = null }) => (
  <Alert severity="warning">
    <AlertTitle>Access Denied</AlertTitle>
    <Typography>
      You are logged in, but your account does not have access to this page.
    </Typography>
    {(reason && (
      <Typography>
        <strong>Reason:</strong> {reason}
      </Typography>
    )) || (
      <Typography>
        <Typography>This could be for a number of reasons:</Typography>
        <ul>
          <li>Your privilege level is insufficient (wrong role).</li>
          <li>Your account is not the owner of the requested resource.</li>
          <li>
            Your account is not affiliated with the same organization as the
            requested resource.
          </li>
        </ul>
        <Typography>
          This error is durable and will persist unless your account parameters
          are modified by an administrator.
        </Typography>
      </Typography>
    )}
  </Alert>
);

export default AccessDenied;
