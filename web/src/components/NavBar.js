import React from "react";
import {
  AppBar,
  Button,
  Box,
  Toolbar,
  Typography,
  makeStyles,
} from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
}));

const NavBar = () => {
  const classes = useStyles();

  const doLogout = () => {
    fetch("/api/logout", { method: "POST" });
  };

  return (
    <Box className={classes.root}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" className={classes.title}></Typography>
          <Button onClick={doLogout}>Logout</Button>
        </Toolbar>
      </AppBar>
    </Box>
  );
};

export default NavBar;
