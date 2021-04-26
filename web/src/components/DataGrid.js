import React from "react";
import { DataGrid as MUIDataGrid } from "@material-ui/data-grid";
import { makeStyles, Paper, withStyles } from "@material-ui/core";
import NoRows from "./NoRows";

const DataGridContainer = withStyles((theme) => ({
  root: {
    marginTop: theme.spacing(2),
  },
}))(Paper);

const useStyles = makeStyles(
  (theme) => ({
    odd: {
      backgroundColor: theme.palette.grey[100],
    },
  }),
  { name: "Mui" }
);

const DataGrid = ({
  onRowClick = () => {},
  columns = {},
  rows = [],
  sortField = "",
  sortDesc = false,
  pageSize = 10,
  rowHeight = 52,
  ...props
}) => {
  const classes = useStyles();

  if (sortField) {
    props.sortModel = [{ field: sortField, sort: sortDesc ? "desc" : "asc" }];
  }

  return (
    <DataGridContainer>
      <MUIDataGrid
        classes={{
          "Mui-odd": classes.odd,
        }}
        columns={columns}
        rows={rows}
        onRowClick={onRowClick}
        pageSize={pageSize}
        components={{ NoRowsOverlay: NoRows }}
        rowsPerPageOptions={[5, 10, 15, 25, 50]}
        rowHeight={rowHeight}
        autoHeight
        {...props}
      />
    </DataGridContainer>
  );
};

export default DataGrid;
