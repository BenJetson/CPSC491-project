import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

import { WithUser } from "./api/Auth";
import Roles from "./api/Roles";
import AccessDenied from "./components/AccessDenied";

const AppSubrouterAdmin = () => {
  const match = useRouteMatch();

  return (
    <WithUser>
      {({ isAuthenticated }) =>
        (!isAuthenticated() && (
          <Switch>
            <Route path={`${match.path}/register`}>
              <NotImplemented feature={"Account - Register"} />
            </Route>
            <Route path={`${match.path}/forgot`}>
              <NotImplemented feature={"Account - Forgot Password"} />
            </Route>
            <Route path={"*"}>
              {/* If no route matches, show a not found page. */}
              <NotFound />
            </Route>
          </Switch>
        )) || (
          <AccessDenied
            reason={
              "Must be logged out to use registration " +
              "or forgot password utilities."
            }
          />
        )
      }
    </WithUser>
  );
};

export default AppSubrouterAdmin;
