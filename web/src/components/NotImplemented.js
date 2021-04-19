import { makeStyles } from "@material-ui/core";
import { useLocation } from "react-router-dom";
import { Alert, AlertTitle } from "@material-ui/lab";
import React from "react";

const useStyles = makeStyles((theme) => ({
  error: {
    marginBottom: theme.spacing(2),
  },
}));

const NotImplemented = ({ feature = null }) => {
  const classes = useStyles();
  const location = useLocation();

  return (
    <>
      <Alert severity="error" className={classes.error}>
        <AlertTitle>Not Implemented</AlertTitle>
        The page you requested is part of a feature that has not yet been
        implemented by the developers of this application.
        {feature && (
          <>
            <br />
            <br />
            <strong>Feature:</strong> {feature}
          </>
        )}
      </Alert>
      <Alert severity="info">
        <AlertTitle>Debugging Information</AlertTitle>
        <strong>Path:</strong> {location.pathname}
      </Alert>
    </>
  );
};

export default NotImplemented;
