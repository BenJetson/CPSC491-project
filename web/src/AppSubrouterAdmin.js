import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

import { WithUser } from "./api/Auth";
import Roles from "./api/Roles";
import AccessDenied from "./components/AccessDenied";

import UserList from "./components/UserList";
import AdminProfileEditor from "./components/AdminProfileEditor";

const AppSubrouterAdmin = () => {
  const match = useRouteMatch();

  return (
    <WithUser>
      {({ isOneOfRoles }) =>
        (isOneOfRoles([Roles.IDOf.ADMIN]) && (
          <Switch>
            <Route exact path={`${match.path}/users`}>
              <UserList />
            </Route>
            <Route path={`${match.path}/users/:userID`}>
              <AdminProfileEditor />
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
