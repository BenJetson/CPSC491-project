import React, { useEffect, useState } from "react";
import { useFormik } from "formik";
import { useHistory } from "react-router-dom";
import { GetOrganizations } from "../api/Organizations";
import { SubmitApplication } from "../api/Applications";
import * as yup from "yup";
import { Box, MenuItem, Select, TextField } from "@material-ui/core";
import { Alert } from "@material-ui/lab";

const validationSchema = yup.object({
  sponsor: yup.number().required("Must select a sponsor."),
  comment: yup
    .string("Enter a comment for your application.")
    .required("Comment cannot be blank."),
});

const ApplicationForm = () => {
  const history = useHistory();

  const [error, setError] = useState(null);
  const [organizations, setOrganizations] = useState([]);

  const formik = useFormik({
    initialValues: {
      organization: null,
      comment: "",
    },
    validationSchema: validationSchema,
    onSubmit: async (values) => {
      const res = await SubmitApplication(values);
      setError(res.error);

      if (!res.error) {
        // Application submitted successfully.
        history.push("/applications"); // FIXME
      }
    },
  });

  useEffect(() => {
    (async () => {
      setOrganizations(await GetOrganizations());
    })();
  });

  return (
    <Box>
      <form onSubmit={formik.handleSubmit} noValidate>
        {error && <Alert severity="error">{error}</Alert>}

        <Select
          value={formik.values.organization}
          onChange={formik.handleChange}
          error={
            formik.touched.organization && Boolean(formik.errors.organization)
          }
          helperText={formik.touched.organization && formik.errors.organization}
        >
          {organizations.map((org) => (
            <MenuItem value={org.id}>{org.title}</MenuItem>
          ))}
        </Select>

        <TextField
          color="secondary"
          variant="outlined"
          margin="normal"
          required
          fullWidth
          id="comment"
          label="Comment"
          name="comment"
          value={formik.values.comment}
          onChange={formik.handleChange}
          error={formik.touched.comment && Boolean(formik.errors.comment)}
          helperText={formik.touched.comment && formik.errors.comment}
        />
      </form>
    </Box>
  );
};

export default ApplicationForm;
