import React from "react";
import { Link as RouterLink } from "react-router-dom";

import {
  Box,
  Divider,
  Card,
  Avatar,
  CardActionArea,
  CardContent,
  Grid,
  Typography,
  makeStyles,
} from "@material-ui/core";

import Navigation from "./Navigation";

import { WithUser } from "../api/Auth";

const useStyles = makeStyles((theme) => ({
  groupBox: {
    marginTop: theme.spacing(3),
  },
  groupTitle: {
    marginTop: theme.spacing(2),
  },
  iconGrid: {
    marginTop: theme.spacing(1),
  },
  itemIcon: {
    backgroundColor: theme.palette.secondary.main,
  },
  itemTitle: {
    marginTop: theme.spacing(1),
  },
}));

const Homepage = () => {
  const classes = useStyles();

  return (
    <WithUser>
      {({ isAuthenticated, getName, isOneOfRoles }) => (
        <>
          <Typography variant="h2">
            Welcome{isAuthenticated() && `, ${getName()[0]}`}!
          </Typography>
          {Navigation.map(
            (group) =>
              (!group.roles || isOneOfRoles(group.roles)) && (
                <Box className={classes.groupBox} key={group.title}>
                  <Divider />
                  {group.title && (
                    <Typography variant="h5" className={classes.groupTitle}>
                      {group.title}
                    </Typography>
                  )}
                  <Grid container spacing={3} className={classes.iconGrid}>
                    {group.items.map(
                      (item) =>
                        // Prevent links back to the homepage (/).
                        item.link !== "/" && (
                          <Grid
                            item
                            xs={12}
                            sm={6}
                            md={4}
                            xl={3}
                            key={item.link}
                          >
                            <Card>
                              <CardActionArea
                                component={RouterLink}
                                to={item.link}
                              >
                                <CardContent>
                                  <Avatar className={classes.itemIcon}>
                                    {item.icon}
                                  </Avatar>
                                  <Typography
                                    variant="subtitle1"
                                    className={classes.itemTitle}
                                  >
                                    {item.name}
                                  </Typography>
                                </CardContent>
                              </CardActionArea>
                            </Card>
                          </Grid>
                        )
                    )}
                  </Grid>
                </Box>
              )
          )}
        </>
      )}
    </WithUser>
  );
};

export default Homepage;
