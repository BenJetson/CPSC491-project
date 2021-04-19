import React from "react";

import { Link as RouterLink } from "react-router-dom";
import {
  Divider,
  Drawer,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Typography,
  makeStyles,
  Box,
} from "@material-ui/core";

import Navigation from "./Navigation";
import LoginStatus from "./LoginStatus";

import { WithUser } from "../api/Auth";

let useStyles = makeStyles((theme) => ({
  header: {
    margin: theme.spacing(3),
  },
}));

let NavDrawer = ({ open, toggle }) => {
  const classes = useStyles();

  return (
    <WithUser>
      {({ user, logout, getName, getInitials }) => (
        <Drawer open={open} onClose={toggle}>
          <Box className={classes.header}>
            <Typography variant="h6">Driver Incentive Program</Typography>
            <LoginStatus
              user={user}
              logout={logout}
              getName={getName}
              getInitials={getInitials}
              stacked
            />
          </Box>
          {Navigation.map(
            (group) =>
              (!group.roles ||
                group.roles.includes(user?.role_id ?? false)) && (
                <>
                  <Divider />
                  <List>
                    {group.title && (
                      <ListItem dense={true}>
                        <Typography variant="overline">
                          {group.title}
                        </Typography>
                      </ListItem>
                    )}
                    {group.items.map((item) => (
                      <ListItem
                        component={RouterLink}
                        button
                        to={item.link}
                        onClick={toggle}
                      >
                        <ListItemIcon>{item.icon}</ListItemIcon>
                        <ListItemText primary={item.name} />
                      </ListItem>
                    ))}
                  </List>
                </>
              )
          )}
        </Drawer>
      )}
    </WithUser>
  );
};

export default NavDrawer;
