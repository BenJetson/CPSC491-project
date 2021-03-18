import React from "react";
import Login from "./Login";

export default {
  title: "Login",
  component: Login,
};

const Template = (args) => <Login {...args} />;

export const Preview = Template.bind({});
Preview.args = {};
