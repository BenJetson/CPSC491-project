import React from "react";
import NavBar from "./NavBar";

const StoryConfig = {
  title: "NavBar",
  component: NavBar,
};

export default StoryConfig;

const Template = (args) => <NavBar {...args} />;

export const Preview = Template.bind({});
Preview.args = {};
