import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  DeactivateUser,
  GetOrgByID,
  UpdateUserEmail,
  UpdateUserName,
  UpdateUserPassword,
  UpdateConversionRate,
} from "../api/Organizations";
import Roles from "../api/Roles";

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

const emptyUser = {
  id: 0,
  name: "",
  email: "",
  role_id: Roles.IDOf.SPONSOR,
  is_deactivated: false,
};

const nameValidationSchema = yup.object({
  lastName: yup
    .string("Enter the new organization name.")
    .required("Organization name is required."),
});
const emailValidationSchema = yup.object({
  email: yup
    .string("Enter the new email.")
    .email("Enter a valid email.")
    .required("Email is required."),
});
const passwordValidationSchema = yup.object({
  password: yup
    .string("Enter the new password.")
    .min(8, "Password should be of minimum 8 characters in length.")
    .required("Password is required."),
  confirm: yup
    .string("Re-enter the new password.")
    .required("Password confirmation is not optional ."),
});

const OrgProfileEditor = () => {
  const params = useParams();
  const userID = parseInt(params["userID"]) ?? false; // FIXME unchecked cast
  const isUpdate = userID !== false && userID > 0;

  const [user, setUser] = useState(emptyUser);

  useEffect(() => {
    if (!isUpdate) {
      setUser(emptyUser);
    }

    (async () => {
      const data = await GetOrgByID(userID);
      setUser(data);
    })();
  }, [userID, isUpdate]);

  const [nameStatus, setNameStatus] = useState(null);
  const nameForm = useFormik({
    initialValues: {
      name: user.name,
    },
    enableReinitialize: true,
    validationSchema: nameValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateUserName(userID, values.name);

      setUser({
        // Force dirty state validation.
        ...user,
        name: !res.error ? values.name : user.name,
      });

      setNameStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Name changed successfully." }
      );
    },
  });

  const [emailStatus, setEmailStatus] = useState(null);
  const emailForm = useFormik({
    initialValues: {
      email: user.email,
    },
    enableReinitialize: true,
    validationSchema: emailValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateUserEmail(userID, values.email);

      setUser({
        // Force dirty state validation.
        ...user,
        email: !res.error ? values.email : user.email,
      });

      setEmailStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Email changed successfully." }
      );
    },
  });

  // Password confirm logic inspired by: https://stackoverflow.com/a/62189211
  const validateConfirmPassword = (expect, actual) =>
    expect && actual && expect !== actual ? "Passwords do not match." : "";

  const [passwordStatus, setPasswordStatus] = useState(null);
  const passwordForm = useFormik({
    initialValues: {
      password: "",
      confirm: "",
    },
    validationSchema: passwordValidationSchema,
    validate: async (values) => {
      let errors = {};

      let confirmPasswordError = validateConfirmPassword(
        values.password,
        values.confirm
      );
      if (confirmPasswordError) {
        errors.confirm = confirmPasswordError;
      }

      return errors;
    },
    onSubmit: async (values) => {
      const res = await UpdateUserPassword(
        userID,
        values.password,
        values.newpwd
      );

      setPasswordStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Password changed successfully." }
      );
    },
  });

  const [rateStatus, setConversionRate] = useState(null);
  const rateForm = useFormik({
    initialValues: {
      rate: user.rate,
    },
    enableReinitialize: true,
    validationSchema: nameValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateConversionRate(userID, values.rate);

      setUser({
        // Force dirty state validation.
        ...user,
        rate: !res.error ? values.rate : user.rate,
      });

      setNameStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Name changed successfully." }
      );
    },
  });

  const [activationStatus, setActivationStatus] = useState(null);
  const doDeactivation = async () => {
    const res = await DeactivateUser();

    setUser({
      // Force dirty state validation.
      ...user,
      is_deactivated: !res.error ? true : user.is_deactivated,
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
              id="orgName"
              name="name"
              label="Organization Name"
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
          <Typography variant="h5">Email</Typography>
          {emailStatus && (
            <Alert severity={emailStatus.success ? "success" : "error"}>
              {emailStatus.message}
            </Alert>
          )}
          <form noValidate onSubmit={emailForm.handleSubmit}>
            <TextField
              variant="outlined"
              required
              fullWidth
              margin="normal"
              id="email"
              name="email"
              label="Email Address"
              type="email"
              value={emailForm.values.email}
              onChange={emailForm.handleChange}
              error={emailForm.touched.email && Boolean(emailForm.errors.email)}
              helperText={emailForm.touched.email && emailForm.errors.email}
            />
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={!emailForm.dirty}
            >
              Save
            </Button>
          </form>
        </CardContent>
      </FormCard>
      <FormCard>
        <CardContent>
          <Typography variant="h5">Password</Typography>
          {passwordStatus && (
            <Alert severity={passwordStatus.success ? "success" : "error"}>
              {passwordStatus.message}
            </Alert>
          )}
          <form noValidate onSubmit={passwordForm.handleSubmit}>
            <TextField
              variant="outlined"
              required
              fullWidth
              type="password"
              margin="normal"
              id="password"
              name="password"
              label="Password"
              value={passwordForm.values.password}
              onChange={passwordForm.handleChange}
              error={
                passwordForm.touched.password &&
                Boolean(passwordForm.errors.password)
              }
              helperText={
                passwordForm.touched.password && passwordForm.errors.password
              }
            />
            <TextField
              variant="outlined"
              required
              fullWidth
              type="password"
              margin="normal"
              id="confirm"
              name="confirm"
              label="Confirm Password"
              value={passwordForm.values.confirm}
              onChange={passwordForm.handleChange}
              error={
                passwordForm.touched.confirm &&
                Boolean(passwordForm.errors.confirm)
              }
              helperText={
                passwordForm.touched.confirm && passwordForm.errors.confirm
              }
            />
            <TextField
              variant="outlined"
              required
              fullWidth
              type="password"
              margin="normal"
              id="newpwd"
              name="newpwd"
              label="New Password"
              value={passwordForm.values.newpwd}
              onChange={passwordForm.handleChange}
              error={Boolean(passwordForm.errors.confirm)}
              helperText={
                passwordForm.touched.confirm && passwordForm.errors.confirm
              }
            />
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={!passwordForm.dirty}
            >
              Save
            </Button>
          </form>
        </CardContent>
      </FormCard>
      <FormCard>
        <CardContent>
          <Typography variant="h5">Point Conversion Rate</Typography>
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
              label="Points Per US Dollar"
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
            <strong>{user.is_deactivated ? "deactivated" : "activated"}</strong>
            .
          </Typography>
          <Button
            onClick={user.is_deactivated ? false : doDeactivation}
            variant="contained"
            color={user.is_deactivated ? "secondary" : "primary"}
            style={{ marginTop: 15 }} // FIXME
          >
            {user.is_deactivated
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
