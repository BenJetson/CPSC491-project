import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  GetSponsorOrganization,
  UpdateSponsorOrganization,
} from "../api/Sponsor";
import { DeleteOrganization } from "../api/Admin";
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
  name: "",
  rate: 1,
};

const validationSchema = yup.object({
  name: yup
    .string("Enter the new organization name.")
    .required("Organization name is required."),
  rate: yup.number().min(1, "Exchange rate must be a positive integer."),
});

const OrgProfileEditor = () => {
  const params = useParams();
  const orgID = parseInt(params["orgID"]) ?? false; // FIXME unchecked cast
  const isUpdate = orgID !== false && orgID > 0;

  const [org, setOrg] = useState(emptyOrg);

  useEffect(() => {
    if (!isUpdate) {
      setOrg(emptyOrg);
      return;
    }

    (async () => {
      const data = await GetSponsorOrganization();
      setOrg(data);
    })();
  }, [orgID, isUpdate]);

  const [status, setStatus] = useState(null);
  const formik = useFormik({
    initialValues: {
      name: org.name,
      rate: org.point_value,
    },
    enableReinitialize: true,
    validationSchema: validationSchema,
    onSubmit: async (values) => {
      const res = await UpdateSponsorOrganization(
        orgID,
        values.name,
        values.rate
      );

      setOrg({
        // Force dirty state validation.
        ...org,
        name: !res.error ? values.name : org.name,
      });

      setStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Name changed successfully." }
      );
    },
  });

  const [activationStatus, setActivationStatus] = useState(null);
  const doDeactivation = async () => {
    const res = await DeleteOrganization();
    setOrg({
      // Force dirty state validation.
      ...org,
      is_deactivated: !res.error ? true : org.is_deactivated,
    });
    setActivationStatus(
      res.error
        ? { success: false, message: res.error }
        : { success: true, message: "Account deactivated successfully." }
    );
  };

  return (
    <>
      <Typography variant="h4">Edit My Organization</Typography>
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
              id="orgName"
              name="name"
              label="Organization Name"
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
              label="Points Per US Dollar"
              type="rate"
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
      <FormCard>
        <CardContent>
          <Typography variant="h5">Account Status</Typography>
          {activationStatus && (
            <Alert severity={activationStatus.success ? "success" : "error"}>
              {activationStatus.message}
            </Alert>
          )}

          <Typography
            style={{ marginTop: 15 }} // FIXME
          >
            This account is currently{" "}
            <strong>{org.is_deactivated ? "deactivated" : "activated"}</strong>.
          </Typography>
          <Button
            onClick={org.is_deactivated ? false : doDeactivation}
            variant="contained"
            color={org.is_deactivated ? "secondary" : "primary"}
            style={{ marginTop: 15 }} // FIXME
          >
            {org.is_deactivated
              ? "contact an administrator to reactivate"
              : "deactivate"}{" "}
            account
          </Button>
        </CardContent>
      </FormCard>
    </>
  );
};

export default OrgProfileEditor;
