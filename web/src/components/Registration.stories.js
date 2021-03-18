import React from "react";
import Registration from "./Registration";

export default {
  title: "Registration",
  component: Registration,
};

const Template = (args) => <Registration {...args} />;

export const Preview = Template.bind({});
Preview.args = {};
