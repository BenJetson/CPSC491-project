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
      {({ isOneOfRoles }) =>
        (isOneOfRoles([Roles.ADMIN]) && (
          <Switch>
            <Route path={`${match.path}/users`}>
              <NotImplemented feature={"Admin - Manage Users"} />
            </Route>
            <Route path={`${match.path}/organizations`}>
              <NotImplemented feature={"Admin - Manage Organizations"} />
            </Route>
            <Route path={"*"}>
              {/* If no route matches, show a not found page. */}
              <NotFound />
            </Route>
          </Switch>
        )) || <AccessDenied />
      }
    </WithUser>
  );
};

export default AppSubrouterAdmin;
