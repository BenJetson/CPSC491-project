import React, { useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import {
  Box,
  Button,
  Checkbox,
  Container,
  FormControlLabel,
  Grid,
  Link,
  TextField,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { Alert } from "@material-ui/lab";
import * as yup from "yup";
import { useFormik } from "formik";

import { DoLogin } from "../api/Auth";

const validationSchema = yup.object({
  email: yup
    .string("Enter your email.")
    .email("Enter a valid email.")
    .required("Email is required."),
  password: yup
    .string("Enter your password.")
    // .min(8, "Password should be of minimum 8 characters length")
    .required("Password is required."),
});

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(20),
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
  },
  form: {
    width: "100%",
    marginTop: theme.spacing(1),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));

let Login = () => {
  const [error, setError] = useState(null);
  const classes = useStyles();
  const formik = useFormik({
    initialValues: {
      email: "",
      password: "",
    },
    validationSchema: validationSchema,
    onSubmit: async (values) => {
      const res = await DoLogin(values.email, values.password);
      setError(res.error);

      if (!res.error) {
        // Notice that this will trigger a full reload, not just using the
        // React Router here. This wlll force the context to reload.
        window.location.href = "/";
      }
    },
  });

  return (
    <Container component="main" maxWidth="xs">
      <Box className={classes.paper}>
        <img
          alt="logo"
          src="https://iconape.com/wp-content/files/zk/93042/svg/react.svg"
          height="192"
          width="192"
        />
        <form
          className={classes.form}
          noValidate
          onSubmit={formik.handleSubmit}
        >
          {error && <Alert severity="error">{error}</Alert>}

          <TextField
            color="secondary"
            variant="outlined"
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email Address"
            name="email"
            autoComplete="email"
            autoFocus
            value={formik.values.email}
            onChange={formik.handleChange}
            error={formik.touched.email && Boolean(formik.errors.email)}
            helperText={formik.touched.email && formik.errors.email}
          />
          <TextField
            color="secondary"
            variant="outlined"
            margin="normal"
            required
            fullWidth
            id="password"
            name="password"
            label="Password"
            type="password"
            autoComplete="current-password"
            value={formik.values.password}
            onChange={formik.handleChange}
            error={formik.touched.password && Boolean(formik.errors.password)}
            helperText={formik.touched.password && formik.errors.password}
          />
          <FormControlLabel
            control={<Checkbox value="remember" color="secondary" />}
            label="Remember me"
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            className={classes.submit}
            disabled={formik.isSubmitting}
          >
            Sign In
          </Button>
          <Grid container>
            <Grid item xs>
              <Link component={RouterLink} to="/forgotpassword" variant="body2">
                Forgot password?
              </Link>
            </Grid>
            <Grid item>
              <Typography variant="body2">
                Don't have an account?&nbsp;
                <Link component={RouterLink} to="/register">
                  Register
                </Link>
              </Typography>
            </Grid>
          </Grid>
        </form>
      </Box>
    </Container>
  );
};
export default Login;
