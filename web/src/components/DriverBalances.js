import { Typography } from "@material-ui/core";
import { Alert } from "@material-ui/lab";
import React, { useEffect, useState } from "react";
import { GetBalances } from "../api/Driver";
import DataGrid from "./DataGrid";

const DriverBalances = () => {
  const columns = [
    { field: "person_id", headerName: "User ID", hide: true },
    { field: "person_first_name", headerName: "First Name", hide: true },
    { field: "person_last_name", headerName: "Last Name", hide: true },
    { field: "organization_id", headerName: "Organization ID", hide: true },
    { field: "organization_name", headerName: "Organization", flex: 1 },
    { field: "balance", headerName: "Balance", width: 170 },
  ];

  const [rows, setRows] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    (async () => {
      const res = await GetBalances();

      if (res.error) {
        setError(res.error);
        return;
      }

      let balances = [];
      for (const b of res.data) {
        balances.push({
          id: `${b.person_id}_${b.organization_id}`,
          ...b,
        });
      }

      setRows(balances);
    })();
  }, []);

  return (
    <>
      <Typography variant="h4">View Balances</Typography>
      {error && <Alert severity="error">{error}</Alert>}
      <DataGrid columns={columns} rows={rows} sortField={"organization_name"} />
    </>
  );
};

export default DriverBalances;
