import React from "react";
import FAQcard from "./FAQcard";
import { Typography, Box } from "@material-ui/core";

const FAQ = () => {
  return (
    <Box>
      <Typography variant="h5" component="h2">
        Frequently Asked Questions From Drivers
      </Typography>
      <FAQcard
        question="How do I apply to a sponsor?"
        answer="You can apply to a sponsor by going to driver tools -> Applications"
      />
      <FAQcard
        question="How can I view my points?"
        answer="You can view your balance by going to driver tools -> View Balance"
      />
      <FAQcard
        question="Where can I view the catalog?"
        answer="You can view your catalogs by going to driver tools -> Incentive Shop"
      />
      <FAQcard
        question="How do I earn points?"
        answer="Points are earned by following your sponsor's requirements for good driving"
      />
      <FAQcard
        question="How can I purchase an item?"
        answer="Items can be purchased fromt he incentive shop"
      />
      <Typography variant="h5" component="h2">
        Frequently Asked Questions From Sponsors
      </Typography>
      <FAQcard
        question="Where can I see my driver's information and add points?"
        answer="Driver information is located in sponsor tools -> Manage Drivers"
      />
      <FAQcard
        question="Where can I view and edit my catalog?"
        answer="You can edit the catelouge by going to sponsor tools -> Manage Catalog"
      />
      <FAQcard
        question="How can I make changes to my organization?"
        answer="You can make changes by going to sponsor tools -> Manage Organization"
      />
      <FAQcard
        question="How can I see driver's applications?"
        answer="you can view and accept applications by going to sponsor tools -> Manage Applications"
      />
    </Box>
  );
};

export default FAQ;
