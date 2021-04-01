import React from "react";
import Registration from "./Registration";

const StoryConfig = {
  title: "Registration",
  component: Registration,
};

export default StoryConfig;

const Template = (args) => <Registration {...args} />;

export const Preview = Template.bind({});
Preview.args = {};
