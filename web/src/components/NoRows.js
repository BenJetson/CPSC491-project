import React from "react";

import { Avatar, Box, Typography, makeStyles } from "@material-ui/core";
import { FindInPage as FindInPageIcon } from "@material-ui/icons";

const useStyles = makeStyles((theme) => ({
  surface: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    paddingTop: theme.spacing(16),
    paddingBottom: theme.spacing(10),
  },
  text: {
    marginTop: theme.spacing(2),
  },
}));

const NoRows = () => {
  const classes = useStyles();

  return (
    <Box className={classes.surface}>
      <Avatar>
        <FindInPageIcon />
      </Avatar>
      <Typography className={classes.text}>No data to display.</Typography>
    </Box>
  );
};

export default NoRows;
