import React from "react";
import Login from "./Login";

const StoryConfig = {
  title: "Login",
  component: Login,
};

export default StoryConfig;

const Template = (args) => <Login {...args} />;

export const Preview = Template.bind({});
Preview.args = {};
