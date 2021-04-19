import React from "react";
import { HashRouter as Router, Route, Switch } from "react-router-dom";
import { Container, CssBaseline, ThemeProvider } from "@material-ui/core";

import { AuthProvider, WithUser } from "./api/Auth";

import Application from "./components/Application";
import Homepage from "./components/Homepage";
import Login from "./components/Login";
import NavBar from "./components/NavBar";
import Registration from "./components/Registration";
import ScrollSpy from "./components/ScrollSpy";
import NotFound from "./components/NotFound";
import LoginRequired from "./components/LoginRequired";

import defaultTheme from "./Theme";

let App = () => {
  return (
    <ThemeProvider theme={defaultTheme}>
      <CssBaseline />
      <Router>
        <AuthProvider>
          <ScrollSpy />
          <WithUser>
            {({ isAuthenticated }) => (
              <Switch>
                <Route
                  exact
                  // Going to /login always shows the login page.
                  // If nobody is logged in, going to / will also show the
                  // login page.
                  path={isAuthenticated() ? "/login" : ["/login", "/"]}
                >
                  <Container>
                    <Login />
                  </Container>
                </Route>
                <Route>
                  <NavBar />
                  <Container>
                    {(isAuthenticated() && (
                      <Switch>
                        <Route path={"/register"}>
                          <Registration />
                        </Route>
                        <Route path={"/apply"}>
                          <Application
                            status={"pending"}
                            companyName={"Clemson Shipping"}
                            reason={"just too cool"}
                          />
                        </Route>
                        <Route exact path={"/"}>
                          <Homepage />
                        </Route>
                        <Route path={"*"}>
                          {/* If no route matches, show a not found page. */}
                          <NotFound />
                        </Route>
                      </Switch>
                    )) || (
                      // If nobody is logged in, simply show an error message
                      // on all pages directing to the login page.
                      <LoginRequired />
                    )}
                  </Container>
                </Route>
              </Switch>
            )}
          </WithUser>
        </AuthProvider>
      </Router>
    </ThemeProvider>
  );
};

export default App;
