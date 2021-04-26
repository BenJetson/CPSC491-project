import React from "react";
import { Typography, Box } from "@material-ui/core";

const About = () => {
  return (
    <Box>
      <Typography variant="h4" component="h2">
        About our app
      </Typography>
      <Typography variant="h6" component="h2">
        This application was created as a senior design project by team XIV in
        the spring semester of 2021 for Computer Science 4910 at Clemson
        University. The driver incentive program intended to reward comercial
        drivers for driving safely throught a system of reward points provided
        by sponsors for use in the sponsors incentive shop to purchase items
        listed by the sponsors from Etsy.
      </Typography>
      <Typography variant="h4" component="h2">
        Team members
      </Typography>
      <Typography variant="h6" component="h2">
        Ben Godfrey, Team Leader
      </Typography>
      <Typography variant="h6" component="h2">
        Chloe Caples, Backend Developer
      </Typography>
      <Typography variant="h6" component="h2">
        Cynthia Brazil, Frontend Developer/tester
      </Typography>
      <Typography variant="h6" component="h2">
        Cameron Sharpe, Frontend Developer/tester
      </Typography>
      <Typography variant="h4" component="h2">
        Contact us
      </Typography>
      <Typography variant="h6" component="h2">
        help@teamxiv.space
      </Typography>
    </Box>
  );
};

export default About;
