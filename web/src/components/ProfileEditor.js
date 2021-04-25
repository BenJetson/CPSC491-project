import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  ActivateUser,
  DeactivateUser,
  GetMyUser,
  UpdateUserEmail,
  UpdateUserName,
  UpdateUserPassword,
} from "../api/My";
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
  first_name: "",
  last_name: "",
  email: "",
  role_id: Roles.IDOf.DRIVER,
  is_deactivated: false,
};

const nameValidationSchema = yup.object({
  firstName: yup
    .string("Enter the new first name.")
    .required("First name is required."),
  lastName: yup
    .string("Enter the new last name.")
    .required("Last name is required."),
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

const MyProfileEditor = () => {
  const params = useParams();
  const userID = parseInt(params["userID"]) ?? false; // FIXME unchecked cast
  const isUpdate = userID !== false && userID > 0;

  const [user, setUser] = useState(emptyUser);

  useEffect(() => {
    if (!isUpdate) {
      setUser(emptyUser);
    }

    (async () => {
      const data = await GetMyUser();
      setUser(data);
    })();
  }, [userID, isUpdate]);

  const [nameStatus, setNameStatus] = useState(null);
  const nameForm = useFormik({
    initialValues: {
      firstName: user.first_name,
      lastName: user.last_name,
    },
    enableReinitialize: true,
    validationSchema: nameValidationSchema,
    onSubmit: async (values) => {
      const res = await UpdateUserName(
        userID,
        values.firstName,
        values.lastName
      );

      setUser({
        // Force dirty state validation.
        ...user,
        first_name: !res.error ? values.firstName : user.first_name,
        last_name: !res.error ? values.lastName : user.last_name,
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
      const res = await UpdateUserPassword(userID, values.password);

      setPasswordStatus(
        res.error
          ? { success: false, message: res.error }
          : { success: true, message: "Password changed successfully." }
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
      <Typography variant="h4">Edit User #{user.id}</Typography>
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
              id="firstName"
              name="firstName"
              label="First Name"
              value={nameForm.values.firstName}
              onChange={nameForm.handleChange}
              error={
                nameForm.touched.firstName && Boolean(nameForm.errors.firstName)
              }
              helperText={
                nameForm.touched.firstName && nameForm.errors.firstName
              }
            />
            <TextField
              variant="outlined"
              required
              fullWidth
              margin="normal"
              id="lastName"
              name="lastName"
              label="Last Name"
              value={nameForm.values.lastName}
              onChange={nameForm.handleChange}
              error={
                nameForm.touched.lastName && Boolean(nameForm.errors.lastName)
              }
              helperText={nameForm.touched.lastName && nameForm.errors.lastName}
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
            onClick={user.is_deactivated ? doActivation : doDeactivation}
            variant="contained"
            color={user.is_deactivated ? "secondary" : "primary"}
            style={{ marginTop: 15 }} // FIXME
          >
            {user.is_deactivated ? "activate" : "deactivate"} account
          </Button>
        </CardContent>
      </FormCard>
    </>
  );
};

export default MyProfileEditor;
