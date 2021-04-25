import { Typography } from "@material-ui/core";
import { Alert } from "@material-ui/lab";
import React, { useEffect, useState } from "react";
import { GetApplications } from "../api/Driver";
import DataGrid from "./DataGrid";

const DriverApplicationList = () => {
  const columns = [
    { field: "person_id", headerName: "User ID", hide: true },
    { field: "person_first_name", headerName: "First Name", hide: true },
    { field: "person_last_name", headerName: "Last Name", hide: true },
    { field: "organization_id", headerName: "Organization ID", hide: true },
    { field: "organization_name", headerName: "Organization", flex: 1 },
    { field: "comment", headerName: "Comment", flex: 1.5 },
    { field: "submitted_at", type: "dateTime", flex: 0.5 },
    {
      field: "status",
      headerName: "Status",
      valueGetter: (params) => {
        // Attention! JS switch statements use strict equality (===).
        switch (params.value) {
          case true:
            return "Accepted";
          case false:
            return "Rejected";
          default:
            // case for null status.
            return "Pending";
        }
      },
    },
    { field: "reason", headerName: "Reason", flex: 1 },
  ];

  const [rows, setRows] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    (async () => {
      const res = await GetApplications();

      if (res.error) {
        setError(res.error);
        return;
      }

      setRows(res.data);
    })();
  }, []);

  return (
    <>
      <Typography variant="h4">My Applications</Typography>
      {error && <Alert severity="error">{error}</Alert>}
      <DataGrid columns={columns} rows={rows} sortField={"submitted_at"} />
    </>
  );
};

export default DriverApplicationList;
