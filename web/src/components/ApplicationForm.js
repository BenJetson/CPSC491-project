import React, { useEffect, useState } from "react";
import { useFormik } from "formik";
import { useHistory } from "react-router-dom";
import { GetOrganizations } from "../api/Organizations";
import { SubmitApplication } from "../api/Applications";
import * as yup from "yup";
import {
  Button,
  Container,
  MenuItem,
  TextField,
  Typography,
  makeStyles,
} from "@material-ui/core";
import { Alert, AlertTitle } from "@material-ui/lab";

const useStyles = makeStyles((theme) => ({
  privacyWarning: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(2),
  },
  title: {
    marginBottom: theme.spacing(1),
  },
}));

const validationSchema = yup.object({
  organization: yup
    .number("Select a sponsoring organization.")
    .required("Must select a sponsor organization."),
  comment: yup
    .string("Enter a comment for your application.")
    .required("Comment cannot be blank."),
});

const ApplicationForm = () => {
  const history = useHistory();
  const classes = useStyles();

  const [error, setError] = useState(null);
  const [organizations, setOrganizations] = useState([]);

  const formik = useFormik({
    initialValues: {
      organization: "",
      comment: "",
    },
    validationSchema: validationSchema,
    onSubmit: async (values) => {
      const res = await SubmitApplication(values);
      setError(res.error);

      if (!res.error) {
        // Application submitted successfully.
        history.push("/driver/applications"); // FIXME
      }
    },
  });

  useEffect(
    () => {
      (async () => {
        try {
          const foundOrgs = await GetOrganizations();
          setOrganizations(foundOrgs);
        } catch {
          setError("Failed to fetch list of organizations.");
        }
      })();
    },
    [
      // Run this hook only once.
    ]
  );

  return (
    <Container maxWidth="xs">
      <Typography variant="h4" className={classes.title}>
        New Driver Application
      </Typography>
      <form onSubmit={formik.handleSubmit} noValidate>
        {error && <Alert severity="error">{error}</Alert>}

        <TextField
          color="secondary"
          variant="outlined"
          margin="normal"
          required
          fullWidth
          select
          name="organization"
          id="organization"
          label="Organization"
          value={formik.values.organization}
          onChange={formik.handleChange("organization")}
          error={
            formik.touched.organization && Boolean(formik.errors.organization)
          }
          helperText={formik.touched.organization && formik.errors.organization}
        >
          <MenuItem value={""}>None</MenuItem>
          {organizations.map((org) => (
            <MenuItem value={org.id} key={org.id}>
              {org.name}
            </MenuItem>
          ))}
        </TextField>

        <TextField
          color="secondary"
          variant="outlined"
          margin="normal"
          required
          fullWidth
          multiline
          rows={5}
          rowsMax={10}
          id="comment"
          label="Comment"
          name="comment"
          value={formik.values.comment}
          onChange={formik.handleChange}
          error={formik.touched.comment && Boolean(formik.errors.comment)}
          helperText={formik.touched.comment && formik.errors.comment}
        />

        <Alert severity="info" className={classes.privacyWarning}>
          <AlertTitle>Privacy Notice</AlertTitle>
          The organization you select will have access to your personal
          information after you submit this form.
        </Alert>

        <Button fullWidth type="submit" variant="contained" color="primary">
          Submit
        </Button>
      </form>
    </Container>
  );
};

export default ApplicationForm;
