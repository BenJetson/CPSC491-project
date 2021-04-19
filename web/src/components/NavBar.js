import React, { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import {
  AppBar,
  Box,
  IconButton,
  Toolbar,
  Typography,
  Link,
  makeStyles,
  Hidden,
} from "@material-ui/core";
import { Menu as MenuIcon } from "@material-ui/icons";

import { WithUser } from "../api/Auth";
import LoginStatus from "./LoginStatus";
import NavDrawer from "./NavDrawer";

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    marginBottom: theme.spacing(3),
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

  const [drawerOpen, setDrawerOpen] = useState(false);
  const toggleDrawer = () => {
    setDrawerOpen(!drawerOpen);
  };

  return (
    <>
      <NavDrawer open={drawerOpen} toggle={toggleDrawer} />
      <Box className={classes.root}>
        <AppBar position="static">
          <Toolbar>
            <IconButton
              edge="start"
              className={classes.menuButton}
              color="inherit"
              aria-label="menu"
              onClick={toggleDrawer}
            >
              <MenuIcon />
            </IconButton>
            <Typography variant="h6" component="h1" className={classes.title}>
              <Link component={RouterLink} color="inherit" to={"/"}>
                Driver Incentive Program
              </Link>
            </Typography>
            <Hidden smDown>
              <WithUser>
                {({ user, getName, getInitials, logout }) => (
                  <LoginStatus
                    user={user}
                    logout={logout}
                    getName={getName}
                    getInitials={getInitials}
                  />
                )}
              </WithUser>
            </Hidden>
          </Toolbar>
        </AppBar>
      </Box>
    </>
  );
};

export default NavBar;
