import React, { useEffect, useState } from "react";

import { useHistory, Link as RouterLink } from "react-router-dom";
import { Typography, Button } from "@material-ui/core";
import { Add as AddIcon } from "@material-ui/icons";
import { GetAllUsers } from "../api/Admin";
import DataGrid from "./DataGrid";
import Roles from "../api/Roles";

const UserList = () => {
  const history = useHistory();

  const columns = [
    { field: "id", headerName: "ID" },
    {
      field: "first_name",
      headerName: "First Name",
      flex: 0.5,
    },
    {
      field: "last_name",
      headerName: "Last Name",
      flex: 0.5,
    },
    {
      field: "email",
      headerName: "Email Address",
      flex: 1,
    },
    {
      field: "role",
      headerName: "Role",
      valueGetter: (params) => Roles.Describe[params.row.role_id] ?? "Unknown",
    },
    {
      field: "active",
      headerName: "Active?",
      type: "boolean",
      valueGetter: (params) => !params.row.is_deactivated,
    },
  ];

  const [rows, setRows] = useState([]);

  useEffect(() => {
    (async () => {
      const users = await GetAllUsers();
      setRows(users);
    })();
  }, []);

  const handleRowClick = (gridRowParams, event) => {
    // Get the user ID out of GridRowParams object.
    // GridRowParams.row is equal to the data object passed to that row.
    const id = gridRowParams.row.id; // FIXME

    history.push(`/admin/users/${id}`);
  };

  return (
    <>
      <Typography variant="h4">Users</Typography>
      <Button
        variant="contained"
        color="primary"
        component={RouterLink}
        to="/admin/users/create"
        style={{ marginTop: 10 }}
      >
        <AddIcon /> Create User
      </Button>
      <DataGrid
        columns={columns}
        rows={rows}
        onRowClick={handleRowClick}
        sortField={"last_name"}
      />
    </>
  );
};

export default UserList;
