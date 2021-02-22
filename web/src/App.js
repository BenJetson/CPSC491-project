import React from "react";
import { HashRouter as Router, Switch, Route } from "react-router-dom";
import { Container, ThemeProvider } from "@material-ui/core";

import defaultTheme from "./Theme";
import ScrollSpy from "./components/ScrollSpy";
import Login from "./components/Login";

//app should always go to login page unless logged in
//if logged in, app should default to home page
//need conditional render for navbar, should be available on pages except login

//New pages can be added by adding a new <Route> with the desired path name
let App = () => {
  return (
    <ThemeProvider theme={defaultTheme}>
      <Router>
        <ScrollSpy />
        <Container>
          <Switch>
            <Route path={"/login"}>
              <Login />
            </Route>
          </Switch>
        </Container>
      </Router>
    </ThemeProvider>
  );
};

export default App;