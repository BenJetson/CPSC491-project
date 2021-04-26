import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  GetOrganizationByID,
  //CreateOrganization,
  UpdateOrgName,
  UpdateOrgRate,
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

  const [nameStatus, setNameStatus] = useState(null);
  const nameForm = useFormik({
    initialValues: {
      firstName: org.first_name,
      lastName: org.last_name,
    },
    enableReinitialize: true,
    validationSchema: nameValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateOrgName(orgID, values.name);

      setOrg({
        // Force dirty state validation.
        ...org,
        name: !res.error ? values.name : org.name,
      });

      setNameStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Name changed successfully." }
      );
    },
  });

  const [rateStatus, setRateStatus] = useState(null);
  const rateForm = useFormik({
    initialValues: {
      rate: org.rate,
    },
    enableReinitialize: true,
    validationSchema: rateValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateOrgRate(orgID, values.rate);

      setOrg({
        // Force dirty state validation.
        ...org,
        rate: !res.error ? values.rate : org.rate,
      });

      setRateStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Exchange rate changed successfully." }
      );
    },
  });

  return (
    <>
      <Typography variant="h4">Edit Org #{org.id}</Typography>
      <FormCard>
        <CardContent>
          <Typography variant="h5">Name</Typography>
          {nameStatus && (
            <Alert severity={nameStatus.success ? "success" : "error"}>
              {nameStatus.message}
            </Alert>
          )}
          <form noValidate onSubmit={nameForm.handleSubmit}>
            <TextField
              variant="outlined"
              required
              fullWidth
              margin="normal"
              id="name"
              name="name"
              label="name"
              value={nameForm.values.name}
              onChange={nameForm.handleChange}
              error={nameForm.touched.name && Boolean(nameForm.errors.name)}
              helperText={nameForm.touched.name && nameForm.errors.name}
            />
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={!nameForm.dirty}
            >
              Save
            </Button>
          </form>
        </CardContent>
      </FormCard>
      <FormCard>
        <CardContent>
          <Typography variant="h5">Exchange Rate</Typography>
          {rateStatus && (
            <Alert severity={rateStatus.success ? "success" : "error"}>
              {rateStatus.message}
            </Alert>
          )}
          <form noValidate onSubmit={rateForm.handleSubmit}>
            <TextField
              variant="outlined"
              required
              fullWidth
              margin="normal"
              id="rate"
              name="rate"
              label="Exchange Rate"
              type="rate"
              value={rateForm.values.rate}
              onChange={rateForm.handleChange}
              error={rateForm.touched.rate && Boolean(rateForm.errors.rate)}
              helperText={rateForm.touched.rate && rateForm.errors.rate}
            />
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={!rateForm.dirty}
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
