import React, { useState } from "react";
import { useHistory, Link as RouterLink } from "react-router-dom";
import { Button, makeStyles, Typography } from "@material-ui/core";
import { Add as AddIcon } from "@material-ui/icons";
import DataGrid from "./DataGrid";

const useStyles = makeStyles((theme) => ({
  title: {
    marginBottom: theme.spacing(2),
  },
  addBtn: {
    marginBottom: theme.spacing(2),
    marginLeft: "auto",
  },
}));

const ApplicationList = ({ applications, isSponsor = false }) => {
  const history = useHistory();
  const classes = useStyles();

  const columns = [
    { field: "id", headerName: "ID" },
    {
      field: "first_name",
      headerName: "First Name",
      hide: !isSponsor,
      flex: 0.5,
    },
    {
      field: "last_name",
      headerName: "Last Name",
      hide: !isSponsor,
      flex: 0.5,
    },
    { field: "organization_id", headerName: "Organization ID", hide: true },
    {
      field: "organization_name",
      headerName: "Organization",
      hide: isSponsor,
      flex: 1,
    },
    { field: "submitted_at", headerName: "Submitted At", flex: 1 },
    { field: "status", headerName: "Status" },
  ];

  const [rows, setRows] = useState([]);

  const handleRowClick = (gridRowParams, event) => {
    // get the ID out of gridrowparams
    const id = 0; // FIXME

    // redirect to appropriate viewer page
    const viewPath = isSponsor
      ? `/sponsor/applications/${id}`
      : `/driver/applications/${id}`;

    history.push(viewPath);
  };

  return (
    <>
      <Typography variant="h4" className={classes.title}>
        Applications
      </Typography>
      {!isSponsor && (
        <Button
          component={RouterLink}
          to="/driver/applications/new"
          variant="contained"
          color="primary"
          className={classes.addBtn}
        >
          <AddIcon />
          New Application
        </Button>
      )}

      <DataGrid columns={columns} rows={rows} onRowClick={handleRowClick} />
    </>
  );
};

export default ApplicationList;
