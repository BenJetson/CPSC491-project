import React, { useState } from "react";
import { DataGrid, GridRowsProp, GridColDef } from "@material-ui/data-grid";
import { useHistory } from "react-router-dom";
import NoRows from "./NoRows";

const ApplicationList = ({ applications, isSponsor = false }) => {
  const history = useHistory();

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
    const id = 0;

    // redirect to appropriate viewer page
    const viewPath = isSponsor
      ? `/sponsor/applications/${id}`
      : `/driver/applications/${id}`;

    history.push(viewPath);
  };

  return (
    <DataGrid
      columns={columns}
      rows={rows}
      components={{ NoRowsOverlay: NoRows }}
      onRowClick={handleRowClick}
    />
  );
};

export default ApplicationList;
