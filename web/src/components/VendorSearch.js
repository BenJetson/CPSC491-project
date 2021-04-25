import React, { useEffect, useState } from "react";
import {
  AddVendorProductToCatalog,
  SearchVendorProducts,
} from "../api/Sponsor";
import { useFormik } from "formik";
import { Alert } from "@material-ui/lab";
import { FormatMoney } from "../api/Money";
import { Button, TextField, Typography } from "@material-ui/core";
import DataGrid from "./DataGrid";

const VendorSearch = () => {
  const [keywords, setKeywords] = useState("");
  const [products, setProducts] = useState([]);
  const [status, setStatus] = useState(null);

  const formik = useFormik({
    initialValues: { keywords: "" },
    onSubmit: (values) => setKeywords(values.keywords),
  });

  useEffect(() => {
    if (!keywords) {
      setProducts([]);
      return;
    }

    (async () => {
      const res = await SearchVendorProducts(keywords);
      setProducts(!res.error ? res.data : []);
      setStatus(!res.error ? null : { success: false, message: res.error });
    })();
  }, [keywords]);

  const makeAddHandler = (productID) => async () => {
    const res = await AddVendorProductToCatalog(productID);
    setStatus(
      !res.error
        ? {
            success: true,
            message: `Added product #${productID} to the catalog.`,
          }
        : { success: false, message: res.error }
    );
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
      field: "price",
      headerName: "Price",
      valueFormatter: (params) => FormatMoney(params.value),
    },
    {
      field: "image_url",
      headerName: "Image",
      renderCell: (params) => (
        <img src={params.value} alt={`${params.getValue("title")}`} />
      ),
      width: 170,
    },
    {
      field: "add",
      headerName: "Add",
      renderCell: (params) => {
        const productID = params.getValue("id");
        return (
          <Button
            variant="contained"
            color="primary"
            onClick={makeAddHandler(productID)}
          >
            Add
          </Button>
        );
      },
      width: 120,
      sortable: false,
      filterable: false,
    },
  ];

  return (
    <>
      <Typography variant="h4">Search Vendor Products</Typography>
      <form noValidate onSubmit={formik.handleSubmit}>
        <TextField
          variant="outlined"
          fullWidth
          id="keywords"
          name="keywords"
          label="Keywords"
          value={formik.values.keywords}
          onChange={formik.handleChange}
          error={formik.touched.keywords && Boolean(formik.errors.keywords)}
          helperText={formik.touched.keywords && formik.errors.keywords}
        />
        <Button variant="contained" color="primary" type="submit">
          Search
        </Button>
      </form>
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

export default VendorSearch;
