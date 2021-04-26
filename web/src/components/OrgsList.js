import React, { useEffect, useState } from "react";

import { useHistory, Link as RouterLink } from "react-router-dom";
import { Typography, Button } from "@material-ui/core";
import { Add as AddIcon } from "@material-ui/icons";
import { GetOrganizations } from "../api/Admin";
import DataGrid from "./DataGrid";

const OrgsList = () => {
  const history = useHistory();

  const columns = [
    { field: "id", headerName: "ID" },
    {
      field: "name",
      headerName: "Organization Name",
      flex: 0.5,
    },
    {
      field: "point_value",
      headerName: "Exchange Rate",
      flex: 0.5,
    },
  ];

  const [rows, setRows] = useState([]);

  useEffect(() => {
    (async () => {
      const res = await GetOrganizations();
      if (res.error) {
        return;
      }
      setRows(res.data);
    })();
  }, []);

  const handleRowClick = (gridRowParams, event) => {
    // Get the user ID out of GridRowParams object.
    // GridRowParams.row is equal to the data object passed to that row.
    const id = gridRowParams.row.id; // FIXME

    history.push(`/admin/organizations/${id}`);
  };

  return (
    <>
      <Typography variant="h4">Organizations</Typography>
      <Button
        variant="contained"
        color="primary"
        component={RouterLink}
        to="/admin/organizations/create"
        style={{ marginTop: 10 }}
      >
        <AddIcon /> Create Organization
      </Button>
      <DataGrid
        columns={columns}
        rows={rows}
        onRowClick={handleRowClick}
        sortField={"name"}
      />
    </>
  );
};

export default OrgsList;
