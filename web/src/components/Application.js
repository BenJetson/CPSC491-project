import React from "react";
import { makeStyles, Card, CardContent, Typography } from "@material-ui/core";

const useStyles = makeStyles({
  root: {
    minWidth: 275,
  },
  bullet: {
    display: "inline-block",
    margin: "0 2px",
    transform: "scale(0.8)",
  },
  title: {
    fontSize: 14,
  },
  pos: {
    marginBottom: 12,
  },
});

export default function SimpleCard() {
  const classes = useStyles();
  const bull = <span className={classes.bullet}>•</span>;

  return (
    <Card className={classes.root}>
      <CardContent>
        <Typography
          className={classes.title}
          color="textSecondary"
          gutterBottom
        >
          Application Overview
        </Typography>
        <Typography variant="h5" component="h2">
          Clemson Global Shipping company
        </Typography>
        <Typography className={classes.pos} color="textSecondary">
          Status
          <br />
          {"In Progress"}
        </Typography>
        <Typography variant="body2" component="p">
          Reason
          <br />
          {'"Awaiting more information requested from driver"'}
        </Typography>
      </CardContent>
    </Card>
  );
}
