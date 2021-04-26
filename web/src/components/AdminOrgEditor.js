import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  GetOrganizationByID,
  //CreateOrganization,
  UpdateOrganization,
  //DeleteOrganization,
} from "../api/Admin";

import * as yup from "yup";
import { useFormik } from "formik";
import {
  Button,
  Card,
  CardContent,
  TextField,
  Typography,
  withStyles,
} from "@material-ui/core";
import { Alert } from "@material-ui/lab";

const FormCard = withStyles((theme) => ({
  root: {
    marginTop: theme.spacing(2),
  },
}))(Card);

const emptyOrg = {
  id: 0,
  name: "",
  rate: 1,
};

const nameValidationSchema = yup.object({
  name: yup
    .string("Enter the new organization name.")
    .required("Organization name is required."),
});
const rateValidationSchema = yup.object({
  rate: yup.number().min(1, "Exchange rate must be a positive integer."),
});

const AdminOrgEditor = () => {
  const params = useParams();
  const orgID = parseInt(params["orgID"]) ?? false; // FIXME unchecked cast
  const isUpdate = orgID !== false && orgID > 0;

  const [org, setOrg] = useState(emptyOrg);

  useEffect(() => {
    if (!isUpdate) {
      setOrg(emptyOrg);
    }

    (async () => {
      const data = await GetOrganizationByID(orgID);
      setOrg(data);
    })();
  }, [orgID, isUpdate]);

  const [status, setStatus] = useState(null);
  const formik = useFormik({
    initialValues: {
      name: org.name,
      rate: org.rate,
    },
    enableReinitialize: true,
    validationSchema: nameValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateOrganization(orgID, values.name, values.rate);

      setOrg({
        // Force dirty state validation.
        ...org,
        name: !res.error ? values.name : org.name,
        point_value: !res.error ? values.rate : org.point_value,
      });

      setStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Updated organization successfully." }
      );
    },
  });

  return (
    <>
      <Typography variant="h4">Edit Org #{org.id}</Typography>
      <FormCard>
        <CardContent>
          <Typography variant="h5">Name</Typography>
          {status && (
            <Alert severity={status.success ? "success" : "error"}>
              {status.message}
            </Alert>
          )}
          <form noValidate onSubmit={formik.handleSubmit}>
            <TextField
              variant="outlined"
              required
              fullWidth
              margin="normal"
              id="name"
              name="name"
              label="name"
              value={formik.values.name}
              onChange={formik.handleChange}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
            />
            <TextField
              variant="outlined"
              required
              fullWidth
              margin="normal"
              id="rate"
              name="rate"
              label="Exchange Rate"
              type="number"
              value={formik.values.rate}
              onChange={formik.handleChange}
              error={formik.touched.rate && Boolean(formik.errors.rate)}
              helperText={formik.touched.rate && formik.errors.rate}
            />
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={!formik.dirty}
            >
              Save
            </Button>
          </form>
        </CardContent>
      </FormCard>
    </>
  );
};

export default AdminOrgEditor;
