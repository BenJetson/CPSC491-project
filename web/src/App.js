import React from "react";
import { HashRouter as Router, Route, Switch } from "react-router-dom";
import { Container, CssBaseline, ThemeProvider } from "@material-ui/core";

import { AuthProvider, WithUser } from "./api/Auth";

import Homepage from "./components/Homepage";
import Login from "./components/Login";
import NavBar from "./components/NavBar";
import ScrollSpy from "./components/ScrollSpy";
import NotFound from "./components/NotFound";
import LoginRequired from "./components/LoginRequired";

import defaultTheme from "./Theme";
import AppSubrouterAdmin from "./AppSubrouterAdmin";
import AppSubrouterDriver from "./AppSubrouterDriver";
import AppSubrouterSponsor from "./AppSubrouterSponsor";
import AppSubrouterStatic from "./AppSubrouterStatic";
import AppSubrouterMy from "./AppSubrouterMy";
import AppSubrouterAccount from "./AppSubrouterAccount";

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
                    <Switch>
                      <Route path={"/account"}>
                        <AppSubrouterAccount />
                      </Route>
                      <Route path={"/static"}>
                        <AppSubrouterStatic />
                      </Route>

                      {isAuthenticated() && (
                        <Route exact path={"/"}>
                          <Homepage />
                        </Route>
                      )}

                      {isAuthenticated() && (
                        <Route path={"/my"}>
                          <AppSubrouterMy />
                        </Route>
                      )}
                      {isAuthenticated() && (
                        <Route path={"/admin"}>
                          <AppSubrouterAdmin />
                        </Route>
                      )}
                      {isAuthenticated() && (
                        <Route path={"/sponsor"}>
                          <AppSubrouterSponsor />
                        </Route>
                      )}
                      {isAuthenticated() && (
                        <Route path={"/driver"}>
                          <AppSubrouterDriver />
                        </Route>
                      )}

                      {(isAuthenticated() && (
                        // If the user is logged in, but no route matches,
                        // show a not found page.
                        <Route path={"*"}>
                          <NotFound />
                        </Route>
                      )) || (
                        // If nobody is logged in, but no route matches, simply
                        // show an error message on all pages directing to the
                        // login page.
                        <Route path={"*"}>
                          <LoginRequired />
                        </Route>
                      )}
                    </Switch>
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
