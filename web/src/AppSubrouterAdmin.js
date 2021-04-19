import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import NotFound from "./components/NotFound";
import NotImplemented from "./components/NotImplemented";

const AppSubrouterAdmin = () => {
  const match = useRouteMatch();

  return (
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
  );
};

export default AppSubrouterAdmin;
