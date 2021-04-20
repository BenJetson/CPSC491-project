import React from "react";
import { Avatar, Button, Box, Chip, makeStyles } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";

const useStyles = makeStyles((theme) => ({
  inlineChip: {
    display: "inline",
    marginRight: theme.spacing(1),
  },
  flexStack: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    flexDirection: "column",
  },
  stackedItem: {
    marginBottom: theme.spacing(2),
  },
}));

const LoginStatus = ({
  user,
  logout,
  getName,
  getInitials,
  stacked = false,
}) => {
  const classes = useStyles();
  const [firstName, lastName] = getName();

  return (
    <Box className={stacked ? classes.flexStack : null}>
      {(user && (
        <>
          <Box className={stacked ? classes.stackedItem : classes.inlineChip}>
            <Chip
              avatar={<Avatar>{getInitials()}</Avatar>}
              label={`${firstName} ${lastName}`}
            />
          </Box>
          <Button size="small" variant="contained" onClick={logout}>
            Logout
          </Button>
        </>
      )) || (
        <Button
          size="small"
          variant="contained"
          component={RouterLink}
          to="/login"
        >
          Login
        </Button>
      )}
    </Box>
  );
};

export default LoginStatus;
