import React, { useEffect, useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import {
  AddVendorProductToCatalog,
  GetCatalog,
  RemoveCatalogProduct,
  SearchVendorProducts,
} from "../api/Sponsor";
import * as yup from "yup";
import { useFormik } from "formik";
import { Alert } from "@material-ui/lab";
import { FormatMoney } from "../api/Money";
import { Button, TextField, Typography } from "@material-ui/core";
import DataGrid from "./DataGrid";

const SponsorCatalog = () => {
  const [products, setProducts] = useState([]);
  const [status, setStatus] = useState(null);

  useEffect(() => {
    (async () => {
      const res = await GetCatalog();
      setProducts(!res.error ? res.data : []);
      setStatus(!res.error ? null : { success: false, message: res.error });
    })();
  }, []);

  const makeRemoveHandler = (rowIndex, productID) => async () => {
    const res = await RemoveCatalogProduct(productID);
    setStatus(
      !res.error
        ? {
            success: true,
            message: `Removed product #${productID} from the catalog.`,
          }
        : { success: false, message: res.error }
    );

    console.log(rowIndex, products);
    if (!res.error) {
      setProducts(
        products
          .slice(0, rowIndex)
          .concat(products.slice(rowIndex + 1, products.length))
      );
    }
  };

  const columns = [
    {
      field: "id",
      headerName: "ID",
      hide: true,
    },
    {
      field: "title",
      headerName: "Title",
      flex: 1,
    },
    {
      field: "description",
      headerName: "Description",
      flex: 2,
      hide: true,
    },
    {
      field: "points",
      headerName: "Cost (points)",
      type: "number",
      valueGetter: (params) => params.value.amount,
    },
    {
      field: "image_url",
      headerName: "Image",
      renderCell: (params) => <img src={params.value} />,
      width: 170,
    },
    {
      field: "remove",
      headerName: "Remove",
      renderCell: (params) => {
        const rowIndex = params.rowIndex;
        const productID = params.getValue("id");
        return (
          <Button
            variant="contained"
            color="secondary"
            onClick={makeRemoveHandler(rowIndex, productID)}
          >
            Remove
          </Button>
        );
      },
      width: 125,
      sortable: false,
      filterable: false,
    },
  ];

  return (
    <>
      <Typography variant="h4">Manage Catalog</Typography>

      <Button
        variant="contained"
        color="primary"
        component={RouterLink}
        to="/sponsor/catalog/vendor"
      >
        Add Products
      </Button>

      {status && (
        <Alert severity={status.success ? "success" : "error"}>
          {status.message}
        </Alert>
      )}

      <DataGrid
        columns={columns}
        rows={products}
        sortField={"title"}
        rowHeight={135}
        pageSize={5}
      />
    </>
  );
};

export default SponsorCatalog;
