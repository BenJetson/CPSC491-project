import Application from "./Application";

const StoryConfig = {
  title: "Application",
  component: Application,
};

export default StoryConfig;

const Template = (args) => <Application {...args} />;

export const Denied = Template.bind({});
Denied.args = {
  companyName: "Very Cool Shipping Co",
  status: "Denied",
  reason: "Not cool enough",
};

export const Accepted = Template.bind({});
Accepted.args = {
  companyName: "Cool Shipping Co",
  status: "accepted",
  reason: "cool enough",
};
