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

const ApplicationCard = ({ companyName, status, reason }) => {
  const classes = useStyles();
  //const bull = <span className={classes.bullet}>•</span>;

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
        <Typography variant="h4" component="h2">
          {companyName}
        </Typography>
        <Typography variant="h5" component="h2">
          Status
        </Typography>
        <Typography className={classes.pos} color="textSecondary">
          {status}
        </Typography>
        <Typography variant="h5" component="h2">
          Reason
        </Typography>
        <Typography className={classes.pos} color="textSecondary">
          {reason}
        </Typography>
      </CardContent>
    </Card>
  );
};

export default ApplicationCard;
