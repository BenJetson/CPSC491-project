import React from "react";
import { HashRouter as Router, Switch, Route } from "react-router-dom";
import { Container, ThemeProvider } from "@material-ui/core";
import CssBaseline from "@material-ui/core/CssBaseline";
import NavBar from "./components/NavBar";
import Registration from "./components/Registration";

import defaultTheme from "./Theme";
import ScrollSpy from "./components/ScrollSpy";
import Login from "./components/Login";

//app should always go to login page unless logged in
//if logged in, app should default to home page
//need conditional render for navbar, should be available on pages except login

//New pages can be added by adding a new <Route> with the desired path name
//Home is a placeholder for later, when navbar is needed
let App = () => {
  return (
    <ThemeProvider theme={defaultTheme}>
      <CssBaseline />
      <Router>
        <ScrollSpy />
        <Switch>
          <Route path={"/login"}>
            <Container>
              <Login />
            </Container>
          </Route>
          <Route path={"/registration"}>
            <NavBar />
            <Container>
              <Registration />
            </Container>
          </Route>
          <Route path={"/home"}>
            <NavBar />
            <Container></Container>
          </Route>
        </Switch>
      </Router>
    </ThemeProvider>
  );
};

export default App;
