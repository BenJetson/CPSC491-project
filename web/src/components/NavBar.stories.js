import React from "react";
import NavBar from "./NavBar";

export default {
  title: "NavBar",
  component: NavBar,
};

const Template = (args) => <NavBar {...args} />;

export const Preview = Template.bind({});
Preview.args = {};
